package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tongsq/gin-example/component"
	"github.com/tongsq/gin-example/component/logger"
	"github.com/tongsq/gin-example/component/proxy"
	_ "github.com/tongsq/gin-example/init"
	"github.com/tongsq/gin-example/model"
	"github.com/tongsq/gin-example/service"
	"sync"
	"time"
)

func main() {
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
			pool.RunTask(func() { checkProxyAndSave(ip, port, db) })
		}

		wg.Wait()
	}
	time.Sleep(time.Second * 20)
}

func checkProxyAndSave(host string, port string, db *gorm.DB) {
	result := proxy.CheckIpStatus(host, port)
	if result {
		logger.Success(result, host, port)
	} else {
		logger.Warning(result, host, port)
	}
	var status int8 = 1
	if !result {
		status = 0
		return
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
