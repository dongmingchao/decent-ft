package main

import (
	"bytes"
	"github.com/dongmingchao/decent-ft@JSlike"
	"github.com/dongmingchao/decent-ft@scraper"
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

const idSourceUrl = "https://ppe-httpizza.ele.me/bdi.pinpoint_warehouse/aoi/crawler_fetch?timestamp=1571801944&limit=10"
const mapDataUrl = "https://ditu.amap.com/detail/get/detail"
const saveDataUrl = "https://ppe-httpizza.ele.me/bdi.pinpoint_warehouse/aoi/crawler_record"

func getIDs() []string {
	var ids []string
	scr := scraper.Scraper{}
	scr.Result.Json = &ids
	err := scr.Get(idSourceUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ids)
	return ids
}

func getMapData(id string) JSlike.Any {
	scr := scraper.Scraper{}
	var data mapAoi
	scr.Result.Json = &data
	err := scr.Get(mapDataUrl, map[string]string{
		"id": id,
	})
	data.pickData(scr.Result.Buf.Bytes())
	if err != nil {
		log.Printf("[ID %s]请求出错", ids[0])
		log.Println(err)
	} else {
		for key, _ := range data.rest["data"].(map[string]interface{}) {
			fmt.Println(key)
		}
		//res, err := json.MarshalIndent(data.rest["data"], "", "\t")
		//if err == nil {
		//	fmt.Println(string(res))
		//}
	}
	return data.rest["data"]
}

func saveData(data JSlike.Any) {
	scr := scraper.Scraper{}
	var back JSlike.Object
	scr.Result.Json = &back
	res, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	// 输出前16位检查返回正常
	println(string(res[0:16]))
	err = scr.Post(saveDataUrl, nil, bytes.NewReader(res))
	println(scr.Result.Buf.String())
}

func main() {
	ids := getIDs()
	if len(ids) == 0 {
		return
	}
	//for _, id := range ids {
	recData := getMapData(ids[0])
	//}
	saveData(recData)
}
