package main

import (
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var identityKey = "id"

type UserLogin struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserDb struct {
	gorm.Model
	Name     string `gorm:"primaryKey"`
	Password string
}
type UserClaims struct {
	Name string
}

func makeAuthMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "clash zone",
		Key:         []byte(os.Getenv("JWT_KEY")),
		Timeout:     time.Hour * 24 * 180,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*UserClaims); ok {
				return jwt.MapClaims{
					identityKey: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &UserClaims{
				Name: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals UserLogin
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userName := loginVals.Name
			password := loginVals.Password

			var dbVals UserDb
			if db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{}); err == nil {
				db.AutoMigrate(&UserDb{})

				var userCount int64
				db.Model(&UserDb{}).Count(&userCount)

				if userCount == 0 && userName == "admin" {
					db.Create(&UserDb{Name: userName, Password: password})
				}

				db.First(&dbVals, "name = ?", userName)
			} else {
				log.Panic(err)
			}

			if password == dbVals.Password {
				return &UserClaims{
					Name: userName,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*UserClaims); ok && v.Name == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
