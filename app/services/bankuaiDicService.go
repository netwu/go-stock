package services

import (
	"fmt"
	"goravel/lib/redisUtil"
	"strings"

	"github.com/goravel/framework/facades"
)

type BankuaiDicService struct {
	redisUtil *redisUtil.RedisUtil
	attr      map[string]interface{}
}

var bankuaiDicKey string = "bankuai_dic"

func NewBankuaiDicService() *BankuaiDicService {

	return &BankuaiDicService{
		redisUtil.NewRedisUtil(),
		make(map[string]interface{}),
	}
}
func (bankuaiDicService *BankuaiDicService) PushDicRedis() {
	var results []map[string]interface{}
	facades.Orm.Query().Table("bankuais").Order("id").Scan(&results)
	bankuaiMap := make(map[string]int)
	for _, symbol := range results {
		bankuai := symbol["bankuai"].(string)
		bankuaiList := strings.Split(bankuai, ",")
		for _, bankuaiItem := range bankuaiList {
			if _, ok := bankuaiMap[bankuaiItem]; ok {
				bankuaiMap[bankuaiItem] = bankuaiMap[bankuaiItem] + 1
			} else {
				bankuaiMap[bankuaiItem] = 1
			}
		}
		// facades.Log.Info("pushRedis bankuaiSymbolKey", bankuaiMap)
		// bankuaiService.redisUtil.PushRedis(bankuaiSymbolKey, string(jsonString))
	}
	for k, v := range bankuaiMap {
		sql := "INSERT INTO bankuai_dics (name,count) VALUES( '%s', %d ) ON DUPLICATE KEY UPDATE name='%s', count = %d"
		s := fmt.Sprintf(sql, k, v, k, v)
		facades.Log.Info("bankuaidic sql", s)
		facades.Orm.Query().Exec(s)
		// facades.Orm.Query().Raw(s).Exec(s)
		// facades.Orm.Query().Table("bankuai_dics").Where("name", "tom")

		// fmt.Println(k, v)
	}

	facades.Log.Info("pushRedis bankuaiSymbolKey")
	return
}
