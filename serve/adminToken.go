package serve

import (
	"context"
	"database/sql"
	"frp-auth/constant"
	"frp-auth/model"
	"frp-auth/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func adminListToken(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, constant.SelectTokenTable)
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
	sendJson(c, 200, "数据获取成功", responseData{
		"tokens": tokens,
	})
}
func adminUpdateToken(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminUpdateToken struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Enable   bool   `json:"enable"`
	}
	var req _adminUpdateToken
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	var row model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, req.Name).Scan(&row.ID, &row.Name, &row.Token, &row.User, &row.Enable); err != nil {
		sendJson(c, 400, "Token不存在", nil)
		return
	}
	var userRow model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.Username).Scan(&userRow.ID, &userRow.Username, &userRow.Password, &userRow.Nickname, &userRow.Admin, &userRow.Enable, &userRow.TokenVersion); err != nil {
		sendJson(c, 400, "用户不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	if _, err := db.Exec(constant.UpdateTokenTableById, req.Enable, req.Username, row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	row.Enable = func() int {
		if req.Enable {
			return 1
		}
		return 0
	}()
	sendJson(c, 200, "Token更新成功", responseData{
		"token": row,
	})
}
func adminRegenerateToken(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminRegenerateToken struct {
		Name string `json:"name"`
	}
	var req _adminRegenerateToken
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	var row model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, req.Name).Scan(&row.ID, &row.Name, &row.Token, &row.User, &row.Enable); err != nil {
		sendJson(c, 400, "Token不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	token, err := utils.GenerateRandomString(9)
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _, err := db.Exec(constant.UpdateTokenTableTokenById, token, row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "Token重新生成成功", responseData{
		"token": token,
	})
}
func adminCreateToken(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminCreateToken struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}
	var req _adminCreateToken
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Name == "" {
		sendJson(c, 400, "名称不能为空", nil)
		return
	}
	var row model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, req.Name).Scan(&row.ID, &row.Name, &row.Token, &row.User, &row.Enable); err == nil {
		sendJson(c, 400, "Token已存在", nil)
		return
	}
	var userRow model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.Username).Scan(&userRow.ID, &userRow.Username, &userRow.Password, &userRow.Nickname, &userRow.Admin, &userRow.Enable, &userRow.TokenVersion); err != nil {
		sendJson(c, 400, "用户不存在", responseData{
			"err": err.Error(),
		})
		return
	}
	_token, err := utils.GenerateRandomString(9)
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _, err := db.Exec(constant.InsertTokenTable, req.Name, _token, req.Username); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "Token创建成功", responseData{
		"token": _token,
		"name":  req.Name,
		"user":  req.Username,
	})
}
func adminDeleteToken(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	name := c.Query("name")

	var row model.TokenTable
	if err := db.QueryRow(constant.SelectTokenTableByName, name).Scan(&row.ID, &row.Name, &row.Token, &row.User, &row.Enable); err != nil {
		sendJson(c, 400, "Token不存在", nil)
		return
	}
	if _, err := db.Exec(constant.DeleteTokenTableRowById, row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "删除Token成功", nil)
}

func InitAdminTokenApi(admin *gin.RouterGroup) {
	token := admin.Group("/token")
	token.GET("/", adminListToken)
	token.PUT("/", adminUpdateToken)
	token.POST("/", adminCreateToken)
	token.PATCH("/", adminRegenerateToken)
	token.DELETE("/", adminDeleteToken)
}
