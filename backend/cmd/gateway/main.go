package main

import (
	"log"
	"net/http"
	"net/http/proxy"
	"net/url"
	"strings"

	"github.com/fileshare/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Name string
	URL  *url.URL
}

var services = map[string]*Service{
	"user":   {Name: "user-service", URL: mustParseURL("http://localhost:8081")},
	"file":   {Name: "file-service", URL: mustParseURL("http://localhost:8082")},
	"share":  {Name: "share-service", URL: mustParseURL("http://localhost:8083")},
}

var whiteList = []string{
	"/api/user/register",
	"/api/user/login",
	"/share/",
	"/health",
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

func main() {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 代理路由
	r.Any("/api/*path", handleProxy)
	r.GET("/share/*path", handleShareProxy)

	log.Println("网关服务启动在 :8000")
	r.Run(":8000")
}

func handleProxy(c *gin.Context) {
	path := c.Param("path")
	segments := strings.Split(strings.Trim(path, "/"), "/")

	if len(segments) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}

	serviceName := segments[0]
	service, ok := services[serviceName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	// 白名单检查
	fullPath := "/api/" + strings.Trim(path, "/")
	for _, wp := range whiteList {
		if strings.HasPrefix(fullPath, wp) {
			proxyRequest(c, service.URL, path)
			return
		}
	}

	// JWT 验证
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	token := authHeader[7:]
	claims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 设置用户信息到请求头
	c.Request.Header.Set("X-User-Id", string(rune(claims.UserId)))
	c.Request.Header.Set("X-Username", claims.Username)
	c.Request.Header.Set("X-Role", claims.Role)

	proxyRequest(c, service.URL, path)
}

func handleShareProxy(c *gin.Context) {
	path := c.Param("path")
	proxyRequest(c, services["share"].URL, "/share"+path)
}

func proxyRequest(c *gin.Context, target *url.URL, path string) {
	target.Path = path
	target.RawQuery = c.Request.URL.RawQuery

	// 创建代理
	director := func(req *http.Request) {
		req.URL = target
		req.Host = target.Host
	}

	p := &proxy.ReverseProxy{Director: director}
	p.ServeHTTP(c.Writer, c.Request)
}
