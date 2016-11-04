package slack

import (
	"encoding/json"
	"errors"
	"reflect"

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

func buildIndexConfig(numOfRep, numOfShards int, object interface{}) string {
	// settings
	index := map[string]interface{}{"number_of_shards": numOfShards, "number_of_replicas": numOfRep}
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
