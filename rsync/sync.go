package rsync

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pubgo/golug/pkg/golug_utils"
	"github.com/pubgo/ossync/models"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
	"github.com/spf13/cobra"
	"github.com/twmb/murmur3"
	"go.uber.org/atomic"
)

func Hash(data []byte) (hash string) {
	var h = murmur3.New64()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func getCrc64Sum(path string) string {
	dt, err := ioutil.ReadFile(path)
	xerror.Panic(err)

	c := crc64.New(crc64.MakeTable(crc64.ECMA))
	xerror.PanicErr(c.Write(dt))
	return fmt.Sprintf("%d", c.Sum64())
}

// 本地文件加载
// 本地存储中，如果已经同步了，那么就不用同步了
var syncPrefix = "sync_files"
var delPrefix = "trash"
var backupPrefix = "backup"

func Md5(path string) string {
	dt, err := ioutil.ReadFile(path)
	xerror.Panic(err)

	c := md5.New()
	xerror.PanicErr(c.Write(dt))
	return base64.StdEncoding.EncodeToString(c.Sum(nil))
}

func CheckAndBackup(dir string, kk *oss.Bucket) {
	var handle = func(path string) {
		xlog.Infof("backup: %s", path)

		var g = xprocess.NewGroup()
		defer g.Wait()
		xerror.Exit(filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if info.Name()[0] == '.' {
					return filepath.SkipDir
				}
				return nil
			}

			// 隐藏文件
			if info.Name()[0] == '.' {
				return nil
			}

			key := filepath.Join(backupPrefix, path)
			xlog.Infof("backup: %s", path)
			g.Go(func(ctx context.Context) {
				xerror.Panic(kk.PutObjectFromFile(key, path, oss.ContentMD5(Md5(path))))
				time.AfterFunc(time.Second*5, func() { _ = os.Remove(path) })
			})

			return nil
		}))
	}

	xerror.Exit(filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if info.Name()[0] == '.' {
			return filepath.SkipDir
		}

		if info.Name() == backupPrefix {
			handle(path)
			return filepath.SkipDir
		}

		return nil
	}))
}

func CheckAndSync(dir string, kk *oss.Bucket, c *atomic.Bool) {
	if !golug_utils.PathExist(dir) {
		xlog.Warnf("path %s not found", dir)
		return
	}

	xlog.Infof("checking %s", dir)
	var handle = func(ctx context.Context, sf *models.SyncFile) {
		key := filepath.Join(syncPrefix, sf.Path)

		if !sf.Synced {
			head, err := kk.GetObjectMeta(key)
			if err != nil && !strings.Contains(err.Error(), "StatusCode=404") {
				xerror.Panic(err)
			}

			var ccc string
			if head != nil {
				ccc = head.Get("X-Oss-Hash-Crc64ecma")
			}

			if ccc != sf.Crc64ecma {
				xlog.Infof("sync: %s %s", key, sf.Path)
				xerror.Exit(kk.PutObjectFromFile(key, sf.Path, oss.ContentMD5(Md5(sf.Path))))
			}
			sf.Changed = true
			sf.Synced = true
		}

		if sf.Changed {
			c.Store(true)
			sf.Changed = false
			xlog.Infof("store: %s %s", key, sf.Path)
			models.SyncFileUpdate(sf, "path_hash=?", sf.PathHash)
		}
	}

	var g = xprocess.NewGroup()
	defer g.Wait()
	xerror.Exit(filepath.Walk(dir, func(path string, info os.FileInfo, err error) (gErr error) {
		defer xerror.RespErr(&gErr)

		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name()[0] == '.' || info.Name() == backupPrefix {
				return filepath.SkipDir
			}

			return nil
		}

		// 隐藏文件
		if info.Name()[0] == '.' {
			return nil
		}

		//if !strings.HasSuffix(info.Name(), ext) {
		//	return nil
		//}

		pathHash := Hash([]byte(path))
		sf := models.SyncFileFindOne("path_hash=?", pathHash)
		if sf == nil || sf.Path == "" {
			xlog.Debugf("ErrKeyNotFound: %s", path)
			sf = &models.SyncFile{
				Name:      info.Name(),
				Size:      info.Size(),
				Mode:      info.Mode(),
				ModTime:   info.ModTime().Unix(),
				IsDir:     info.IsDir(),
				Synced:    false,
				Changed:   true,
				Path:      path,
				PathHash:  pathHash,
				Crc64ecma: getCrc64Sum(path),
				Md5:       Md5(path),
			}
			models.SyncFileCreate(sf)

			g.Go(func(ctx context.Context) { handle(ctx, sf) })
			return nil
		}

		if sf.ModTime == info.ModTime().Unix() {
			return nil
		}

		sf.Name = info.Name()
		sf.Size = info.Size()
		sf.Mode = info.Mode()
		sf.ModTime = info.ModTime().Unix()
		sf.IsDir = info.IsDir()
		sf.Changed = true

		if hash := getCrc64Sum(path); sf.Crc64ecma != hash {
			sf.Synced = false
			sf.Crc64ecma = hash
		}

		g.Go(func(ctx context.Context) { handle(ctx, sf) })
		return nil
	}))
}

func CheckAndDelete(kk *oss.Bucket, c *atomic.Bool) {
	models.SyncFileEach(func(sf models.SyncFile) {
		if golug_utils.PathExist(sf.Path) {
			return
		}

		c.Store(true)
		xlog.Infof("delete:%s", sf.Path)

		xerror.Panic(ossRemove(kk, filepath.Join(syncPrefix, sf.Path), filepath.Join(delPrefix, sf.Path)))
		models.SyncFileDelete("path_hash=?", sf.PathHash)
	})
}

func GetDbCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "db"}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var prefix string
		if len(args) > 0 {
			prefix = args[0]
		}

		models.SyncFileEach(func(sf models.SyncFile) {
			if !strings.HasPrefix(sf.Path, prefix) {
				return
			}

			fmt.Println(sf)
		})
	}
	return cmd
}

func ossRemove(k *oss.Bucket, srcObjectKey, destObjectKey string) error {
	xlog.Infof("copy: %s %s", srcObjectKey, destObjectKey)
	_, err := k.CopyObject(srcObjectKey, destObjectKey)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode=404") {
			return nil
		}

		return xerror.Wrap(err)
	}

	xlog.Infof("delete: %s", srcObjectKey)
	return xerror.Wrap(k.DeleteObject(srcObjectKey))
}

func init() {
	//kk := golug_oss.GetClient()
	//resp := xerror.PanicErr(kk.ListObjectsV2(oss.Prefix(syncPrefix))).(oss.ListObjectsResultV2)
	//fmt.Println(resp.Prefix)
	//fmt.Println(resp.XMLName)
	//fmt.Println(resp.MaxKeys)
	//fmt.Println(resp.MaxKeys)
	//fmt.Println(resp.Delimiter)
	//fmt.Println(resp.IsTruncated)
	//fmt.Println(resp.CommonPrefixes)
	//for _, k := range resp.Objects {
	//	fmt.Printf("%#v\n", k)
	//}
}
