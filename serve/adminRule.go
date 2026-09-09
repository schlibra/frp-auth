package serve

import (
	"context"
	"database/sql"
	"frp-auth/constant"
	"frp-auth/model"
	"time"

	"github.com/gin-gonic/gin"
)

func adminListPort(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, constant.SelectPortTable)
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
	var ports []model.PortTable
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
func adminCreatePort(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminCreatePort struct {
		Min  int    `json:"min"`
		Max  int    `json:"max"`
		Name string `json:"name"`
		User string `json:"user"`
	}
	var req _adminCreatePort
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
	var userRow model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.User).Scan(&userRow.ID, &userRow.Username, &userRow.Password, &userRow.Nickname, &userRow.Admin, &userRow.Enable, &userRow.TokenVersion); err != nil {
		sendJson(c, 400, "用户不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if _, err := db.Exec(constant.InsertPortTable, req.Min, req.Max, req.Name, req.User); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "创建成功", nil)
}
func adminUpdatePort(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminCreatePort struct {
		ID   int    `json:"id"`
		Min  int    `json:"min"`
		Max  int    `json:"max"`
		Name string `json:"name"`
		User string `json:"user"`
	}
	var req _adminCreatePort
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
		sendJson(c, 400, "端口规则记录不存在", nil)
		return
	}
	var token model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, req.Name).Scan(&token.ID, &token.Name, &token.Token, &token.User, &token.Enable); err != nil {
		sendJson(c, 400, "Token不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	var userRow model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.User).Scan(&userRow.ID, &userRow.Username, &userRow.Password, &userRow.Nickname, &userRow.Admin, &userRow.Enable, &userRow.TokenVersion); err != nil {
		sendJson(c, 400, "用户不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if _, err := db.Exec(constant.UpdatePortTableById, req.Min, req.Max, req.Name, req.User, req.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "更新成功", responseData{
		"min":   req.Min,
		"max":   req.Max,
		"token": req.Name,
		"user":  req.User,
	})
}
func adminDeletePort(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminDeletePort struct {
		ID int `uri:"id"`
	}
	var req _adminDeletePort
	if err := c.ShouldBindUri(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	var port model.PortTable
	if err := db.QueryRow(constant.SelectPortTableById, req.ID).Scan(&port.ID, &port.Min, &port.Min, &port.Token, &port.User); err != nil {
		sendJson(c, 400, "端口规则记录不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if _, err := db.Exec(constant.DeletePortTableById, req.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "删除成功", nil)
}

func InitAdminRuleApi(admin *gin.RouterGroup) {
	port := admin.Group("/rule")
	port.GET("/", adminListPort)
	port.POST("/", adminCreatePort)
	port.PUT("/", adminUpdatePort)
	port.DELETE("/:id", adminDeletePort)
}
