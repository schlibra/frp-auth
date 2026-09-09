package serve

import (
	"context"
	"database/sql"
	"fmt"
	"frp-auth/constant"
	"frp-auth/model"
	"frp-auth/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func adminListUser(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userRow, err := db.QueryContext(ctx, constant.SelectUserTable)
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	defer func(userRow *sql.Rows) {
		err := userRow.Close()
		if err != nil {
			sendJson(c, 500, err.Error(), nil)
			return
		}
	}(userRow)
	var users []model.UserTable
	for userRow.Next() {
		var _u model.UserTable
		if err := userRow.Scan(&_u.ID, &_u.Username, &_u.Password, &_u.Nickname, &_u.Admin, &_u.Enable, &_u.TokenVersion); err != nil {
			sendJson(c, 500, err.Error(), nil)
			return
		}
		_u.Password = ""
		_u.TokenVersion = ""
		users = append(users, _u)
	}
	sendJson(c, 200, "数据获取成功", responseData{
		"users": users,
	})
}
func adminUpdateUser(c *gin.Context) {
	data, _row := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminUpdate struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Admin    bool   `json:"admin"`
		Enable   bool   `json:"enable"`
	}
	var req _adminUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _row.Username == req.Username {
		if !req.Admin || !req.Enable {
			sendJson(c, 400, "不可以修改该用户的Admin或Enable", nil)
			return
		}
	}
	var row model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.Username).Scan(&row.ID, &row.Username, &row.Password, &row.Nickname, &row.Admin, &row.Enable, &row.TokenVersion); err != nil {
		sendJson(c, 400, "用户不存在", nil)
		return
	}
	if req.Nickname != "" {
		row.Nickname = req.Nickname
	}
	row.Admin = req.Admin
	row.Enable = req.Enable
	if req.Password != "" {
		_p, err := utils.PasswordHash(req.Password)
		if err != nil {
			sendJson(c, 500, err.Error(), nil)
			return
		}
		row.Password = _p
	}
	fmt.Println(row)
	if _, err := db.Exec(constant.UpdateUserTableById, row.Password, row.Nickname, row.Admin, row.Enable, row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	row.Password = ""
	row.TokenVersion = ""
	sendJson(c, 200, "用户更新成功", responseData{
		"user": row,
	})
}
func adminCreateUser(c *gin.Context) {
	data, _ := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminCreateUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Admin    bool   `json:"admin"`
		Enable   bool   `json:"enable"`
	}
	var req _adminCreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	var row model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.Username).Scan(&row.ID, &row.Username, &row.Password, &row.Nickname, &row.Admin, &row.Enable, &row.TokenVersion); err == nil {
		sendJson(c, 400, "用户已存在", nil)
		return
	}
	_password, err := utils.PasswordHash(req.Password)
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _, err := db.Exec(constant.InsertUserTable, req.Username, _password, req.Nickname, req.Admin, req.Enable); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "用户创建成功", nil)
}

func adminDeleteUser(c *gin.Context) {
	data, _row := checkAdmin(c)
	if data == nil {
		return
	}
	type _adminDeleteUser struct {
		ID string `uri:"id"`
	}
	var req _adminDeleteUser
	if err := c.ShouldBindUri(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	var row model.UserTable
	if err := db.QueryRow(constant.SelectUserTableById, req.ID).Scan(&row.ID, &row.Username, &row.Password, &row.Nickname, &row.Admin, &row.Enable, &row.TokenVersion); err != nil {
		sendJson(c, 400, "用户不存在", nil)
		return
	}
	if _row.Username == row.Username {
		sendJson(c, 400, "不能删除该用户", nil)
		return
	}
	if _, err := db.Exec(constant.DeleteUserTableById, row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "删除用户成功", nil)
}

func InitAdminUserApi(admin *gin.RouterGroup) {
	group := admin.Group("/user")
	group.GET("/", adminListUser)
	group.PUT("/", adminUpdateUser)
	group.POST("/", adminCreateUser)
	group.DELETE("/:id", adminDeleteUser)
}
