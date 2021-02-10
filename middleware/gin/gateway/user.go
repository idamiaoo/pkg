package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pescaria/pkg/util"
)

type user struct {
	Id       interface{} `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	SchoolId interface{} `json:"school_id"`
	Type     string      `json:"type"`
}

// 不验证请求来源
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		profile := c.GetHeader("Ms-Profile")
		if len(profile) == 0 {
			c.Next()
			return
		}

		var userInfo user
		if err := json.Unmarshal(util.StringToBytes(profile), &userInfo); err != nil {
			c.Next()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "Authorization", c.GetHeader("Authorization")))

		c.Set("school_id", userInfo.SchoolId)
		c.Set("id", userInfo.Id)
		c.Set("name", userInfo.Name)
		c.Set("email", userInfo.Email)
		c.Set("user_type", userInfo.Type)
		c.Next()
	}
}

// 学生端
func StudentMiddleware() gin.HandlerFunc {
	return checkPlatform("student", "student_all_in_one")
}

// 辅导端
func TeacherMiddleware() gin.HandlerFunc {
	return checkPlatform("live_teacher")
}

// 学生端 辅导端 PC
func checkPlatform(platform ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		profile := c.GetHeader("Ms-Profile")
		fail := map[string]string{
			"message": "认证失败, 请重新登录",
		}
		if len(profile) == 0 {
			c.JSON(http.StatusUnauthorized, fail)
			c.Abort()
			return
		}

		var userInfo user
		if err := json.Unmarshal(util.StringToBytes(profile), &userInfo); err != nil {
			c.JSON(http.StatusUnauthorized, fail)
			c.Abort()
			return
		}

		if !util.InArray(userInfo.Type, platform) {
			c.JSON(http.StatusUnauthorized, fail)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "Authorization", c.GetHeader("Authorization")))

		c.Set("school_id", userInfo.SchoolId)
		c.Set("id", userInfo.Id)
		c.Set("name", userInfo.Name)
		c.Set("email", userInfo.Email)
		c.Set("user_type", userInfo.Type)
		return
	}
}
