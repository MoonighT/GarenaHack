package slack

const (
	ERROR_NONE    = 0
	ERROR_PARAM   = 1
	ERROR_UNKNOWN = -1
)

type Message struct {
	Id          string `json:"id" estype:"string" ana:"n"`
	Seqid       int64  `json:"seqid" estype:"long"`
	Text        string `json:"text" estype:"string"`
	Filecontent string `json:"filecontent" estype:"string"`
	Rawtext     string `json:"rawtext" estype:"string"`

	Userid    string `json:"userid" estype:"string" ana:"n"`
	Channelid string `json:"channelid" estype:"string" ana:"n"`
	Type      string `json:"type" estype:"string" ana:"n"`
	Subtype   int32  `json:"subtype" estype:"long"`
	Filemeta  string `json:"filemeta" estype:"string" ana:"n"`
	Timestamp int64  `json:"timestamp" estype:"long"`
}

type SearchRequest struct {
	Keyword   string `json:"keyword"`
	Cursor    int64  `json:"cursor"`
	Limit     int64  `json:"limit"`
	Channelid string `json:"channelid"`
	Userid    string `json:"userid"`
}

type SearchResponse struct {
	Errorcode  int32      `json:"errcode"`
	Totalcount int64      `json:"totalcount"`
	Messages   []*Message `json:"messages"`
}
