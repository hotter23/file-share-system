package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fileshare/internal/dto"
	"github.com/fileshare/internal/model"
	"github.com/fileshare/pkg/jwt"
	"github.com/fileshare/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	dsn := "root:root123@tcp(localhost:3306)/file_share?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id, X-Username, X-Role")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 路由
	r.POST("/api/user/register", register)
	r.POST("/api/user/login", login)
	r.GET("/api/user/info", authMiddleware(), getUserInfo)

	log.Println("用户服务启动在 :8081")
	r.Run(":8081")
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader("X-User-Id")
		if userId == "" {
			response.Unauthorized(c, "未认证")
			c.Abort()
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}

func register(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 检查用户名是否存在
	var count int64
	db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		response.BadRequest(c, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.ServerError(c, "密码加密失败")
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "USER",
		Status:   "ACTIVE",
	}

	if err := db.Create(&user).Error; err != nil {
		response.ServerError(c, "创建用户失败")
		return
	}

	response.SuccessWithMessage(c, "注册成功", gin.H{
		"userId":   user.ID,
		"username": user.Username,
	})
}

func login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	if user.Status == "DISABLED" {
		response.Forbidden(c, "账号已被禁用")
		return
	}

	token, err := jwt.GenerateToken(int64(user.ID), user.Username, user.Role)
	if err != nil {
		response.ServerError(c, "生成Token失败")
		return
	}

	response.Success(c, gin.H{
		"token":    token,
		"userId":   user.ID,
		"username": user.Username,
	})
}

func getUserInfo(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userId, _ := strconv.ParseUint(userIdStr.(string), 10, 64)

	var user model.User
	if err := db.First(&user, userId).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"userId":   user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}
