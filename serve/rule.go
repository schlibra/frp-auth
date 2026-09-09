package serve

import (
	"context"
	"database/sql"
	"frp-auth/constant"
	"frp-auth/model"
	"time"

	"github.com/gin-gonic/gin"
)

func listPort(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, constant.SelectPortTableByUser, row.Username)
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
	var ports = make([]model.PortTable, 0)
	for rows.Next() {
		var _p model.PortTable
		if err := rows.Scan(&_p.ID, &_p.Min, &_p.Max, &_p.Token, &_p.User); err != nil {
			sendJson(c, 500, err.Error(), nil)
			return
		}
		ports = append(ports, _p)
	}
	sendJson(c, 200, "数据获取成功", responseData{
		"ports": ports,
	})
}
func createPort(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _createPort struct {
		Name string `json:"name"`
		Min  int    `json:"min"`
		Max  int    `json:"max"`
	}
	var req _createPort
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Max <= 0 || req.Max >= 65536 || req.Min <= 0 || req.Min >= 65536 {
		sendJson(c, 400, "端口必须为1-65535", responseData{
			"min": req.Min,
			"max": req.Max,
		})
		return
	}
	if req.Min > req.Max {
		sendJson(c, 400, "min端口不能大于max端口", nil)
		return
	}
	var token model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, req.Name).Scan(&token.ID, &token.Name, &token.Token, &token.User, &token.Enable); err != nil {
		sendJson(c, 400, "Token不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if token.User != row.Username {
		sendJson(c, 403, "没有权限为该Token创建端口规则", nil)
		return
	}
	if _, err := db.Exec(constant.InsertPortTable, req.Min, req.Max, token.Name, row.Username); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "端口规则创建成功", nil)

}

func updatePort(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _updatePort struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Min  int    `json:"min"`
		Max  int    `json:"max"`
	}
	var req _updatePort
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Max <= 0 || req.Max >= 65536 || req.Min <= 0 || req.Min >= 65536 {
		sendJson(c, 400, "端口必须为1-65535", responseData{
			"min": req.Min,
			"max": req.Max,
		})
		return
	}
	if req.Min > req.Max {
		sendJson(c, 400, "min端口不能大于max端口", nil)
		return
	}
	var port model.PortTable
	if err := db.QueryRow(constant.SelectPortTableById, req.ID).Scan(&port.ID, &port.Min, &port.Max, &port.Token, &port.User); err != nil {
		sendJson(c, 400, "端口规则记录不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	var token model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, port.Token).Scan(&token.ID, &token.Name, &token.Token, &token.User, &token.Enable); err != nil {
		sendJson(c, 400, "Token不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if token.User != row.Username || port.User != row.Username {
		sendJson(c, 403, "没有权限操作这条记录", nil)
		return
	}
	if err := db.QueryRow(constant.SelectTokenTableByName, req.Name).Scan(&token.ID, &token.Name, &token.Token, &token.User, &token.Enable); err != nil {
		sendJson(c, 400, "Token不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if token.User != row.Username {
		sendJson(c, 403, "没有权限绑定到这个Token", nil)
		return
	}
	if _, err := db.Exec(constant.UpdatePortTableById, req.Min, req.Max, req.Name, row.Username, req.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "端口规则修改成功", responseData{
		"id":    req.ID,
		"min":   req.Min,
		"max":   req.Max,
		"token": token,
		"user":  row.Username,
	})
}

func deletePort(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _deletePort struct {
		ID int `uri:"id"`
	}
	var req _deletePort
	if err := c.ShouldBindUri(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	var port model.PortTable
	if err := db.QueryRow(constant.SelectPortTableById, req.ID).Scan(&port.ID, &port.Min, &port.Max, &port.Token, &port.User); err != nil {
		sendJson(c, 400, "记录不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	var token model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, port.Token).Scan(&token.ID, &token.Name, &token.Token, &token.User, &token.Enable); err != nil {
		sendJson(c, 400, "Token不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if port.User != row.Username || token.User != row.Username {
		sendJson(c, 403, "没有权限操作", nil)
		return
	}
	if _, err := db.Exec(constant.DeletePortTableById, req.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "端口规则删除成功", nil)
}

func InitRuleApi(group *gin.RouterGroup) {
	port := group.Group("/rule")
	port.GET("/", listPort)
	port.POST("/", createPort)
	port.PUT("/", updatePort)
	port.DELETE("/:id", deletePort)
}
