package services

import (
	"encoding/json"
	"fmt"
	"goravel/app/models"
	"goravel/lib"
	"goravel/lib/httpUtil"
	"goravel/lib/redisUtil"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goravel/framework/facades"
)

var sql string = "truncate symbols"
var sinaApi = "https://finance.sina.com.cn/realstock/company/"

type BankuaiService struct {
	redisUtil *redisUtil.RedisUtil
	attr      map[string]interface{}
}

var bankuaiSymbolKey string = "bankuai_symbols"
var bankuaiKey string = "bankuais"

func NewBankuaiService() *BankuaiService {

	return &BankuaiService{
		redisUtil.NewRedisUtil(),
		make(map[string]interface{}),
	}
}
func (bankuaiService *BankuaiService) PushRedis() {
	var results []map[string]interface{}
	facades.Orm.Query().Table("symbols").Order("id").Scan(&results)
	bankuaiService.redisUtil.DelKeyRedis(bankuaiSymbolKey)
	for _, symbol := range results {
		jsonString, err := json.Marshal(symbol)
		if err != nil {
			panic(err)
		}
		bankuaiService.redisUtil.PushRedis(bankuaiSymbolKey, string(jsonString))
	}

	return
}
func (bankuaiService *BankuaiService) GetAllBankuaiMulity() {
	count := int(bankuaiService.redisUtil.RedisLLen(bankuaiSymbolKey))
	fmt.Print("count is ", count)
	if count == 0 {
		bankuaiService.PushRedis()
	}
	pageSize := 40
	var pageCount int
	if count < pageSize {
		pageCount = 1
	} else {
		pageCount = count / pageSize
	}

	for page := 0; page <= pageCount; page++ {
		wg := sync.WaitGroup{}
		for perPage := 0; perPage < pageSize; perPage++ {
			wg.Add(1)
			go func() {
				bankuaiService.GetBankuaiMulity()
				wg.Done()
			}()
		}
		wg.Wait()
		bankuaiCount := int(bankuaiService.redisUtil.RedisLLen(bankuaiKey))
		if bankuaiCount > 0 {
			bankuaiModel := models.Bankuais{}
			bankuaiStrs := bankuaiService.redisUtil.RedisBatchPop(bankuaiKey, bankuaiCount)
			var bankuaiSlice []models.Bankuais
			var bankuai models.Bankuais
			for _, v := range bankuaiStrs {
				json.Unmarshal([]byte(v), &bankuai)
				bankuaiSlice = append(bankuaiSlice, bankuai)
			}
			// facades.Log.Info("bankuai 入库:", bankuaiSlice)
			bankuaiModel.UpdateOrCreate(bankuaiSlice)
			time.Sleep(1 * time.Second)
		}
	}
	bankuaiService.SaveBankuaiDic()

}

func (bankuaiService *BankuaiService) GetBankuaiMulity() {
	var symbol models.Symbols
	// var bankuai stock.Bankuai
	symbolStr := bankuaiService.redisUtil.RedisRPop(bankuaiSymbolKey)
	if symbolStr != "" {
		json.Unmarshal([]byte(symbolStr), &symbol)
		bankuai := getBankuaisFromSymbol(symbol)
		jsonString, _ := json.Marshal(bankuai)
		bankuaiService.redisUtil.PushRedis(bankuaiKey, string(jsonString))
	}
	return
}
func getBankuaisFromSymbol(symbol models.Symbols) models.Bankuais {
	var bankuaiInfo models.Bankuais
	bankuaiInfo.Market = symbol.Market
	bankuaiInfo.Code = symbol.Code
	bankuaiInfo.Name = symbol.Name
	MarketMap := make(map[string]string)
	MarketMap["0"] = "sz"
	MarketMap["3"] = "sz"
	MarketMap["6"] = "sh"
	MarketMap["8"] = "bj"
	MarketMap["4"] = "bj"
	market := MarketMap[strings.Split(symbol.Code, "")[0]]
	api := sinaApi + market + symbol.Code + "/nc.shtml"

	facades.Log.Info(api)
	// api := sinaApi + symbol.Market + "/nc.shtml"
	content := httpUtil.Request(api, nil, "")
	bankuaiRet := getBankuaisFromHtml(content)
	bankuaiInfo.Zhuying = bankuaiRet["zhuying"]
	bankuaiInfo.Bankuai = bankuaiRet["bankuai"]
	return bankuaiInfo
}
func getBankuaisFromHtml(htmlContent string) map[string]string {
	companyMap := make(map[string]string)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	overviewText, _ := doc.Find(".com_overview").Html()
	// <p><b>
	var zhuyingSlice []string
	reg := regexp.MustCompile(`主营业务：.*?\n<p>(?s:(.*?))<\/p>`)

	//提取关键信息
	zhuyingStr := reg.FindAllStringSubmatch(overviewText, -1)

	for _, text := range zhuyingStr {
		zhuyingSlice = append(zhuyingSlice, lib.ReplaceSpace(text[1]))
	}
	reg = regexp.MustCompile(`主营业务：.*?\n.*?title="(.*?)">`)

	//提取关键信息
	zhuyingStr = reg.FindAllStringSubmatch(overviewText, -1)

	for _, text := range zhuyingStr {
		zhuyingSlice = append(zhuyingSlice, lib.ReplaceSpace(text[1]))
	}
	reg = regexp.MustCompile(`<a href="http://vip.stock.finance.sina.com.cn/.*?" target="_blank">(?s:(.*?))</a>`)
	if reg == nil {
		facades.Log.Error("MustCompile err")
	}
	//提取关键信息
	bankuaiStr := reg.FindAllStringSubmatch(overviewText, -1)
	var bankuaiSlice []string
	for _, text := range bankuaiStr {
		bankuaiSlice = append(bankuaiSlice, lib.ReplaceSpace(text[1]))
	}
	companyMap["zhuying"] = strings.Join(zhuyingSlice, ",")
	companyMap["bankuai"] = strings.Join(bankuaiSlice, ",")

	return companyMap
}

func (bankuaiService *BankuaiService) SaveBankuaiDic() {
	var results []map[string]interface{}
	facades.Orm.Query().Table("bankuais").Order("id").Scan(&results)
	bankuaiMap := make(map[string]int)
	for _, symbol := range results {
		bankuai := symbol["bankuai"].(string)
		bankuaiList := strings.Split(bankuai, ",")
		for _, bankuaiItem := range bankuaiList {
			if bankuaiItem != "" {
				if _, ok := bankuaiMap[bankuaiItem]; ok {
					bankuaiMap[bankuaiItem] = bankuaiMap[bankuaiItem] + 1
				} else {
					bankuaiMap[bankuaiItem] = 1
				}
			}
		}
	}
	for k, v := range bankuaiMap {
		sql := "INSERT INTO bankuai_dics (name,count) VALUES( '%s', %d ) ON DUPLICATE KEY UPDATE name='%s', count = %d"
		s := fmt.Sprintf(sql, k, v, k, v)
		// facades.Log.Info("bankuaidic sql", s)
		facades.Orm.Query().Exec(s)
	}

	facades.Log.Info("SaveBankuaiDic Success")
	return
}
