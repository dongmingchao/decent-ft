package main

import (
	"decent-ft/src/JSlike"
	"decent-ft/src/scraper"
	"encoding/json"
	"fmt"
	"log"
)

type mapAoi struct {
	Status string        `json:"status"`
	rest   JSlike.Object `json:"-"`
}

func (aoi *mapAoi) pickData(buf []byte) error {
	var err error
	switch aoi.Status {
	case "1":
		err = json.Unmarshal(buf, &aoi.rest)
	case "8":
		//httpErr.ErrCode = http.StatusNotFound
		//httpErr.ErrMsg = "404 Not found!"
		//err = &httpErr
	case "3":
		//httpErr.ErrCode = http.StatusBadRequest
		//httpErr.ErrMsg = "参数错误!"
		//err = &httpErr
	}
	return err
}

func main() {
	var ids []string
	scr := scraper.Scraper{}
	scr.Result.Json = &ids
	err := scr.Get(idSourceUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ids)
	if len(ids) == 0 {
		return
	}
	//for _, id := range ids {
	//httpErr := Error{
	//	ErrUrl: reqUrl,
	//}
	scr = scraper.Scraper{}
	scr.AfterInit = func() {
		println(scr.Url.String(), "after init")
		scr.EventBus.On(scraper.Event_BeforeRequest, func(any ...JSlike.Any) {
			println("before request")
			println(scraper.Event_BeforeRequest)
		})
	}
	var data mapAoi
	scr.Result.Json = &data
	err = scr.Get(mapDataUrl, map[string]string{
		"id": ids[0],
	})
	data.pickData(scr.Result.Buf.Bytes())
	if err != nil {
		log.Printf("[ID %s]请求出错", ids[0])
		log.Println(err)
	} else {
		//for key, _ := range data.rest["data"].(map[string]interface{}) {
		//	fmt.Println(key)
		//}
		//res, err := json.MarshalIndent(data.rest["data"], "", "\t")
		//if err == nil {
		//	fmt.Println(string(res))
		//}
	}
	//}
}
