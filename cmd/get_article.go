package cmd

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/tongsq/gin-example/component/logger"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"os"
	"time"
)

func init() {
	rootCmd.AddCommand(articleCmd)
}

var articleCmd = &cobra.Command{
	Use:   "get_article",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run:   getArticle,
}

func getArticle(cmd *cobra.Command, args []string) {

	//pool := component.NewTaskPool(20)
	for i := 1; i < 12; i++ {
		requestUrl := fmt.Sprintf("http://sydw.huatu.com/beikao/ms/%d.html", i)
		req, _ := http.NewRequest("GET", requestUrl, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
		req.Header.Set("Host", "sydw.huatu.com")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		client := http.Client{
			Timeout: time.Second * 5,
		}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error("http get error", err)
			return
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			logger.Error("http status error ", resp.StatusCode)
			return
		}
		body := resp.Body
		defer body.Close()

		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			logger.Error(err)
			return
		}
		//var proxy_list [][]string
		doc.Find("ul.listSty01 > li").Each(func(i int, selection *goquery.Selection) {
			td := selection.ChildrenFiltered("a").First()
			title := td.Text()
			url, _ := td.First().Attr("href")

			titleStr, err := simplifiedchinese.GBK.NewDecoder().String(title)

			if err != nil {
				panic(err)
			}
			fmt.Println(titleStr, url)
			getArticleContent(url, title)
		})

	}
	time.Sleep(time.Second * 20)
}

func getArticleContent(requestUrl string, title string) {

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	req.Header.Set("Host", "sydw.huatu.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "http://sydw.huatu.com/beikao/ms/")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("http get error", err)
		return
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		logger.Error("http status error ", resp.StatusCode)
		return
	}
	body := resp.Body
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		logger.Error(err)
		return
	}
	//var proxy_list [][]string
	var content []string
	doc.Find("div.artBcon > p").Each(func(i int, selection *goquery.Selection) {

		p := selection.Text()
		//a,_ := simplifiedchinese.GBK.NewEncoder().String("“")
		//b,_ := simplifiedchinese.GBK.NewEncoder().String("”")
		//p = strings.ReplaceAll(p, a, "\"")
		//p = strings.ReplaceAll(p, b, "\"")
		//pStr, err :=simplifiedchinese.GBK.NewDecoder().String(p)

		if err != nil {
			panic(err)
		}
		content = append(content, p)
	})
	if len(content) < 2 {
		doc.Find("div.mainWords > p").Each(func(i int, selection *goquery.Selection) {

			p := selection.Text()
			//a,_ := simplifiedchinese.GBK.NewEncoder().String("“")
			//b,_ := simplifiedchinese.GBK.NewEncoder().String("”")
			//p = strings.ReplaceAll(p, a, "\"")
			//p = strings.ReplaceAll(p, b, "\"")
			//pStr, err :=simplifiedchinese.GBK.NewDecoder().String(p)

			if err != nil {
				panic(err)
			}
			content = append(content, p)
		})
	}
	if len(content) < 2 {
		logger.Error("parse error")
		return
	}
	content = content[1 : len(content)-1]
	contentStr := title
	for _, str := range content {
		contentStr = contentStr + str
	}
	contentStr = contentStr + "\n\n\n\n"
	dir, _ := os.Getwd()
	file, err := os.OpenFile(dir+"/article.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error("打开文件失败")
	}
	file.WriteString(contentStr)
	file.Close()
}
