package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/MoonighT/GarenaHack/common"
)

func JsonEncode(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		common.LogWarningf("json encode error %s", err)
		return ""
	}
	return string(data)
}

func HandleIndex(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		common.LogWarningf("read message error %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message := &Message{}
	err = json.Unmarshal(body, message)
	if err != nil {
		common.LogWarningf("unmarshal message error %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	common.LogDetailf("message = %v", message)
	err = IndexMessage(message)
	if err != nil {
		common.LogWarningf("index message error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func HandleSearchMessage(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		common.LogWarningf("read message error %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	searchreq := &SearchRequest{}
	err = json.Unmarshal(body, searchreq)
	if err != nil {
		common.LogWarningf("unmarshal message error %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	common.LogDetailf("search message request = %s", JsonEncode(searchreq))
	resp := SearchMessage(searchreq)
	common.LogDetailf("search message resp = %s", JsonEncode(resp))
	var data []byte
	data, err = json.Marshal(resp)
	if err != nil {
		common.LogWarningf("marshal resp error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func HandleGetMessageByChannel(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		common.LogWarningf("read message error %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	searchreq := &SearchRequest{}
	err = json.Unmarshal(body, searchreq)
	if err != nil {
		common.LogWarningf("unmarshal message error %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	common.LogDetailf("detail message request = %s", JsonEncode(searchreq))
	resp := GetMessagesByCursor(searchreq)
	common.LogDetailf("detail message resp = %s", JsonEncode(resp))
	var data []byte
	data, err = json.Marshal(resp)
	if err != nil {
		common.LogWarningf("marshal resp error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
