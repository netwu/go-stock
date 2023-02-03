package dfcf

import (
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"goravel/lib"
	"goravel/lib/httpUtil"
)

type Diff struct {
	F1   int
	F2   float64
	F3   float64
	F4   float64
	F5   float64
	F6   float64
	F7   float64
	F8   float64
	F9   float64
	F10  string
	F11  float64
	F12  string
	F13  int
	F14  string
	F15  float64
	F16  float64
	F17  float64
	F18  float64
	F20  float64
	F21  float64
	F22  float64
	F23  float64
	F24  float64
	F25  float64
	F62  float64
	F115 float64
	F128 string
	F140 string
	F141 string
	F136 string
	F152 int
}
type DfcfResData struct {
	Total int
	Diff  []Diff
}
type DfcfResponse struct {
	Rc     int
	Rt     int
	Svr    int
	Lt     int
	Full   int
	Dlmkts string
	Data   DfcfResData
}

var dfcfListApiHost = "http://4.push2.eastmoney.com/api/qt/clist/get?po=1&np=1&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13,f14,f15,f16,f17,f18,f20,f21,f22,f23,f24,f25,f62,f128,f136,f115,f152&pz="

func GetStock() error {

	symbolModel := models.Symbols{}
	fmt.Println("start")
	pages := 27
	pagesize := "200"
	// pages := 1
	// pagesize := "1"
	api := dfcfListApiHost + pagesize

	var j int32
	for j = 1; j <= int32(pages); j++ {
		apiTmp := api + "&pn=" + lib.String(j)
		// fmt.Println(apiTmp)
		symbols := dfcfRequestApi(apiTmp)
		if len(symbols) > 0 {
			symbolModel.UpdateOrCreate(symbols)
		}
	}

	return nil

}
func dfcfRequestApi(apiTmp string) []models.Symbols {
	var dfcfSymbolJsons DfcfResponse
	var symbols []models.Symbols
	result := httpUtil.HttpGet(apiTmp)
	json.Unmarshal([]byte(result), &dfcfSymbolJsons)
	for _, symbolJson := range dfcfSymbolJsons.Data.Diff {
		var symbol models.Symbols
		symbol.Market = lib.Int2String(symbolJson.F13)
		symbol.Code = symbolJson.F12
		symbol.Name = symbolJson.F14
		symbol.Trade = 0
		symbol.Pricechange = symbolJson.F4
		symbol.Changepercent = symbolJson.F3
		symbol.Buy = symbolJson.F2
		symbol.Sell = 0
		symbol.Settlement = symbolJson.F18
		symbol.Open = symbolJson.F17
		symbol.High = symbolJson.F15
		symbol.Low = symbolJson.F16
		symbol.Volume = symbolJson.F5
		symbol.Amount = symbolJson.F6
		symbol.Ticktime = ""
		symbol.Per = symbolJson.F9
		symbol.Pb = symbolJson.F23
		symbol.Mktcap = symbolJson.F20
		symbol.Nmc = symbolJson.F21
		symbol.Turnoverratio = symbolJson.F8
		symbols = append(symbols, symbol)
	}

	return symbols
}
