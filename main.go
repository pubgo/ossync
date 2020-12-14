package main

import (
	"github.com/pubgo/golug"
	"github.com/pubgo/ossync/rsync"
)

func main() {
	golug.Init()
	golug.Run(rsync.GetEntry())
}
