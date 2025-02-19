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

// var detailApiHost = "http://push2.eastmoney.com/api/qt/stock/get?fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13,f57,f107,f162,f152,f167,f92,f59,f183,f184,f105,f185,f186,f187,f173,f188,f84,f116,f85,f117,f190,f189,f62,f55&secid=0.002271"

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

	if count <= 0 {
		return
	}
	// 优化并发处理
	pageSize := 50
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))

	for page := 0; page < pageCount; page++ {
		wgChddata := &sync.WaitGroup{}
		for perPage := 0; perPage < pageSize; perPage++ {
			wgChddata.Add(1)
			go chddataService.processSymbol(wgChddata)
		}
		wgChddata.Wait()
		// time.Sleep(time.Duration(100) * time.Millisecond)
	}

	// 更新月度数据
	chddateModel := models.Chddata{}
	chddateModel.UpdateChddataMonth()
}

// 新增方法：处理单个股票数据
func (chddataService *ChddataService) processSymbol(wg *sync.WaitGroup) {
	defer wg.Done()

	symbolStr := chddataService.redisUtil.RedisRPop(symbolKey)
	if symbolStr == "" {
		return
	}

	var symbol models.Symbols
	if err := json.Unmarshal([]byte(symbolStr), &symbol); err != nil {
		facades.Log.Error("解析股票数据失败:", err)
		return
	}

	if chddatas := chddataService.getChddataFromSymbol(symbol, ""); chddatas != nil {
		chddateModel := models.Chddata{}
		chddateModel.Store(chddatas)
	}
}

func (chddataService *ChddataService) getChddataFromSymbol(symbol models.Symbols, proxy string) []models.Chddata {

	if symbol.Code == "" {
		return nil
	}

	// 获取日期范围
	start := chddataService.getStartDate(symbol.Code)
	end := time.Now().Format("20060102")
	// 构建请求
	api := fmt.Sprintf("%s&secid=%s.%s&beg=%s&end=%s",
		chddataApiHost,
		symbol.Market,
		symbol.Code,
		start,
		end)
	// fmt.Println(api)
	content := httpUtil.Request(api, getHeader(), proxy)
	chddataMaps := chddataService.ParseChddata(content)

	if len(chddataMaps) == 0 {
		return nil
	}

	return chddataService.convertToChddatas(symbol, chddataMaps)
}

// 新增方法：获取开始日期
func (chddataService *ChddataService) getStartDate(code string) string {
	if maxDate := chddataService.redisUtil.RedisGet("max_date:" + code); maxDate != "" {
		return maxDate
	}
	return "20211001"
}

// 新增方法：转换数据
func (chddataService *ChddataService) convertToChddatas(symbol models.Symbols, chddataMaps []map[string]interface{}) []models.Chddata {
	chddatas := make([]models.Chddata, 0, len(chddataMaps))

	for _, chddataMap := range chddataMaps {
		chddata := models.Chddata{
			Market: symbol.Market,
			Code:   symbol.Code,
			Name:   symbol.Name,
			Date:   chddataMap["date"].(string),
		}

		// 设置数值字段
		if v, ok := chddataMap["TCLOSE"].(float64); ok {
			chddata.Tclose = v
		}
		if v, ok := chddataMap["HIGH"].(float64); ok {
			chddata.High = v
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

	// 清理缓存
	chddataService.redisUtil.DelKeyRedis("max_date:" + symbol.Code)
	return chddatas
}

func (chddataService *ChddataService) PushRedis() {
	// 判断队列长度
	count := int(chddataService.redisUtil.RedisLLen(symbolKey))
	if count > 0 {
		return
	}
	chddataService.redisUtil.DelKeyRedis(symbolKey)
	var results []map[string]interface{}
	facades.Orm.Query().Table("symbols").Where("name not like '%退市%'").Order("id").Scan(&results)
	for _, symbol := range results {
		jsonString, err := json.Marshal(symbol)
		if err != nil {
			panic(err)
		}

		chddataService.redisUtil.PushRedis(symbolKey, string(jsonString))
	}
	fmt.Println("pushRedis success")
}

func (chddataService *ChddataService) getLastDate() {
	var results []map[string]interface{}
	facades.Orm.Query().Table("chddata").Select("code,max(date) as max_date,max(updated_at) as updated_at").Group("code").Scan(&results)
	today := time.Now().Format("20060102")
	for _, v := range results {
		lastDay := v["max_date"].(time.Time).Format("20060102")
		if today == lastDay {
			updated_at := v["updated_at"].(time.Time)
			if updated_at.Hour() >= 15 {
				continue
			} else {
				// 删除今天的数据
				facades.Orm.Query().Table("chddata").
					Where("code = ?", v["code"]).
					Where("date = ?", lastDay).
					Delete(&models.Chddata{})
			}
		}
		chddataService.redisUtil.RedisSet("max_date:"+v["code"].(string), v["max_date"].(time.Time).Format("20060102"), 0)
	}
}

func (chddataService *ChddataService) GetChddataMulity() {
	var symbol models.Symbols
	symbolStr := chddataService.redisUtil.RedisRPop(symbolKey)
	json.Unmarshal([]byte(symbolStr), &symbol)
	chddatas := chddataService.getChddataFromSymbol(symbol, "")

	if chddatas != nil {
		chddateModel := models.Chddata{}
		chddateModel.Store(chddatas)
	}
	time.Sleep(time.Duration(1) * time.Second)
	// return

}

// func (chddataService *ChddataService) getChddataFromSymbol(symbol models.Symbols, proxy string) []models.Chddata {
// 	var chddata models.Chddata
// 	var chddatas []models.Chddata
// 	if symbol.Code == "" {
// 		return nil
// 	}
// 	start := "20211001"
// 	max_date := chddataService.redisUtil.RedisGet("max_date:" + symbol.Code)
// 	if max_date != "" {
// 		start = max_date
// 	}
// 	end := time.Now().Format("20060102")

// 	secid := symbol.Market + "." + symbol.Code
// 	// symbolnode := strings.Replace(symbol.Market, "sh", "0", -1)
// 	// symbolnode = strings.Replace(symbolnode, "sz", "1", -1)
// 	// secid=1.600517&beg=20221114&end=20221114
// 	api := chddataApiHost
// 	api += "&secid=" + secid + "&beg=" + start + "&end=" + end
// 	facades.Log.Info("requestapi:", api)

// 	content := httpUtil.Request(api, getHeader(), proxy)
// 	fmt.Println(content)
// 	// panic("getChddataFromSymbol")
// 	chddataMaps := chddataService.ParseChddata(content)

// 	for _, chddataMap := range chddataMaps {
// 		chddata.Market = symbol.Market
// 		chddata.Code = symbol.Code
// 		chddata.Name = symbol.Name
// 		chddata.Date = chddataMap["date"].(string)
// 		chddata.Tclose = 0
// 		if nil != chddataMap["TCLOSE"] {
// 			chddata.Tclose = chddataMap["TCLOSE"].(float64)
// 		}

// 		chddata.High = 0
// 		if nil != chddataMap["HIGH"] {
// 			chddata.High = chddataMap["HIGH"].(float64)
// 		}
// 		chddata.Low = 0
// 		if nil != chddataMap["LOW"] {
// 			chddata.Low = chddataMap["LOW"].(float64)
// 		}
// 		chddata.Topen = 0
// 		if nil != chddataMap["TOPEN"] {
// 			chddata.Topen = chddataMap["TOPEN"].(float64)
// 		}

// 		chddata.Chg = 0
// 		if nil != chddataMap["CHG"] {
// 			chddata.Chg = chddataMap["CHG"].(float64)
// 		}
// 		chddata.Pchg = 0
// 		if nil != chddataMap["PCHG"] {
// 			chddata.Pchg = chddataMap["PCHG"].(float64)

// 		}

// 		chddata.Turnover = 0
// 		if nil != chddataMap["TURNOVER"] {
// 			chddata.Turnover = chddataMap["TURNOVER"].(float64)
// 		}
// 		chddata.Voturnover = 0
// 		if nil != chddataMap["VOTURNOVER"] {
// 			chddata.Voturnover = chddataMap["VOTURNOVER"].(float64)
// 		}
// 		chddata.Vaturnover = 0
// 		if nil != chddataMap["VATURNOVER"] {
// 			chddata.Vaturnover = chddataMap["VATURNOVER"].(float64)
// 		}

// 		chddata.Tcap = 0
// 		if nil != chddataMap["TCAP"] {
// 			chddata.Tcap = chddataMap["TCAP"].(float64)
// 		}

// 		chddata.Mcap = 0
// 		if nil != chddataMap["MCAP"] {
// 			chddata.Mcap = chddataMap["MCAP"].(float64)
// 		}
// 		chddata.Amplitude = 0
// 		chddatas = append(chddatas, chddata)
// 	}

// 	chddataService.redisUtil.DelKeyRedis("max_date:" + symbol.Code)
// 	return chddatas
// }

func getHeader() map[string]string {
	httpHeader := make(map[string]string)
	httpHeader["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	httpHeader["Cookie"] = "perf_dv6Tr4n=1"
	httpHeader["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	// httpHeader["hexin-v"] = "AzQWCN396Xe14X26f07dtN8vA_mlDVj3mjHsO86VwL9COdov9h0oh-pBvMkd"
	// httpHeader["Referer"] = "quotes.money.163.com"
	// httpHeader["X-Requested-With"] = "XMLHttpRequest"

	return httpHeader
}

func (chddataService *ChddataService) ParseChddata(htmlContent string) []map[string]interface{} {
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
