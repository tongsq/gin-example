package service

import proxy_service "github.com/tongsq/gin-example/service/proxy-service"

var (
	GetProxyService  proxy_service.GetProxyInterface
	SaveProxyService *proxy_service.SaveProxyService
)

func init() {
	GetProxyService = &proxy_service.GetProxyXila{}
	SaveProxyService = &proxy_service.SaveProxyService{}
}
