package services

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

	"github.com/go-redis/redis"
	"github.com/goravel/framework/support/facades"
	"github.com/panjf2000/ants/v2"
)

var symbolKey string = "Chddatasymbols"
var ChddataKey string = "Chddata"
var Fields = []string{"TCLOSE", "HIGH", "LOW", "TOPEN", "CHG", "PCHG", "TURNOVER", "VOTURNOVER", "VATURNOVER", "TCAP", "MCAP"}

// var apiHost = "https://data.eastmoney.com/stockcomment/api/600032.json"
var chddataApiHost = "http://quotes.money.163.com/service/chddata.html?" + "fields=" + strings.Join(Fields, ";")

type ChddataService struct {
	redisConn *redis.Client
	attr      map[string]interface{}
}

func NewChddataService() *ChddataService {
	redisConn := redisUtil.InitRedisClient()
	return &ChddataService{redisConn, make(map[string]interface{})}
}
func (chddataService *ChddataService) GetAllChddataMulity() {
	chddataService.PushRedis()
	chddataService.getLastDate()
	count := int(redisUtil.RedisLLen(chddataService.redisConn, symbolKey))
	if count > 0 {
		pageSize := 5
		pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
		var wg sync.WaitGroup
		for page := 0; page < pageCount; page++ {
			p, _ := ants.NewPoolWithFunc(10, func(itf interface{}) {
				chddataService.GetChddataMulity(itf)
				wg.Done()
			})
			defer p.Release()
			for perPage := 0; perPage < pageSize; perPage++ {
				wg.Add(1)
				_ = p.Invoke(perPage)

			}
			wg.Wait()

			// time.Sleep(time.Duration(1) * time.Second)
		}
	}
	return
}

func (chddataService *ChddataService) PushRedis() {
	redisUtil.DelKeyRedis(chddataService.redisConn, symbolKey)
	var results []map[string]interface{}
	facades.DB.Table("symbols").Order("id").Scan(&results)
	for _, symbol := range results {
		jsonString, err := json.Marshal(symbol)
		if err != nil {
			panic(err)
		}
		redisUtil.DelKeyRedis(chddataService.redisConn, "max_date:"+symbol["code"].(string))
		redisUtil.PushRedis(chddataService.redisConn, symbolKey, string(jsonString))
	}

	fmt.Println("pushRedis success")
}

func (chddataService *ChddataService) getLastDate() {
	// sql := "from chddata group by code"
	var results []map[string]interface{}
	facades.DB.Table("chddata").Select("code,max(date) max_date ").Group("code").Scan(&results)
	for _, v := range results {
		redisUtil.RedisSet(chddataService.redisConn, "max_date:"+v["code"].(string), v["max_date"].(time.Time).Format("20060102"), 0)
	}
}

func (chddataService *ChddataService) GetChddataMulity(i interface{}) {
	// func (chddataService *ChddataService) GetChddataMulity() {
	// facades.Log.Info("getChddataMulity")
	var symbol models.Symbols

	// symbolStr := redisUtil.RedisRPop(chddataService.redisConn, symbolKey)
	symbolStr, err := chddataService.redisConn.RPop(symbolKey).Result()
	if err != nil {
		facades.Log.Error(err)
		// return err
	}
	json.Unmarshal([]byte(symbolStr), &symbol)

	chddatas := chddataService.getChddataFromSymbol(symbol, "")
	// ChddataRepo := _ChddataRepository.NewMysqlChddataRepository(chddataService.dbConn)

	if chddatas != nil {
		chddateModel := models.Chddata{}
		chddateModel.Store(&chddatas)
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
	max_data := redisUtil.RedisGet(chddataService.redisConn, "max_date:"+symbol.Code)
	if max_data != "" {
		start = max_data
	}
	end := time.Now().Format("20060102")
	if start == end {
		return nil
	}

	symbolnode := strings.Replace(symbol.Symbol, "sh", "0", -1)
	symbolnode = strings.Replace(symbolnode, "sz", "1", -1)
	api := chddataApiHost
	api += "&code=" + symbolnode + "&start=" + start + "&end=" + end
	facades.Log.Info("requestapi:", api)

	content := httpUtil.Request(api, getHeader(), proxy)
	chddataMaps := chddataService.ParseChddataFromHtml(content)

	if chddataMaps != nil {
		for _, chddataMap := range chddataMaps {
			chddata.Symbol = symbol.Symbol
			chddata.Code = symbol.Code
			chddata.Name = symbol.Name
			chddata.Date = chddataMap["date"].(string)
			chddata.Tclose = chddataMap["TCLOSE"].(float64)
			chddata.High = chddataMap["HIGH"].(float64)
			chddata.Low = chddataMap["LOW"].(float64)
			chddata.Topen = chddataMap["TOPEN"].(float64)
			chddata.Chg = chddataMap["CHG"].(float64)
			chddata.Pchg = chddataMap["PCHG"].(float64)
			chddata.Turnover = chddataMap["TURNOVER"].(float64)
			chddata.Voturnover = chddataMap["VOTURNOVER"].(float64)
			chddata.Vaturnover = chddataMap["VATURNOVER"].(float64)
			chddata.Tcap = chddataMap["TCAP"].(float64)
			chddata.Mcap = chddataMap["MCAP"].(float64)

			chddatas = append(chddatas, chddata)
		}
	}
	return chddatas
}

func getHeader() map[string]string {
	httpHeader := make(map[string]string)
	httpHeader["User-Agent"] = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	httpHeader["hexin-v"] = "AzQWCN396Xe14X26f07dtN8vA_mlDVj3mjHsO86VwL9COdov9h0oh-pBvMkd"
	httpHeader["Referer"] = "quotes.money.163.com"
	// httpHeader["X-Requested-With"] = "XMLHttpRequest"

	return httpHeader
}

func (chddataService *ChddataService) ParseChddataFromHtml(htmlContent string) []map[string]interface{} {
	var chddataMaps []map[string]interface{}
	csvList := strings.Split(htmlContent, "\r\n")

	csvDataList := csvList[1:]
	for _, dataStr := range csvDataList {
		if dataStr != "" {
			chddata := make(map[string]interface{})
			dataArr := strings.Split(dataStr, ",")

			chddata["date"] = dataArr[0]
			DataItemList := dataArr[3:]

			for k, v := range DataItemList {
				chddata[Fields[k]] = lib.String2float64(v)
			}
			chddataMaps = append(chddataMaps, chddata)
		}
	}
	return chddataMaps
}

func (chddataService *ChddataService) updateChddataMonth() {
	sql := "insert into chddata_months (code,name,month,avg_price) select code,name,left(date,7)as month,avg(tclose) as avg_price from chddata  group by code,month"
	facades.DB.Exec(sql)
}

// select t1.*,t2.month,t2.avg_price,t3.month,t3.avg_price from chddata_months t1,chddata_months t2,chddata_months t3 where t1.code=t2.code and t1.code=t3.code and t1.month='2021-10' and t2.month='2021-11' and t3.month='2021-12' and t1.name not like '%ST%' group by t1.month,t1.code having t3.avg_price>t2.avg_price and t2.avg_price>t1.avg_price
