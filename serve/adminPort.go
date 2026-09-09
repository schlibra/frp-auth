package serve

import (
	"fmt"
	"frp-auth/model"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func getAdminPortList(c *gin.Context) {
	frpsHost := cfg.FrpsWeb.Host
	frpsPort := cfg.FrpsWeb.Port
	frpsToken := cfg.FrpsWeb.Token
	data, _ := parseToken(c)
	if data == nil {
		return
	}

	client := resty.New()
	var clientResult model.FrpsClient
	_, err := client.R().
		SetAuthScheme("Basic").
		SetAuthToken(frpsToken).
		SetResult(&clientResult).
		Get(fmt.Sprintf("http://%s:%d/api/v2/clients?pageSize=200", frpsHost, frpsPort))
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	fmt.Printf("%v", clientResult)
	if clientResult.Code != 200 {
		sendJson(c, clientResult.Code, clientResult.Msg, nil)
		return
	}
	var proxyResult model.FrpsProxy
	_, err = client.R().
		SetAuthScheme("Basic").
		SetAuthToken(frpsToken).
		SetResult(&proxyResult).
		Get(fmt.Sprintf("http://%s:%d/api/v2/proxies?pageSize=200", frpsHost, frpsPort))
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if proxyResult.Code != 200 {
		sendJson(c, proxyResult.Code, proxyResult.Msg, nil)
		return
	}
	sendJson(c, 200, "数据获取成功", responseData{
		"clients": clientResult.Data.Items,
		"proxies": proxyResult.Data.Items,
	})
}

func InitAdminPortApi(group *gin.RouterGroup) {
	port := group.Group("/port")
	port.GET("/", getAdminPortList)
}
