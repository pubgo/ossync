package main

import (
	"fmt"
	"github.com/pubgo/golug"
	"github.com/pubgo/ossync/entry"
	"path/filepath"
)

func main() {
	fmt.Println(filepath.Join("/ss","ff"))
	golug.Run(entry.GetEntry())
}
