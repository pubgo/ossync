package rsync

import (
	"context"
	"github.com/pubgo/golug/golug_config"
	"os"
	"time"

	"github.com/pubgo/golug"
	"github.com/pubgo/golug/golug_entry"
	"github.com/pubgo/golugin/golug_badger"
	"github.com/pubgo/golugin/golug_oss"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xprocess"
	"go.uber.org/atomic"
)

type Cfg struct {
}

func GetCfg() (cfg Cfg) {
	xerror.Next().Panic(golug_config.Decode(name, &cfg))
	return
}

var name = "ossync"

func GetEntry() golug_entry.Entry {
	ent := golug.NewCtlEntry(name)
	ent.Version("v0.0.1")
	ent.Description("sync from local to remote")
	ent.Commands(GetDbCmd())

	ent.Register(func() {
		defer xerror.RespDebug()

		bucket := golug_oss.GetClient()
		db := golug_badger.GetClient()
		defer db.Close()

		var nw = NewWaiter()
		var run = func(path string) {
			key := os.ExpandEnv(path)

			xprocess.GoLoop(func(ctx context.Context) error {
				if nw.Skip(key) {
					time.Sleep(5 * time.Second)
					return nil
				}

				var c = atomic.NewBool(false)
				defer nw.Report(key, c)
				checkAndSync(key, bucket, db, "", c)
				checkAndMove(bucket, db, c)
				checkAndBackup(key, bucket)
				return nil
			})
		}

		run("${HOME}/Documents")
		run("${HOME}/Downloads")
		run("${HOME}/git/docs")
		select {}
	})

	return ent
}
