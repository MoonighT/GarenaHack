package main

import "github.com/MoonighT/GarenaHack/common"

func main() {
	common.LoggerInit("log/indexer.log", 3600*24, 1024*1024*1024, 100, 3)
	common.LogDetailf("hello")
}
