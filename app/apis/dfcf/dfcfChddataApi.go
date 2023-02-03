package dfcf

// 1上证
// 0深证 北证

import (
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"goravel/lib"
	"goravel/lib/httpUtil"
	"goravel/lib/redisUtil"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/goravel/framework/facades"
)

type DfcfDetailResponse struct {
	Data struct {
		F116 float64
		F117 float64
	}
}

type DfcfKlineResponse struct {
	Data struct {
		Klines []string
	}
}

var symbolKey string = "chddata_symbols"
var ChddataKey string = "chddata"
var Fields = []string{"TOPEN", "TCLOSE", "HIGH", "LOW", "VOTURNOVER", "VATURNOVER", "AMPLITUDE", "PCHG", "CHG", "TURNOVER", "TCAP", "MCAP"}

// `topen`,`tclose`,high,  low,   voturnover     vaturnover    振幅    pchg  chg   turnover
// var apiHost = "https://data.eastmoney.com/stockcomment/api/600032.json"
var chddataApiHost = "http://push2his.eastmoney.com/api/qt/stock/kline/get?klt=101&fqt=1&fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61&rtntype=6&"

var detailApiHost = "http://push2.eastmoney.com/api/qt/stock/get?fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13,f57,f107,f162,f152,f167,f92,f59,f183,f184,f105,f185,f186,f187,f173,f188,f84,f116,f85,f117,f190,f189,f62,f55&secid=0.002271"

type ChddataService struct {
	redisUtil *redisUtil.RedisUtil
	attr      map[string]interface{}
}

func NewChddataService() *ChddataService {
	return &ChddataService{
		redisUtil.NewRedisUtil(),
		make(map[string]interface{}),
	}
}

func (chddataService *ChddataService) Test() {

	Text := `{
		"rc": 0,
		"rt": 17,
		"svr": 181669633,
		"lt": 1,
		"full": 0,
		"dlmkts": "",
		"data": {
			"code": "688147",
			"market": 1,
			"name": "微导纳米",
			"decimal": 2,
			"dktotal": 24,
			"preKPrice": 31.75,
			"prePrice": 31.75,
			"qtMiscType": 7,
			"klines": [
				"2023-02-02,33.01,38.10,38.10,32.60,79756,289881743.00,17.32,20.00,6.35,21.22"
			]
		}
	}`
	var dfcfKlineJson DfcfKlineResponse

	json.Unmarshal([]byte(Text), &dfcfKlineJson)
	klines := dfcfKlineJson.Data.Klines
	fmt.Println(dfcfKlineJson)
	var chddataMaps []map[string]interface{}

	for _, dataStr := range klines {
		if dataStr != "" {
			chddata := make(map[string]interface{})
			dataArr := strings.Split(dataStr, ",")
			if len(dataArr) > 3 {
				chddata["date"] = dataArr[0]
				DataItemList := dataArr[1:]

				for k, v := range DataItemList {
					chddata[Fields[k]] = lib.String2float64(v)
				}
				chddataMaps = append(chddataMaps, chddata)
			}
		}
	}
	fmt.Print(chddataMaps)
	// chddataService.redisUtil.RedisSet("Test", "123", 0)
}
func (chddataService *ChddataService) GetAllChddataMulity() {
	chddataService.PushRedis()
	chddataService.getLastDate()
	count := int(chddataService.redisUtil.RedisLLen(symbolKey))
	if count > 0 {
		pageSize := 20
		var pageCount int

		if count < pageSize {
			pageCount = 1
		} else {
			pageCount = int(math.Ceil(float64(count) / float64(pageSize)))
		}

		wgChddata := sync.WaitGroup{}
		for page := 0; page <= pageCount; page++ {
			for perPage := 0; perPage < pageSize; perPage++ {
				wgChddata.Add(1)
				go func() {
					chddataService.GetChddataMulity()
					wgChddata.Done()
				}()
			}
			wgChddata.Wait()
			// time.Sleep(time.Duration(1) * time.Second)
		}
	}
	chddateModel := models.Chddata{}
	chddateModel.UpdateChddataMonth()
	return
}

func (chddataService *ChddataService) PushRedis() {
	chddataService.redisUtil.DelKeyRedis(symbolKey)
	var results []map[string]interface{}
	facades.Orm.Query().Table("symbols").Order("id").Scan(&results)
	for _, symbol := range results {
		jsonString, err := json.Marshal(symbol)
		if err != nil {
			panic(err)
		}
		chddataService.redisUtil.DelKeyRedis("max_date:" + symbol["code"].(string))
		chddataService.redisUtil.PushRedis(symbolKey, string(jsonString))
	}

	fmt.Println("pushRedis success")
}

func (chddataService *ChddataService) getLastDate() {
	// sql := "from chddata group by code"
	var results []map[string]interface{}
	facades.Orm.Query().Table("chddata").Select("code,max(date) max_date ").Group("code").Scan(&results)
	for _, v := range results {
		chddataService.redisUtil.RedisSet("max_date:"+v["code"].(string), v["max_date"].(time.Time).Format("20060102"), 0)
	}
	return
}

func (chddataService *ChddataService) GetChddataMulity() {
	// func (chddataService *ChddataService) GetChddataMulity() {
	// facades.Log.Info("getChddataMulity")
	var symbol models.Symbols

	// symbolStr := chddataService.redisUtil.RedisRPop(chddataService.redisConn, symbolKey)
	symbolStr := chddataService.redisUtil.RedisRPop(symbolKey)

	json.Unmarshal([]byte(symbolStr), &symbol)

	chddatas := chddataService.getChddataFromSymbol(symbol, "")

	if chddatas != nil {
		chddateModel := models.Chddata{}
		chddateModel.Store(chddatas)
	}
	return

}

func (chddataService *ChddataService) getChddataFromSymbol(symbol models.Symbols, proxy string) []models.Chddata {
	var chddata models.Chddata
	var chddatas []models.Chddata
	if symbol.Code == "" {
		return nil
	}
	start := "20211001"
	max_date := chddataService.redisUtil.RedisGet("max_date:" + symbol.Code)
	if max_date != "" {
		start = max_date
	}
	end := time.Now().Format("20060102")
	if start == end {
		return nil
	}

	secid := symbol.Market + "." + symbol.Code
	// symbolnode := strings.Replace(symbol.Market, "sh", "0", -1)
	// symbolnode = strings.Replace(symbolnode, "sz", "1", -1)
	// secid=1.600517&beg=20221114&end=20221114
	api := chddataApiHost
	api += "&secid=" + secid + "&beg=" + start + "&end=" + end
	facades.Log.Info("requestapi:", api)

	content := httpUtil.Request(api, getHeader(), proxy)
	chddataMaps := chddataService.ParseChddataFromHtml(content)

	if chddataMaps != nil {
		for _, chddataMap := range chddataMaps {
			chddata.Market = symbol.Market
			chddata.Code = symbol.Code
			chddata.Name = symbol.Name
			chddata.Date = chddataMap["date"].(string)
			chddata.Tclose = 0
			if nil != chddataMap["TCLOSE"] {
				chddata.Tclose = chddataMap["TCLOSE"].(float64)
			}

			chddata.High = 0
			if nil != chddataMap["HIGH"] {
				chddata.High = chddataMap["HIGH"].(float64)
			}
			chddata.Low = 0
			if nil != chddataMap["LOW"] {
				chddata.Low = chddataMap["LOW"].(float64)
			}
			chddata.Topen = 0
			if nil != chddataMap["TOPEN"] {
				chddata.Topen = chddataMap["TOPEN"].(float64)
			}

			chddata.Chg = 0
			if nil != chddataMap["CHG"] {
				chddata.Chg = chddataMap["CHG"].(float64)
			}
			chddata.Pchg = 0
			if nil != chddataMap["PCHG"] {
				chddata.Pchg = chddataMap["PCHG"].(float64)

			}

			chddata.Turnover = 0
			if nil != chddataMap["TURNOVER"] {
				chddata.Turnover = chddataMap["TURNOVER"].(float64)
			}
			chddata.Voturnover = 0
			if nil != chddataMap["VOTURNOVER"] {
				chddata.Voturnover = chddataMap["VOTURNOVER"].(float64)
			}
			chddata.Vaturnover = 0
			if nil != chddataMap["VATURNOVER"] {
				chddata.Vaturnover = chddataMap["VATURNOVER"].(float64)
			}

			chddata.Tcap = 0
			if nil != chddataMap["TCAP"] {
				chddata.Tcap = chddataMap["TCAP"].(float64)
			}

			chddata.Mcap = 0
			if nil != chddataMap["MCAP"] {
				chddata.Mcap = chddataMap["MCAP"].(float64)
			}
			chddata.Amplitude = 0
			chddatas = append(chddatas, chddata)
		}
	}
	chddataService.redisUtil.DelKeyRedis("max_date:" + symbol.Code)
	return chddatas
}

func getHeader() map[string]string {
	httpHeader := make(map[string]string)
	httpHeader["User-Agent"] = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	// httpHeader["hexin-v"] = "AzQWCN396Xe14X26f07dtN8vA_mlDVj3mjHsO86VwL9COdov9h0oh-pBvMkd"
	// httpHeader["Referer"] = "quotes.money.163.com"
	// httpHeader["X-Requested-With"] = "XMLHttpRequest"

	return httpHeader
}

func (chddataService *ChddataService) ParseChddataFromHtml(htmlContent string) []map[string]interface{} {
	// var chddataMaps []map[string]interface{}
	var dfcfKlineJson DfcfKlineResponse

	json.Unmarshal([]byte(htmlContent), &dfcfKlineJson)
	klines := dfcfKlineJson.Data.Klines
	var chddataMaps []map[string]interface{}

	for _, dataStr := range klines {
		if dataStr != "" {
			chddata := make(map[string]interface{})
			dataArr := strings.Split(dataStr, ",")
			if len(dataArr) > 3 {
				chddata["date"] = dataArr[0]
				DataItemList := dataArr[1:]

				for k, v := range DataItemList {
					chddata[Fields[k]] = lib.String2float64(v)
				}
				chddataMaps = append(chddataMaps, chddata)
			}
		}
	}

	return chddataMaps
}
