package serve

import (
	"encoding/base64"
	"fmt"
	"frp-auth/constant"
	"frp-auth/model"
	"frp-auth/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseToken(c *gin.Context) (*model.JwtToken, model.UserTable) {
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	data, err := utils.JwtUserCheck(token)
	if err != nil {
		sendJson(c, 401, err.Error(), nil)
		return nil, model.UserTable{}
	}
	if data == nil {
		sendJson(c, 401, "数据异常", nil)
		return nil, model.UserTable{}
	}
	var row model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, data.Username).Scan(&row.ID, &row.Username, &row.Password, &row.Nickname, &row.Admin, &row.Enable, &row.TokenVersion); err != nil {
		sendJson(c, 401, err.Error(), nil)
		return nil, model.UserTable{}
	}
	if row.ID != data.UserID {
		sendJson(c, 401, "用户ID错误", nil)
		return nil, model.UserTable{}
	}
	if row.TokenVersion != data.TokenVersion {
		sendJson(c, 401, "Token无效", nil)
		return nil, model.UserTable{}
	}
	return data, row
}

func userLoginPubKey(c *gin.Context) {
	type _userLoginPubKey struct {
		Username string `json:"username"`
	}
	var req _userLoginPubKey
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	priKey, pubKey, err := utils.GenerateRsaKeyPair()
	if err != nil {
		return
	}
	RsaKeyStorage.Store("login_"+req.Username+"_private_key", priKey)
	sendJson(c, 200, "公钥生成成功", responseData{
		"public_key": pubKey,
	})
}

func userLogin(c *gin.Context) {
	var req model.UserLoginRequest
	var row model.UserTable
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Username == "" || req.Password == "" {
		sendJson(c, 400, "用户名或密码不能为空", nil)
		return
	}
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.Username).Scan(&row.ID, &row.Username, &row.Password, &row.Nickname, &row.Admin, &row.Enable, &row.TokenVersion); err != nil {
		sendJson(c, 401, "用户不存在", nil)
		return
	}
	decodeString, err := base64.StdEncoding.DecodeString(req.Password)
	if err != nil {
		sendJson(c, 500, "Base64解析失败", responseData{
			"err": err.Error(),
		})
		return
	}
	priKey, ok := RsaKeyStorage.LoadAndDelete("login_" + req.Username + "_private_key")
	if !ok {
		sendJson(c, 400, "私钥不存在", nil)
		return
	}
	priKeyPem, ok := priKey.(string)
	if !ok {
		sendJson(c, 400, "私钥解析失败", nil)
		return
	}
	data, err := utils.DecryptWithPrivateKey(decodeString, priKeyPem)
	if err != nil {
		sendJson(c, 400, "数据解密失败", responseData{
			"err": err.Error(),
		})
		return
	}
	if utils.PasswordCheck(row.Password, string(data)) {
		if row.Enable {
			tokenVersion, err := utils.GenerateRandomString(16)
			if err != nil {
				sendJson(c, 500, err.Error(), nil)
				return
			}
			token, err := utils.JwtUserGenerate(row.ID, row.Username, row.Admin, row.Enable, tokenVersion)
			if err != nil {
				sendJson(c, 500, err.Error(), nil)
				return
			}
			if _, err := db.Exec(constant.UpdateUserTableTokenVersionById, tokenVersion, row.ID); err != nil {
				sendJson(c, 500, err.Error(), nil)
				return
			}
			sendJson(c, 200, "登录成功", responseData{
				"token": token,
			})
		} else {
			sendJson(c, 401, "账号未启用", nil)
		}
	} else {
		sendJson(c, 401, "密码错误", nil)
	}
}

func userRegisterPubKey(c *gin.Context) {
	type _userLoginPubKey struct {
		Username string `json:"username"`
	}
	var req _userLoginPubKey
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	priKey, pubKey, err := utils.GenerateRsaKeyPair()
	if err != nil {
		return
	}
	RsaKeyStorage.Store("register_"+req.Username+"_private_key", priKey)
	sendJson(c, 200, "公钥生成成功", responseData{
		"public_key": pubKey,
	})
}

func userRegister(c *gin.Context) {
	type _userRegister struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	var req _userRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if req.Username == "" || req.Password == "" {
		sendJson(c, 400, "用户名或密码不能为空", nil)
		return
	}
	var row model.UserTable
	if err := db.QueryRow(constant.SelectUserTableByUsername, req.Username).Scan(&row.ID, &row.Username, &row.Password, &row.Nickname, &row.Admin, &row.Enable, &row.TokenVersion); err == nil {
		sendJson(c, 400, "用户名已存在", nil)
		return
	}
	decodeString, err := base64.StdEncoding.DecodeString(req.Password)
	if err != nil {
		sendJson(c, 500, "Base64解析失败", responseData{
			"err": err.Error(),
		})
		return
	}
	priKey, ok := RsaKeyStorage.LoadAndDelete("register_" + req.Username + "_private_key")
	if !ok {
		sendJson(c, 400, "私钥不存在", nil)
		return
	}
	priKeyPem, ok := priKey.(string)
	if !ok {
		sendJson(c, 400, "私钥解析失败", nil)
		return
	}
	data, err := utils.DecryptWithPrivateKey(decodeString, priKeyPem)
	if err != nil {
		sendJson(c, 400, "数据解密失败", responseData{
			"err": err.Error(),
		})
		return
	}
	password, err := utils.PasswordHash(string(data))
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _, err := db.Exec(constant.InsertUserTable, req.Username, password, req.Nickname, false, false); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "注册成功，请等待管理员确认", nil)
}

func userInfo(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	if row.Enable {
		sendJson(c, 200, "数据获取成功", responseData{
			"userId":   strconv.Itoa(row.ID),
			"username": row.Username,
			"nickname": row.Nickname,
			"role": func() string {
				if row.Admin {
					return "admin"
				}
				return "user"
			}(),
		})
		return
	}
	sendJson(c, 401, "用户未启用", nil)
	return
}

func userUpdate(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _userUpdate struct {
		Nickname string `json:"nickname"`
	}
	var req _userUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	fmt.Println(req.Nickname)
	if _, err := db.Exec(constant.UpdateUserTableNicknameById, req.Nickname, row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "昵称修改成功", nil)
	return
}
func userChangePwd(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	type _userChangePwd struct {
		Password string `json:"password"`
	}
	var req _userChangePwd
	if err := c.ShouldBindJSON(&req); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	passwordHash, err := utils.PasswordHash(req.Password)
	if err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _, err := db.Exec(constant.UpdateUserTablePasswordById, passwordHash, row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	if _, err := db.Exec(constant.UpdateUserTableTokenVersionById, "", row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "密码修改成功", nil)
}
func userLogout(c *gin.Context) {
	data, row := parseToken(c)
	if data == nil {
		return
	}
	if _, err := db.Exec(constant.UpdateUserTableTokenVersionById, "", row.ID); err != nil {
		sendJson(c, 500, err.Error(), nil)
		return
	}
	sendJson(c, 200, "成功退出登录", nil)
}

func InitUserApi(group *gin.RouterGroup) {
	user := group.Group("/user")
	user.POST("/login", userLogin)
	user.PUT("/login", userLoginPubKey)
	user.POST("/register", userRegister)
	user.PUT("/register", userRegisterPubKey)
	user.GET("/", userInfo)
	user.PUT("/", userUpdate)
	user.PATCH("/", userChangePwd)
	user.POST("/logout", userLogout)
}
