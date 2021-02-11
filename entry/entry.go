package entry

import (
	"context"
	"github.com/pubgo/golug"
	"github.com/pubgo/golug/golug_entry"
	"github.com/pubgo/ossync/internal/ossync_db"
	"github.com/pubgo/ossync/internal/ossync_oss"
	"github.com/pubgo/ossync/models"
	"github.com/pubgo/ossync/rsync"
	"github.com/pubgo/ossync/version"
	"github.com/pubgo/xprocess"
	"go.uber.org/atomic"
	"os"
	"time"
)

func GetEntry() golug_entry.Entry {
	ent := golug.NewCtl(name)
	ent.OnCfg(&cfg)
	ent.Version(version.Version)
	ent.Description("sync from local to remote")
	ent.Commands(rsync.GetDbCmd())

	ent.BeforeStart(func() {
		ossync_db.InitDb(cfg.Db)
		ossync_oss.InitBucket(cfg.Oss)
	})

	ent.Register(func() {
		db := ossync_db.GetDb()
		_ = db.Sync2(
			new(models.SyncFile),
		)

		var waiter = rsync.NewWaiter()
		var run = func(path string) {
			key := os.ExpandEnv(path)
			bucket := ossync_oss.GetBucket()

			cancel := xprocess.GoLoop(func(ctx context.Context) {
				if waiter.Skip(key) {
					time.Sleep(5 * time.Second)
					return
				}

				var c = atomic.NewBool(false)
				defer waiter.Report(key, c)
				rsync.CheckAndSync(key, bucket, c)
				rsync.CheckAndDelete(bucket, c)
				rsync.CheckAndBackup(key, bucket)
				return
			})
			golug.AfterStop(func() { cancel() })
		}

		for i := range cfg.Files {
			run(cfg.Files[i])
		}

		select {}
	})

	return ent
}
