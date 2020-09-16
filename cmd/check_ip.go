package cmd

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/tongsq/gin-example/component/proxy"
	"time"
)
import _ "github.com/jinzhu/gorm/dialects/mysql"
import "github.com/tongsq/gin-example/model"
import "github.com/tongsq/gin-example/component"

func init() {
	rootCmd.AddCommand(checkIpCmd)
}

var checkIpCmd = &cobra.Command{
	Use:   "check_ip",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run:   checkIpExist,
}

/*
 * 检查已存在的失效ip
 */
func checkIpExist(cmd *cobra.Command, args []string) {
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

func CheckProxyStatus(proxyModel model.Proxy, db *gorm.DB) {
	fmt.Printf("start check :host:%s, port:%s\n", proxyModel.Host, proxyModel.Port)
	result := proxy.CheckIpStatus(proxyModel.Host, proxyModel.Port)
	fmt.Printf("%s, %s, the result is %v\n", proxyModel.Host, proxyModel.Port, result)
	if result {
		proxyModel.Status = 1
		db.Save(&proxyModel)
	}
}
