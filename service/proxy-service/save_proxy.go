package proxy_service

import (
	"github.com/jinzhu/gorm"
	"github.com/tongsq/gin-example/component/logger"
	"github.com/tongsq/gin-example/component/proxy"
	"github.com/tongsq/gin-example/model"
	"time"
)

type SaveProxyService struct {
}

func (s *SaveProxyService) CheckProxyAndSave(host string, port string, db *gorm.DB) {
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
