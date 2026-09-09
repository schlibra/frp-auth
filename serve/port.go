package serve

import (
	"context"
	"database/sql"
	"fmt"
	"frp-auth/constant"
	"frp-auth/model"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func getPortList(c *gin.Context) {
	frpsHost := cfg.FrpsWeb.Host
	frpsPort := cfg.FrpsWeb.Port
	frpsToken := cfg.FrpsWeb.Token
	data, row := parseToken(c)
	if data == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, constant.SelectTokenTableByUser, row.Username)
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			sendJson(c, 500, err.Error(), nil)
			return
		}
	}(rows)
	var tokens []string
	for rows.Next() {
		var _t model.TokenTable
		if err := rows.Scan(&_t.ID, &_t.Name, &_t.Token, &_t.User, &_t.Enable); err != nil {
			sendJson(c, 500, err.Error(), nil)
			return
		}
		tokens = append(tokens, _t.Name)
	}
	client := resty.New()
	var clientResult model.FrpsClient
	_, err = client.R().
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
	clients := make([]model.FrpsClientDataItem, 0)
	for _, item := range clientResult.Data.Items {
		if slices.Contains(tokens, item.User) {
			clients = append(clients, item)
		}
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
	proxies := make([]model.FrpsProxyDataItem, 0)
	for _, item := range proxyResult.Data.Items {
		if slices.Contains(tokens, item.User) {
			proxies = append(proxies, item)
		}
	}
	sendJson(c, 200, "数据获取成功", responseData{
		"clients": clients,
		"proxies": proxies,
	})
}

func InitPortApi(group *gin.RouterGroup) {
	port := group.Group("/port")
	port.GET("/", getPortList)
}
