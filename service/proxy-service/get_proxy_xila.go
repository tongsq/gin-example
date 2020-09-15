package proxy_service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tongsq/gin-example/component/logger"
	"io"
	"net/http"
	"strings"
	"time"
)

type GetProxyXila struct {
}

func (s *GetProxyXila) GetContentHtml(i int) io.ReadCloser {
	requestUrl := fmt.Sprintf("http://www.xiladaili.com/https/%d/", i)
	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	req.Header.Set("Host", "www.xiladaili.com")
	req.Header.Set("Referer", "http://www.xiladaili.com/https/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	logger.Info("get proxy from xiladaili", i)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("http get error", err)
		return nil
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		logger.Error("http status error ", resp.StatusCode)
		return nil
	}
	return resp.Body
}

func (s *GetProxyXila) ParseHtml(body io.ReadCloser) [][]string {
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		logger.Error(err)
		return nil
	}
	var proxy_list [][]string
	doc.Find("tbody > tr").Each(func(i int, selection *goquery.Selection) {
		td := selection.ChildrenFiltered("td").First()
		proxy_str := td.Text()
		proxy_str = strings.Trim(proxy_str, " ")
		proxy_arr := strings.Split(proxy_str, ":")
		if len(proxy_arr) != 2 {
			logger.Error("格式错误:", proxy_str)
			return
		}
		proxy_list = append(proxy_list, proxy_arr)
	})
	return proxy_list
}
