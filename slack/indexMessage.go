package slack

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MoonighT/GarenaHack/common"
	"github.com/MoonighT/elastic"
)

var (
	channelClient = &ChannelClient{}
	client        *elastic.Client
)

func InitIndexer(addrs []string) error {
	//load channel seq id
	channelClient.Init("127.0.0.1:6379", 1)
	client = NewESClient(addrs)
	if client == nil {
		return errors.New("client is nil")
	}
	return nil
}

func IndexMessage(message *Message) error {
	//get message seqid
	if message.Channelid == "" {
		return errors.New("empty channel id")
	}
	if message.Type == "" {
		message.Type = "message"
	}
	message.Text = strings.ToLower(message.Text)
	message.Filecontent = strings.ToLower(message.Filecontent)
	seqid, err := channelClient.Increase(message.Channelid)
	if err != nil {
		common.LogWarningf("increase seq id error %s", err)
		return err
	}
	message.Seqid = seqid
	docIdStr := fmt.Sprintf("%s_%d", message.Channelid, message.Seqid)
	message.Id = docIdStr
	resp, err := client.Index().
		Index(INDEX_NAME).
		Type(TABLE_NAME).
		Id(docIdStr).
		BodyJson(message).
		Do()
	if err != nil {
		common.LogWarningf("index message error %s", err)
	}
	common.LogDetailf("index message %v, resp=%v", message, resp)
	return nil
}
