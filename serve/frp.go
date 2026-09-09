package serve

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"frp-auth/constant"
	"frp-auth/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func frpLoginApi(c *gin.Context) {
	var req model.RequestBody
	var row model.TokenTable
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"reject":        true,
			"reject_reason": err.Error(),
		})
		return
	}
	user := req.Content.User
	token := req.Content.Metas.Token
	err := db.QueryRow(constant.SelectTokenTableByName, user).Scan(&row.ID, &row.Name, &row.Token, &row.User, &row.Enable)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, gin.H{
				"reject":        true,
				"reject_reason": "用户不存在|User not exist",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"reject":        true,
			"reject_reason": "数据查询失败|Data query error",
		})
		return
	}
	if row.Token == token {
		c.JSON(http.StatusOK, gin.H{
			"reject":   false,
			"unchange": true,
		})
		return
	}
	if row.Enable == 0 {
		c.JSON(http.StatusOK, gin.H{
			"reject":        true,
			"reject_reason": "Token未启用|Token not enable",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"reject":        true,
		"reject_reason": "令牌错误|Token not match",
	})
	return
}

func frpProxyApi(c *gin.Context) {
	type _frpProxyApiContentUser struct {
		User string `json:"user"`
	}
	type _frpProxyApiContent struct {
		User       _frpProxyApiContentUser `json:"user"`
		RemotePort int                     `json:"remote_port"`
	}
	type _frpProxyApi struct {
		Content _frpProxyApiContent `json:"content"`
	}
	var req _frpProxyApi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"reject":        true,
			"reject_reason": err.Error(),
		})
		return
	}
	user := req.Content.User.User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := db.QueryContext(ctx, constant.SelectPortTableByToken, user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"reject":        true,
			"reject_reason": err.Error(),
		})
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"reject":        true,
				"reject_reason": err.Error(),
			})
			return
		}
	}(rows)
	allowPort := make([]string, 0)
	remotePort := req.Content.RemotePort
	for rows.Next() {
		var port model.PortTable
		if err := rows.Scan(&port.ID, &port.Min, &port.Max, &port.Token, &port.User); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"reject":        true,
				"reject_reason": err.Error(),
			})
			return
		}
		if remotePort >= port.Min && remotePort <= port.Max {
			c.JSON(http.StatusOK, gin.H{
				"reject":   false,
				"unchange": true,
			})
			return
		}
		allowPort = append(allowPort, strconv.Itoa(port.Min)+func() string {
			if port.Min == port.Max {
				return ""
			}
			return "-" + strconv.Itoa(port.Max)
		}())
	}
	ports := strings.Join(allowPort, ", ")
	c.JSON(http.StatusOK, gin.H{
		"reject":        true,
		"reject_reason": fmt.Sprintf("用户%s不允许使用端口%d，只允许使用端口：%s|User %s not allow use port %d, only allow use port: %s", user, remotePort, ports, user, remotePort, ports),
	})
}

func InitFrpApi(group *gin.RouterGroup) {
	frp := group.Group("/frp")
	frp.POST("/login", frpLoginApi)
	frp.POST("/proxy", frpProxyApi)
}
