package proxy

import (
	"fmt"
	"github.com/tongsq/gin-example/component/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func CheckIpStatusActive(host, port string) bool {
	request_url := "https://www.baidu.com"
	req, _ := http.NewRequest("GET", request_url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
	proxyServer := fmt.Sprintf("http://%s:%s", host, port)
	proxyUrl, _ := url.Parse(proxyServer)
	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("http get error", err)
		return false
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("http read error", err)
		return false
	}
	return true
}

func CheckIpStatus(host, port string) bool {
	request_url := "https://www.c5game.com/api/product/sale.json?id=2705689&page=1&sort=1&key=1523539522"
	req, _ := http.NewRequest("GET", request_url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
	proxyServer := fmt.Sprintf("http://%s:%s", host, port)
	proxyUrl, _ := url.Parse(proxyServer)
	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return false
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err)
		return false
	}
	//fmt.Println(string(body))
	return true
}
