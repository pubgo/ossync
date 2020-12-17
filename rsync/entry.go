package rsync

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/golug"
	"github.com/pubgo/golug/golug_config"
	"github.com/pubgo/golug/golug_db"
	"github.com/pubgo/golug/golug_entry"
	"github.com/pubgo/golugin/golug_oss"
	"github.com/pubgo/ossync/api"
	"github.com/pubgo/ossync/models"
	"github.com/pubgo/ossync/version"
	"github.com/pubgo/xerror"
)

var name = "ossync"

func GetEntry() golug_entry.Entry {
	ent := golug.NewRestEntry(name)
	ent.Version(version.Version)
	ent.Description("sync from local to remote")
	ent.Commands(GetDbCmd())
	ent.Init(func() {
		golug_config.Decode(name, &cfg)
	})

	ent.Router("/", api.Router)

	golug.WithAfterStart(func(ctx *dix_run.AfterStartCtx) {
		kk := golug_oss.GetClient()
		resp := xerror.PanicErr(kk.ListObjectsV2(oss.Prefix(syncPrefix))).(oss.ListObjectsResultV2)
		fmt.Println(resp.Prefix)
		fmt.Println(resp.XMLName)
		fmt.Println(resp.MaxKeys)
		fmt.Println(resp.MaxKeys)
		fmt.Println(resp.Delimiter)
		fmt.Println(resp.IsTruncated)
		fmt.Println(resp.CommonPrefixes)
		//for _, k := range resp.Objects {
		//	fmt.Printf("%#v\n", k)
		//}
	})

	golug.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
		db := golug_db.GetClient()
		xerror.Exit(db.Sync2(
			new(models.SyncFile),
		))
	})

	return ent
}

func init() {
	//ent.Register(func() {
	//	defer xerror.RespExit()
	//
	//	bucket := golug_oss.GetClient()
	//	db := golug_badger.GetClient()
	//	defer db.Close()
	//
	//	var nw = newWaiter(0)
	//	var run = func(path string) {
	//		key := os.ExpandEnv(path)
	//
	//		xprocess.GoLoop(func(ctx context.Context) error {
	//			if nw.Skip(key) {
	//				time.Sleep(5 * time.Second)
	//				return nil
	//			}
	//
	//			var c = atomic.NewBool(false)
	//			defer nw.Report(key, c)
	//			checkAndSync(key, bucket, db, "", c)
	//			checkAndMove(bucket, db, c)
	//			checkAndBackup(key, bucket)
	//			return nil
	//		})
	//	}
	//	_ = run
	//
	//	//run("${HOME}/Documents")
	//	//run("${HOME}/Downloads")
	//	//run("${HOME}/git/docs")
	//	select {}
	//})
}
