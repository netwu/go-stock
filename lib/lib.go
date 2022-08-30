package lib

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

func GetDatatimeByFrequency(frequency int) string {
	nTime := time.Now()
	if frequency == 0 {
		return nTime.Format("2006-01-02 15:04:01")
	}
	formatDatatime := ""

	day := nTime.Format("2006-01-02 ")
	nowMinute := time.Now().Minute()
	hour := time.Now().Hour()
	if nowMinute < frequency {
		nowMinute = 60
		if hour == 0 {
			hour = 24
			day = nTime.AddDate(0, 0, -1).Format("2006-01-02 ")
		}
		hour -= 1
	}

	getMinute := (int(math.Floor(float64(nowMinute/frequency))) - 1) * frequency

	if getMinute == 0 {
		formatDatatime = day + strconv.Itoa(hour) + ":00:00"
	} else {
		formatDatatime = day + strconv.Itoa(hour) + ":" + strconv.Itoa(getMinute) + ":00"
	}

	return formatDatatime
}

func Decimal2float64(d decimal.Decimal) float64 {
	f, exact := d.Float64()
	if !exact {
		return f
	}
	return 0
}

func String2float64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

func Isutf8(s string) bool {
	return utf8.ValidString(s)
}

func Int2String(i int) string {
	return strconv.Itoa(i)
}
func String2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err == nil {
		return i
	}
	return 0
}

func Read(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(content)
}
func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
func WriteFile(filename string, str string) {

	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	defer f.Close()

	_, err1 = f.WriteString(str) //写入文件(字符串)
	if err1 != nil {
		panic(err1)
	}
	f.Sync()
}
func SplitArray(arr []interface{}, pageSize int) []interface{} {
	if pageSize > len(arr) {
		return arr
	}
	var subArr []interface{}
	for i := 0; i < len(arr); i += pageSize {
		end := i + pageSize
		if (end) > len(arr) {
			end = len(arr)
		}
		subArr = append(subArr, arr[i:end])
	}
	return subArr
}

func ReplaceSpace(str string) string {
	str = strings.Replace(strings.Replace(str, "\n", "", -1), " ", "", -1)
	return str
}
func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func GetValueFromMap(mapObj map[string]interface{}, key string, valueType string, defaultValue interface{}) interface{} {
	if itf, ok := mapObj[key]; ok {
		switch itf.(type) {
		case string:
			value := itf.(string)
			switch valueType {
			case "string":
				return ok
			case "float64":
				return String2float64(value)
			case "int":
				return String2Int(value)
			}

		case int:
			value := itf.(int)
			switch valueType {
			case "string":
				return Int2String(value)
			case "float64":
				return float64(value)
			case "int":
				return value
			}
		case float64:
			value := itf.(float64)
			switch valueType {
			case "string":
				return fmt.Sprintf("%.5f", value)
			case "float64":
				return value
			case "int":
				return int(value)
			}
		}
	}

	if _, ok := mapObj[key]; !ok {
		return defaultValue
	}
	return nil
}

// 并发调用服务，每个handler都会传入一个调用逻辑函数
func GoroutineNotPanic(handlers []func() error) (err error) {

	var wg sync.WaitGroup
	// 假设我们要调用handlers这么多个服务
	for _, f := range handlers {

		wg.Add(1)
		// 每个函数启动一个协程
		go func(handler func() error) {

			defer func() {
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				if e := recover(); e != nil {
					// 某个服务调用协程报错，可以在这里打印一些错误日志
				}
				wg.Done()
			}()

			// 取第一个报错的handler调用逻辑，并最终向外返回
			e := handler()
			if err == nil && e != nil {
				err = e
			}
		}(f)
	}

	wg.Wait()

	return
}
