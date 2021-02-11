package version

import (
	"runtime"

	"github.com/pubgo/golug/golug_version"
)

var GoVersion = runtime.Version()
var GoPath = ""
var GoROOT = ""
var CommitID = ""
var Project = ""

func init() {
	if Version == "" {
		Version = "v0.0.1"
	}

	dix_trace.With(func(ctx *dix_trace.Ctx) {
		ctx.Func("ossync_version", func() interface{} {
			return golug_version.M{
				"build_time": BuildTime,
				"version":    Version,
				"go_version": GoVersion,
				"go_path":    GoPath,
				"go_root":    GoROOT,
				"commit_id":  CommitID,
				"project":    Project,
			}
		})
	})
}
