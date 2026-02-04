package middleware

import (
	"net/http"
	"pi/pkg/utils/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "no tokken",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "incorect format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "incorrect or expired tokken",
			})
			c.Abort()
			return
		}

		c.Set("worker_id", claims.WorkerID)
		c.Set("worker_name", claims.Name)
		c.Set("worker_jobtitle", claims.JobTitle)
		c.Set("department_id", claims.DepartmentID)
		c.Set("department_name", claims.DepartmentName)
		c.Set("access_level", claims.AccessLevel)

		c.Next()
	}
}
