package service

import proxy_service "github.com/tongsq/gin-example/service/proxy-service"

var (
	GetProxyService proxy_service.GetProxyInterface
)

func init() {
	GetProxyService = &proxy_service.GetProxyNima{}
}
