package scraper

import (
	"bytes"
	"decent-ft/src/JSlike"
	"decent-ft/src/event"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func (scr *Scraper) Request() error {
	scr.EventBus.Emit(Event_BeforeRequest)
	resp, err := http.Get(scr.Url.String())
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

func (scr *Scraper) Get(baseUrl string, query map[string]string) error {
	scr.EventBus = event.Bus{}
	scr.Url = CombineURL(baseUrl, query)
	if scr.AfterInit != nil {
		scr.AfterInit()
	}
	return scr.Request()
}

type Result struct {
	Resp *http.Response
	Buf  bytes.Buffer
	Json JSlike.Any
}

type Scraper struct {
	Url       url.URL
	Result    Result
	AfterInit func()
	EventBus  event.Bus
}

func Scrape() {

}
