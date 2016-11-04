package main

import (
	"github.com/MoonighT/GarenaHack/common"
	"github.com/MoonighT/GarenaHack/slack"
)

var (
	addrs = []string{"http://127.0.0.1:9200"}
)

func main() {
	common.LoggerInit("log/create.log", 3600*24, 1024*1024*1024, 100, 3)
	err := slack.CreateIndex(addrs)
	if err != nil {
		common.LogWarningf("create index error")
	}
	return
}
