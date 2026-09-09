package serve

import (
	"frp-auth/model"

	"github.com/gin-gonic/gin"
)

func checkAdmin(c *gin.Context) (*model.JwtToken, model.UserTable) {
	data, row := parseToken(c)
	if data == nil {
		return data, row
	}
	if !row.Admin {
		sendJson(c, 403, "没有管理员权限", nil)
		return nil, row
	}
	return data, row
}

func InitAdminApi(group *gin.RouterGroup) {
	admin := group.Group("/admin")
	InitAdminUserApi(admin)
	InitAdminTokenApi(admin)
	InitAdminRuleApi(admin)
	InitAdminPortApi(admin)
}
