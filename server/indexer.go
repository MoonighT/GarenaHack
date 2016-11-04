package main

import (
	"net/http"

	"github.com/MoonighT/GarenaHack/common"
	"github.com/MoonighT/GarenaHack/slack"
)

var (
	addrs = []string{"http://127.0.0.1:9200"}
)

func main() {
	common.LoggerInit("log/indexer.log", 3600*24, 1024*1024*1024, 100, 3)
	common.LogDetailf("hello")

	err := slack.InitIndexer(addrs)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/index", slack.HandleIndex)
	http.HandleFunc("/search", slack.HandleSearchMessage)
	http.HandleFunc("/detail", slack.HandleGetMessageByChannel)
	common.LogFatal(http.ListenAndServe(":8080", nil))
}
