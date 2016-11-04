package slack

import (
	"encoding/json"
	"strings"

	"github.com/MoonighT/GarenaHack/common"
	"github.com/MoonighT/elastic"
)

func loadMessages(hits []*elastic.SearchHit) []*Message {
	result := make([]*Message, 0, len(hits))
	for _, hit := range hits {
		m := &Message{}
		err := json.Unmarshal(*(hit.Source), m)
		if err != nil {
			common.LogWarningf("load message json unmarshal error %s", err)
			continue
		}
		reshigh := hit.Highlight
		common.LogDetailf("highlight = %v", reshigh)
		if val, ok := reshigh["text"]; ok {
			m.Text = strings.Join(val, " ")
		}

		result = append(result, m)
	}
	return result
}

func SearchMessage(req *SearchRequest) *SearchResponse {
	resp := &SearchResponse{}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}
	keyword := strings.ToLower(req.Keyword)
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		resp.Errorcode = ERROR_NONE
		return resp
	}
	boolquery := elastic.NewBoolQuery()
	query := elastic.NewMultiMatchQuery(req.Keyword, "text", "filecontent")
	query.Operator("AND")
	boolquery.Must(query)
	if req.Cursor > 0 {
		//set cursor filter
		q := elastic.NewRangeQuery("seqid").Lt(req.Cursor)
		boolquery.Filter(q)
	}
	if req.Channelid != "" {
		//set channelid filter
		q := elastic.NewTermQuery("channelid", req.Channelid)
		boolquery.Filter(q)
	}
	if req.Userid != "" {
		//set userid filter
		q := elastic.NewTermQuery("userid", req.Userid)
		boolquery.Filter(q)
	}
	highlight := elastic.NewHighlight()
	//highlight.FragmentSize(500)
	highlight.Field("text")
	highlight.PreTags("餮餮")
	highlight.PostTags("犇犇")
	result, err := client.Search().
		Index(INDEX_NAME).
		Type(TABLE_NAME).
		Query(boolquery).
		Size(int(req.Limit)).
		Sort("timestamp", false).
		Pretty(true).
		Highlight(highlight).
		Do()
	if err != nil {
		common.LogWarningf("search message error %s", err)
		resp.Errorcode = ERROR_UNKNOWN
		return resp
	}
	resp.Totalcount = result.Hits.TotalHits
	resp.Messages = loadMessages(result.Hits.Hits)
	resp.Errorcode = ERROR_NONE
	return resp
}

func GetMessagesByCursor(req *SearchRequest) *SearchResponse {
	resp := &SearchResponse{}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}
	boolquery := elastic.NewBoolQuery()
	if req.Cursor == 0 || req.Channelid == "" {
		resp.Errorcode = ERROR_PARAM
		return resp
	}
	if req.Cursor > 0 {
		//set cursor filter
		q := elastic.NewRangeQuery("seqid").Lte(req.Cursor + req.Limit)
		boolquery.Filter(q)

		q2 := elastic.NewRangeQuery("seqid").Gte(req.Cursor - req.Limit)
		boolquery.Filter(q2)
	}
	if req.Channelid != "" {
		//set channelid filter
		q := elastic.NewTermQuery("channelid", req.Channelid)
		boolquery.Filter(q)
	}
	result, err := client.Search().
		Index(INDEX_NAME).
		Type(TABLE_NAME).
		Query(boolquery).
		Size(int(req.Limit*2+1)).
		Sort("timestamp", false).
		Pretty(true).
		Do()
	if err != nil {
		common.LogWarningf("search message error %s", err)
		resp.Errorcode = ERROR_UNKNOWN
		return resp
	}
	resp.Totalcount = result.Hits.TotalHits
	resp.Messages = loadMessages(result.Hits.Hits)
	resp.Errorcode = ERROR_NONE
	return resp
}
