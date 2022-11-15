package services

import (
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"goravel/lib"
	"goravel/lib/httpUtil"
	"strings"
	"sync"
)

var sum int32

var apiHost = "http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?sort=symbol&asc=1&num="

type SymbolJson struct {
	Symbol        string
	Code          string
	Name          string
	Trade         string
	Pricechange   float64
	Changepercent float64
	Buy           string
	Sell          string
	Settlement    string
	Open          string
	High          string
	Low           string
	Volume        float64
	Amount        float64
	Ticktime      string
	Per           float64
	Pb            float64
	Mktcap        float64
	Nmc           float64
	Turnoverratio float64
}

func GetAllStock() error {
	symbolModel := models.Symbols{}
	fmt.Println("start")
	pages := 40
	pagesize := "110"
	staockWg := sync.WaitGroup{}
	nodes := [2]string{"sh_a", "sz_a"}
	api := apiHost + pagesize

	for i := 0; i < len(nodes); i++ {
		apiTmp := api
		apiTmp += "&node=" + nodes[i]
		var j int32
		for j = 0; j < int32(pages); j++ {
			staockWg.Add(1)
			go func(n int32) {
				symbols := RequestApi(n, apiTmp)
				if len(symbols) > 0 {
					symbolModel.UpdateOrCreate(symbols)
				}
				staockWg.Done()
			}(j)
		}
		staockWg.Wait()
	}
	return nil

}
func RequestApi(i interface{}, apiTmp string) []models.Symbols {
	var symbolJsons []SymbolJson
	var symbols []models.Symbols

	n := i.(int32)
	n += 1
	// atomic.AddInt32(&sum, n)
	apiTmp += "&page=" + lib.String(n)
	result := httpUtil.HttpGet(apiTmp)
	fmt.Println(apiTmp)
	// lib.WriteFile(lib.String(n)+".txt", result)
	result = strings.Replace(result, " ", "", -1)

	json.Unmarshal([]byte(result), &symbolJsons)
	for _, symbolJson := range symbolJsons {
		var symbol models.Symbols
		symbol.Symbol = symbolJson.Symbol
		symbol.Code = symbolJson.Code
		symbol.Name = symbolJson.Name
		symbol.Trade = lib.String2float64(symbolJson.Trade)
		symbol.Pricechange = symbolJson.Pricechange
		symbol.Changepercent = symbolJson.Changepercent
		symbol.Buy = lib.String2float64(symbolJson.Buy)
		symbol.Sell = lib.String2float64(symbolJson.Sell)
		symbol.Settlement = lib.String2float64(symbolJson.Settlement)
		symbol.Open = lib.String2float64(symbolJson.Open)
		symbol.High = lib.String2float64(symbolJson.High)
		symbol.Low = lib.String2float64(symbolJson.Low)
		symbol.Volume = symbolJson.Volume
		symbol.Amount = symbolJson.Amount
		symbol.Ticktime = symbolJson.Ticktime
		symbol.Per = symbolJson.Per
		symbol.Pb = symbolJson.Pb
		symbol.Mktcap = symbolJson.Mktcap
		symbol.Nmc = symbolJson.Nmc
		symbol.Turnoverratio = symbolJson.Turnoverratio
		symbols = append(symbols, symbol)
	}

	return symbols
}
