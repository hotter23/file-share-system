package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fileshare/internal/dto"
	"github.com/fileshare/internal/model"
	"github.com/fileshare/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// 公开路由（分享访问）
	r.GET("/share/:shareUuid", accessShare)
	r.GET("/share/:shareUuid/download", downloadShare)

	// 需要认证的路由
	api := r.Group("/api/share")
	api.Use(authMiddleware())
	{
		api.POST("", createShare)
		api.GET("/my", getMyShares)
		api.DELETE("/:shareUuid", deleteShare)
	}

	log.Println("分享服务启动在 :8083")
	r.Run(":8083")
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

func getUserId(c *gin.Context) uint {
	userIdStr, _ := c.Get("userId")
	var id uint
	// 处理接口
	switch v := userIdStr.(type) {
	case string:
		var parsed int64
		for _, ch := range v {
			if ch >= '0' && ch <= '9' {
				parsed = parsed*10 + int64(ch-'0')
			}
		}
		id = uint(parsed)
	case uint:
		id = v
	case int64:
		id = uint(v)
	}
	return id
}

func createShare(c *gin.Context) {
	userId := getUserId(c)

	var req dto.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取文件
	var file model.File
	if err := db.Where("file_uuid = ? AND user_id = ?", req.FileUuid, userId).First(&file).Error; err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	shareUuid := uuid.New().String()
	shareCode := generateShortCode()

	var expireTime *time.Time
	if req.ExpireType != "" && req.ExpireType != "PERMANENT" {
		expire := calculateExpireTime(req.ExpireType)
		expireTime = &expire
	}

	var hashedPassword *string
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		hashedPassword = new(string)
		*hashedPassword = string(hash)
	}

	share := model.Share{
		ShareUuid:  shareUuid,
		FileID:     file.ID,
		ShareCode:  shareCode,
		Password:   *hashedPassword,
		ExpireTime: expireTime,
		ViewCount:  0,
		CreatedBy:  userId,
	}

	if err := db.Create(&share).Error; err != nil {
		response.ServerError(c, "创建分享失败")
		return
	}

	response.Success(c, gin.H{
		"shareUuid":  shareUuid,
		"shareUrl":   "/share/" + shareUuid,
		"shareCode":  shareCode,
		"expireTime": expireTime,
	})
}

func accessShare(c *gin.Context) {
	shareUuid := c.Param("shareUuid")
	password := c.Query("password")

	var share model.Share
	if err := db.Where("share_uuid = ?", shareUuid).First(&share).Error; err != nil {
		response.NotFound(c, "分享不存在")
		return
	}

	// 检查过期
	if share.ExpireTime != nil && share.ExpireTime.Before(time.Now()) {
		response.Error(c, 410, "分享已过期")
		return
	}

	// 检查密码
	if share.Password != "" {
		if password == "" {
			response.Success(c, gin.H{
				"requirePassword": true,
			})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(share.Password), []byte(password)); err != nil {
			response.Error(c, 403, "访问密码错误")
			return
		}
	}

	// 增加访问次数
	db.Model(&share).Update("view_count", share.ViewCount+1)

	// 获取文件信息
	var file model.File
	db.First(&file, share.FileID)

	response.Success(c, gin.H{
		"requirePassword": false,
		"fileName":        file.FileName,
		"fileSize":        file.FileSize,
		"fileType":        file.FileType,
	})
}

func downloadShare(c *gin.Context) {
	shareUuid := c.Param("shareUuid")
	password := c.Query("password")

	var share model.Share
	if err := db.Where("share_uuid = ?", shareUuid).First(&share).Error; err != nil {
		response.NotFound(c, "分享不存在")
		return
	}

	// 检查过期
	if share.ExpireTime != nil && share.ExpireTime.Before(time.Now()) {
		response.Error(c, 410, "分享已过期")
		return
	}

	// 检查密码
	if share.Password != "" {
		if password == "" || bcrypt.CompareHashAndPassword([]byte(share.Password), []byte(password)) != nil {
			response.Error(c, 403, "访问密码错误")
			return
		}
	}

	// 获取文件
	var file model.File
	if err := db.First(&file, share.FileID).Error; err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	// 这里简化处理，实际需要从 COS 下载
	c.JSON(http.StatusOK, gin.H{
		"message": "请通过文件服务下载",
		"fileName": file.FileName,
	})
}

func getMyShares(c *gin.Context) {
	userId := getUserId(c)

	var shares []model.Share
	if err := db.Where("created_by = ?", userId).Order("created_at DESC").Find(&shares).Error; err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	records := make([]gin.H, 0, len(shares))
	for _, s := range shares {
		var file model.File
		if err := db.First(&file, s.FileID).Error; err == nil {
			records = append(records, gin.H{
				"shareUuid":  s.ShareUuid,
				"fileName":   file.FileName,
				"viewCount":  s.ViewCount,
				"expireTime": s.ExpireTime,
				"createdAt":  s.CreatedAt,
			})
		}
	}

	response.Success(c, records)
}

func deleteShare(c *gin.Context) {
	userId := getUserId(c)
	shareUuid := c.Param("shareUuid")

	var share model.Share
	if err := db.Where("share_uuid = ? AND created_by = ?", shareUuid, userId).First(&share).Error; err != nil {
		response.NotFound(c, "分享不存在")
		return
	}

	db.Delete(&share)
	response.SuccessWithMessage(c, "删除成功", nil)
}

func calculateExpireTime(expireType string) time.Time {
	now := time.Now()
	switch expireType {
	case "1H":
		return now.Add(1 * time.Hour)
	case "24H":
		return now.Add(24 * time.Hour)
	case "7D":
		return now.Add(7 * 24 * time.Hour)
	default:
		return time.Time{}
	}
}

func generateShortCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
