package main

import (
	"github.com/pubgo/golug"
	"github.com/pubgo/ossync/rsync"
)

func main() {
	golug.Run(rsync.GetEntry())
}
