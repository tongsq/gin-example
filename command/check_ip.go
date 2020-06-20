package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)
import _ "github.com/jinzhu/gorm/dialects/mysql"
import "github.com/tongsq/gin-example/model"
import "github.com/tongsq/gin-example/component"

func main() {

	db, err := gorm.Open("mysql", "python:123456@(127.0.0.1:3306)/py?charset=utf8&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SingularTable(true)
	var proxys []model.Proxy
	db.Where("status<>?", 1).Find(&proxys)
	fmt.Printf("count:%d, cap: %d\n", len(proxys), cap(proxys))
	pool := component.NewTaskPool(20)
	for _, proxy := range proxys {
		var proxy_tmp model.Proxy = proxy
		pool.RunTask(func() { CheckProxyStatus(proxy_tmp, db) })
	}

	time.Sleep(time.Second * 50)
}

func CheckProxyStatus(proxy model.Proxy, db *gorm.DB) {
	fmt.Printf("start check :host:%s, port:%s\n", proxy.Host, proxy.Port)
	result := CheckIp(proxy.Host, proxy.Port)
	fmt.Printf("%s, %s, the result is %v\n", proxy.Host, proxy.Port, result)
	if result {
		proxy.Status = 1
		db.Save(&proxy)
	}
}

func CheckIp(host, port string) bool {
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
