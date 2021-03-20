package main

import (
	"github.com/pubgo/golug"
	"github.com/pubgo/ossync/entry"
)

func main() {
	golug.Run(entry.GetEntry())
}
