package cmd

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/tongsq/gin-example/component"
	"github.com/tongsq/gin-example/component/logger"
	"github.com/tongsq/gin-example/model"
	"github.com/tongsq/gin-example/service"
	"sync"
	"time"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "get_ip",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run:   getIp,
}

func getIp(cmd *cobra.Command, args []string) {
	db := model.DB.New()
	defer db.Close()
	pool := component.NewTaskPool(20)
	for i := 1; i < 10; i++ {
		contentBody := service.GetProxyService.GetContentHtml(i)
		if contentBody == nil {
			time.Sleep(time.Second * 5)
			continue
		}
		proxy_list := service.GetProxyService.ParseHtml(contentBody)
		logger.Info("获取到ip:", proxy_list)
		var wg sync.WaitGroup = sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			logger.Info("wait 20s ...")
			time.Sleep(time.Second * 20)
		}(&wg)
		for _, proxy_arr := range proxy_list {
			ip, port := proxy_arr[0], proxy_arr[1]
			pool.RunTask(func() { service.SaveProxyService.CheckProxyAndSave(ip, port, db) })
		}

		wg.Wait()
	}
	time.Sleep(time.Second * 20)
}
