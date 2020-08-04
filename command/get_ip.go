package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tongsq/gin-example/component"
	"github.com/tongsq/gin-example/component/proxy"
	"github.com/tongsq/gin-example/model"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

func getContentHtml(i int) io.ReadCloser {
	requestUrl := fmt.Sprintf("http://www.nimadaili.com/https/%d/", i)
	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
	req.Header.Set("Host", "www.nimadaili.com")
	req.Header.Set("Referer", "http://www.nimadaili.com/https/3/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return nil
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		fmt.Println("http status error ", resp.StatusCode)
		return nil
	}
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("read error", err)
	//	continue
	//}
	//bodyString := string(body)
	//fmt.Println(bodyString)
	return resp.Body
}

func getContentHtmlKuai(i int) io.ReadCloser {
	requestUrl := fmt.Sprintf("https://www.kuaidaili.com/free/inha/%d/", i)
	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
	req.Header.Set("Host", "www.kuaidaili.com")
	req.Header.Set("Referer", "www.kuaidaili.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Dest", "document")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return nil
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		fmt.Println("http status error ", resp.StatusCode)
		return nil
	}
	return resp.Body
}

func parseHtmlKuai(body io.ReadCloser) [][]string {
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var proxy_list [][]string
	doc.Find("tbody > tr").Each(func(i int, selection *goquery.Selection) {
		td := selection.ChildrenFiltered("td").Eq(0)
		proxy_host := td.Text()
		proxy_host = strings.Trim(proxy_host, " ")
		td = selection.ChildrenFiltered("td").Eq(1)
		proxy_port := td.Text()
		proxy_port = strings.Trim(proxy_port, " ")
		proxy_arr := []string{proxy_host, proxy_port}
		if proxy_host == "" || proxy_port == "" {
			fmt.Println("格式错误：", proxy_host, proxy_port)
		}
		proxy_list = append(proxy_list, proxy_arr)
	})
	return proxy_list
}

func parseHtml(body io.ReadCloser) [][]string {
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var proxy_list [][]string
	doc.Find("tbody > tr").Each(func(i int, selection *goquery.Selection) {
		td := selection.ChildrenFiltered("td").First()
		proxy_str := td.Text()
		proxy_str = strings.Trim(proxy_str, " ")
		proxy_arr := strings.Split(proxy_str, ":")
		if len(proxy_arr) != 2 {
			fmt.Println("格式错误:", proxy_str)
			return
		}
		proxy_list = append(proxy_list, proxy_arr)
	})
	return proxy_list
}

func main() {
	db, err := gorm.Open("mysql", "python:123456@(127.0.0.1:3306)/py?charset=utf8&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SingularTable(true)
	for i := 1; i < 20; i++ {
		contentBody := getContentHtmlKuai(i)
		if contentBody == nil {
			continue
		}
		proxy_list := parseHtmlKuai(contentBody)
		fmt.Println(proxy_list)
		//for _, proxy_arr := range proxy_list {
		//	result := proxy.CheckIpStatus(proxy_arr[0], proxy_arr[1])
		//	fmt.Println(result, proxy_arr)
		//}

		pool := component.NewTaskPool(20)
		for _, proxy_arr := range proxy_list {
			ip, port := proxy_arr[0], proxy_arr[1]
			pool.RunTask(func() { checkProxyAndSave(ip, port, db) })
		}
		var wg sync.WaitGroup = sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Println("wait 5s ...")
			time.Sleep(time.Second * 5)
		}(&wg)
		wg.Wait()
	}
	time.Sleep(time.Second * 50)
}

func checkProxyAndSave(host string, port string, db *gorm.DB) {
	result := proxy.CheckIpStatus(host, port)
	fmt.Println(result, host, port)
	var status int8 = 1
	if !result {
		status = 0
	}
	var proxyModel model.Proxy
	err := db.Where("host = ? AND port = ?", host, port).First(&proxyModel).Error

	if err != nil && gorm.IsRecordNotFoundError(err) {
		proxyModel = model.Proxy{
			Host:       host,
			Port:       port,
			Status:     status,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		db.Create(&proxyModel)
		return
	}
	proxyModel.Status = status
	proxyModel.UpdateTime = time.Now().Unix()
	db.Save(&proxyModel)
	return
}