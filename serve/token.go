package serve

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"frp-auth/constant"
	"frp-auth/model"
	"frp-auth/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func tokenCheckPermission(c *gin.Context, tokenName string, username string) (model.TokenTable, error) {
	var tokenRow model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, tokenName).Scan(&tokenRow.ID, &tokenRow.Name, &tokenRow.Token, &tokenRow.User, &tokenRow.Enable); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return tokenRow, err
	}
	if tokenRow.User != username {
		sendJson(c, 403, "没有权限操作这条token", nil)
		return tokenRow, errors.New("no permission")
	}
	return tokenRow, nil
}

func createToken(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _createToken struct {
		Name string `json:"name"`
	}
	var req _createToken
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Name == "" {
		sendJson(c, 400, "名称不能为空", nil)
		return
	}
	var tokenRow model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, req.Name).Scan(&tokenRow.ID, &tokenRow.Name, &tokenRow.Token, &tokenRow.User, &tokenRow.Enable); err == nil {
		sendJson(c, 400, "该名称已存在", nil)
		return
	}
	_token, err := utils.GenerateRandomString(9)
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _, err := db.Exec(constant.InsertTokenTable, req.Name, _token, row.Username); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "Token创建成功", responseData{
		"name":  req.Name,
		"token": _token,
	})
}

func listToken(c *gin.Context) {
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
	var tokens []model.TokenTable
	for rows.Next() {
		var _t model.TokenTable
		if err := rows.Scan(&_t.ID, &_t.Name, &_t.Token, &_t.User, &_t.Enable); err != nil {
			sendJson(c, 500, err.Error(), nil)
			return
		}
		tokens = append(tokens, _t)
	}
	if err := rows.Err(); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if tokens == nil {
		tokens = make([]model.TokenTable, 0)
	}
	sendJson(c, 200, "读取成功", responseData{
		"tokens": tokens,
	})
}

func updateToken(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _updateToken struct {
		Name   string `json:"name"`
		Enable bool   `json:"enable"`
	}
	var req _updateToken
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Name == "" {
		sendJson(c, 400, "名称不能为空", nil)
		return
	}
	tokenRow, err := tokenCheckPermission(c, req.Name, row.Username)
	if err != nil {
		return
	}
	if _, err := db.Exec(constant.UpdateTokenTableEnableById, req.Enable, tokenRow.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "Token状态修改成功", responseData{
		"id":     tokenRow.ID,
		"name":   tokenRow.Name,
		"token":  tokenRow.Token,
		"user":   tokenRow.User,
		"enable": req.Enable,
	})
}

func deleteToken(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	name := c.Query("name")
	if name == "" {
		sendJson(c, 400, "名称不能为空", nil)
		return
	}
	tokenRow, err := tokenCheckPermission(c, name, row.Username)
	if err != nil {
		return
	}
	var port model.PortTable
	if err := db.QueryRow(constant.SelectPortTableByToken, name).Scan(&port.ID, &port.Min, &port.Max, &port.Token, &port.User); err == nil {
		sendJson(c, 400, "该Token有端口规则未删除", responseData{
			"port": port,
		})
		return
	}
	if _, err := db.Exec(constant.DeleteTokenTableRowById, tokenRow.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "删除成功", nil)
}

func regenerateToken(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _regenerateToken struct {
		Name string `json:"name"`
	}
	var req _regenerateToken
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Name == "" {
		sendJson(c, 400, "名称不能为空", nil)
		return
	}
	tokenRow, err := tokenCheckPermission(c, req.Name, row.Username)
	if err != nil {
		return
	}
	_token, err := utils.GenerateRandomString(9)
	if _, err := db.Exec(constant.UpdateTokenTableTokenById, _token, tokenRow.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "Token重新生成成功", responseData{
		"id":     tokenRow.ID,
		"name":   tokenRow.Name,
		"token":  _token,
		"user":   tokenRow.User,
		"enable": tokenRow.Enable,
	})
}

func generateConfig(c *gin.Context) {
	frpsHost := cfg.Frps.Host
	frpsPort := cfg.Frps.Port
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _generateConfig struct {
		ID int `uri:"id"`
	}
	var req _generateConfig
	if err := c.ShouldBindUri(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	var token model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableById, req.ID).Scan(&token.ID, &token.Name, &token.Token, &token.User, &token.Enable); err != nil {
		sendJson(c, 400, "Token不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if token.User != row.Username {
		sendJson(c, 403, "没有权限操作这个Token", nil)
		return
	}
	sendJson(c, 200, "配置生成成功", responseData{
		"config": fmt.Sprintf("serverAddr = \"%s\"\n"+
			"serverPort = %d\n"+
			"user = \"%s\"\n"+
			"metadatas.token = \"%s\"", frpsHost, frpsPort, token.Name, token.Token),
	})
}

func InitTokenApi(group *gin.RouterGroup) {
	token := group.Group("/token")
	token.POST("/", createToken)
	token.GET("/", listToken)
	token.GET("/conf/:id", generateConfig)
	token.PUT("/", updateToken)
	token.DELETE("/", deleteToken)
	token.PATCH("/", regenerateToken)
}
