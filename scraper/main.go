package scraper

import (
	"bytes"
	"encoding/json"
	"github.com/dongmingchao/decent-ft/JSlike"
	"github.com/dongmingchao/decent-ft/event"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (scr *Scraper) Request(data io.Reader, header map[string]string) error {
	scr.EventBus.Emit(Event_BeforeRequest)
	var resp *http.Response
	var err error
	req, _ := http.NewRequest(scr.Method, scr.Url.String(), data)
	if scr.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err = http.DefaultClient.Do(req)
	scr.Result.Resp = resp
	if err == nil {
		scr.Result.Buf.ReadFrom(resp.Body)
		scr.EventBus.Emit(Event_BeforeUnmarshal)
		err = json.Unmarshal(scr.Result.Buf.Bytes(), scr.Result.Json)
	}
	return err
}

func CombineURL(base string, appendQuery map[string]string) url.URL {
	reqUrl, err := url.Parse(base)
	if err != nil {
		log.Fatalf("URL格式不正确！%s", base)
	}
	query := reqUrl.Query()
	for k, v := range appendQuery {
		query.Set(k, v)
	}
	reqUrl.RawQuery = query.Encode()
	return *reqUrl
}

type Method byte

const (
	GET Method = 1 << iota
	POST
)

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	}
	return ""
}

func (scr *Scraper) Ready(method Method, baseUrl string, query map[string]string) {
	scr.EventBus = event.Bus{}
	scr.Url = CombineURL(baseUrl, query)
	scr.Method = method.String()
	if scr.AfterInit != nil {
		scr.AfterInit()
	}
}

func (scr *Scraper) Get(baseUrl string, query map[string]string) error {
	scr.Ready(GET, baseUrl, query)
	return scr.Request(nil, nil)
}

func (scr *Scraper) Post(baseUrl string, query map[string]string, data io.Reader) error {
	scr.Ready(POST, baseUrl, query)
	return scr.Request(data, nil)
}

type Result struct {
	Resp *http.Response
	Buf  bytes.Buffer
	Json JSlike.Any
}

type Scraper struct {
	Url       url.URL
	Method    string
	Result    Result
	AfterInit func()
	EventBus  event.Bus
}

func Scrape() {

}
