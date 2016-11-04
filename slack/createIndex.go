package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/MoonighT/GarenaHack/common"
	"github.com/MoonighT/elastic"
)

const (
	INDEX_NAME = "slackindex"
	TABLE_NAME = "messagetab"
)

func NewESClient(addrs []string) *elastic.Client {
	common.LogDetailf("client urls = %v", addrs)
	client, err := elastic.NewClient(elastic.SetURL(addrs...))
	if err != nil {
		common.LogWarningf("create es client error %s", err)
		return nil
	}
	return client
}

// Add unit test for this function
func getAsciiMappingArray() (charsMappingArray []string) {
	//Ascii characters other than [0-9], [a-zA-Z] will be converted into a white space
	for i := 33; i < 128; i++ {
		if i < 32 || (i > 32 && i < 48) ||
			(i > 57 && i < 65) || (i > 90 && i < 97) || (i > 122) {
			var x string
			if i == 39 {
				//continue
				x = "'=>\\u0020"
			} else {
				r := rune(i)
				s := fmt.Sprintf("%q", r)
				s = strings.Replace(s, "'", "", -1)
				x = fmt.Sprintf("%s=>\\u0020", s)
			}
			charsMappingArray = append(charsMappingArray, x)
		}
	}
	return
}

func buildIndexConfig(numOfRep, numOfShards int, object interface{}) string {
	// settings
	charsMappingName := "char_mapping"
	charsMappingRule := map[string]interface{}{"type": "mapping",
		"mappings": getAsciiMappingArray()}
	charsMapping := map[string]interface{}{charsMappingName: charsMappingRule}
	slackAnalyzer := map[string]interface{}{"tokenizer": "whitespace",
		"char_filter": []string{charsMappingName}} //
	analyzer := map[string]interface{}{"slack_analyzer": slackAnalyzer}
	analysis := map[string]interface{}{"char_filter": charsMapping, "analyzer": analyzer}
	index := map[string]interface{}{"number_of_shards": numOfShards, "number_of_replicas": numOfRep,
		"analysis": analysis}
	settings := map[string]interface{}{"index.cache.query.enable": true, "index": index}
	tp := reflect.TypeOf(object)
	properties := map[string]interface{}{}
	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		estype := field.Tag.Get("estype")
		name := field.Tag.Get("json")
		ana := field.Tag.Get("ana")
		n := map[string]interface{}{"type": estype}
		if ana == "n" {
			n["index"] = "not_analyzed"
		} else if ana == "s" {
			n["analyzer"] = "slack_analyzer"
		}
		properties[name] = n
	}
	table := map[string]interface{}{"properties": properties,
		"_routing": map[string]bool{"required": false}}

	mappings := make(map[string]interface{})
	mappings[TABLE_NAME] = table
	// mappings
	config := map[string]interface{}{"settings": settings, "mappings": mappings}
	js, err := json.Marshal(config)
	if err != nil {
		common.LogWarningf("Encode json error=%s", err)
		return ""
	}
	return string(js)
}

func CreateIndex(addrs []string) error {
	client := NewESClient(addrs)
	if client == nil {
		return errors.New("client is nil")
	}
	mapping := buildIndexConfig(0, 3, Message{})
	createIndex, err := client.CreateIndex(INDEX_NAME).BodyString(mapping).Do()
	if err != nil || !createIndex.Acknowledged {
		common.LogWarningf("Failed to create index, err=%s", err)
		return err
	}
	common.LogDetailf("create index finish")
	return nil
}
