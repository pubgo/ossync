package entry

import (
	"context"
	"os"
	"time"

	"github.com/pubgo/golug"
	"github.com/pubgo/golug/golug_entry"
	"github.com/pubgo/ossync/internal/ossync_db"
	"github.com/pubgo/ossync/internal/ossync_oss"
	"github.com/pubgo/ossync/models"
	"github.com/pubgo/ossync/rsync"
	"github.com/pubgo/ossync/version"
	"github.com/pubgo/xprocess"
	"go.uber.org/atomic"
)

func GetEntry() golug_entry.Entry {
	ent := golug.NewTaskEntry(name)
	ent.OnCfg(&cfg)
	ent.Version(version.Version)
	ent.Description("sync from local to remote")
	ent.Commands(rsync.GetDbCmd())

	golug.BeforeStart(func(ctx *dix_run.BeforeStartCtx) {
		ossync_db.InitDb(cfg.Db)
		ossync_oss.InitBucket(cfg.Oss)
	})

	golug.AfterStart(func(ctx *dix_run.AfterStartCtx) {
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
			golug.AfterStop(func(ctx *dix_run.AfterStopCtx) { cancel() })
		}

		for i := range cfg.Files {
			run(cfg.Files[i])
		}
	})

	return ent
}
