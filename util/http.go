package util

import (
	"net/http"
	"strings"
)

func HttpSchemeToPort(scheme string) int {
	if strings.ToLower(scheme) == "http" {
		return 80
	} else if strings.ToLower(scheme) == "https" {
		return 443
	}
	return 0
}

func GetFristRequest(resp *http.Response) *http.Request {
	if resp == nil {
		return nil
	}

	// 如果响应没有请求，说明这是第一个请求
	if resp.Request == nil {
		return nil
	}

	if resp.Request.Response == nil {
		return resp.Request
	}

	return GetFristRequest(resp.Request.Response)
}
