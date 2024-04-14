package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kellswork/wayfarer/internal/utils"
)

func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "error",
				"message": "Access denied. no authorisation provided",
			})
			c.Abort()
			return
		}
		token := strings.Split(authHeader, " ")[1]

		claims, err := utils.VerifyJwtToken(token)

		if err != nil {
			log.Printf("failed to decode token : %v\n", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "error",
				"message": "Access denied. Authentication failed",
			})
			c.Abort()
			return
		}

		user := map[string]interface{}{
			"id":      claims["id"].(string),
			"isAdmin": claims["isAdmin"].(bool),
			"exp":     claims["exp"].(float64),
		}
		c.Set("user", user)
		c.Next()
	}

}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		value, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "error",
				"message": "Access denied. No authorisation provided",
			})
			c.Abort()
			return

		}

		user, ok := value.(map[string]interface{})

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "error",
				"message": "An error occured please refresh or try again",
			})
			c.Abort()
			return
		}

		if !user["isAdmin"].(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "error",
				"message": "Access denied. You are not authorised to access this resource",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
