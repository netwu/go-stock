package httpUtil

import (
	"crypto/tls"
	"fmt"
	"goravel/lib"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func Request(webUrl string, header map[string]string, proxy string) string {
	/*
		1. 代理请求
		2. 跳过https不安全验证
		3. 自定义请求头 User-Agent

	*/
	// webUrl := "http://ip.gs/"

	request, _ := http.NewRequest("GET", webUrl, nil)
	for k, v := range header {
		v = strings.Replace(v, "\n", "", -1)
		request.Header.Set(k, v)
	}
	// client := &http.Client{
	// 	Timeout: time.Second * 5, //超时时间
	// }
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		// proxyUrl := "http://112.80.248.73:80"
		// proxyUrl := redisGet
		proxy, _ := url.Parse(proxy)
		tr = &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("出错了", err)
		return ""
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	if value, ok := header["content-encoding"]; ok {
		if value == "gzip" {
			body, _ = lib.GzipDecode(body)
		}
	}

	if !lib.Isutf8(string(body)) {
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body)
		return string(utf8Data)
	}
	// fmt.Println(string(body))
	return string(body)

}

func HttpGet(url string) string {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	if !lib.Isutf8(string(body)) {
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body)
		if err != nil {
			// handle error
			panic(err)

		}
		return string(utf8Data)
	}
	return string(body)

}
