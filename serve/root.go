package serve

import (
	"database/sql"
	"embed"
	"fmt"
	"frp-auth/model"
	"frp-auth/utils"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var cfg model.Config
var db *sql.DB

type responseData map[string]any

func sendJson(c *gin.Context, status int, msg string, data responseData) {
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"msg":    msg,
		"data":   data,
	})
}

func Serve(distFS embed.FS) {
	fmt.Println("Loading config...")
	cfg = utils.LoadConfig()
	fmt.Println("Loading database...")
	db = utils.LoadMySQL(cfg)
	host := cfg.Server.Host
	port := cfg.Server.Port
	debug := cfg.Server.Debug
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	api := r.Group("/api")
	InitFrpApi(api)
	InitUserApi(api)
	InitTokenApi(api)
	InitAdminApi(api)
	InitRuleApi(api)
	InitPortApi(api)

	subFS, err := fs.Sub(distFS, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/assets/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		println(filepath)

		data, err := fs.ReadFile(subFS, path.Join("assets", filepath))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		// 根据文件名简单判断 Content-Type
		contentType := func() string {
			if strings.HasSuffix(filepath, ".css") {
				return "text/css"
			}
			if strings.HasSuffix(filepath, ".js") {
				return "text/javascript"
			}
			return http.DetectContentType(data)
		}()

		c.Data(http.StatusOK, contentType, data)
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		data, err := fs.ReadFile(subFS, "favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.Data(http.StatusOK, "image/x-icon", data)
	})

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			sendJson(c, 404, "Api not found", nil)
			return
		}
		data, err := fs.ReadFile(subFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "index.html not found")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	if !debug {
		fmt.Printf("[GIN] Listening and serving HTTP on %s:%d\n", host, port)
	}
	err = r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}
}
