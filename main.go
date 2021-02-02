package main

import (
	"github.com/google/gops/agent"
	"github.com/pubgo/golug"
	"github.com/pubgo/ossync/entry"
	"github.com/pubgo/xerror"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	//xerror.ExitErr(profiler.Start(profiler.Config{
	//	ApplicationName: "ossync",
	//	ServerAddress:   "http://localhost:4040", // this will run inside docker-compose, hence `pyroscope` for hostname
	//}))

	lis, err := net.Listen("tcp", ":8088")
	xerror.Exit(err)
	go http.Serve(lis, nil)
	xerror.Exit(agent.Listen(agent.Options{}))
	golug.Run(entry.GetEntry())
}
