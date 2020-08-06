package proxy_service

import "io"

type GetProxyInterface interface {
	GetContentHtml(i int) io.ReadCloser
	ParseHtml(body io.ReadCloser) [][]string
}
