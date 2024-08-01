package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/ports/httpgin/cookie"
)

func SetUpCheckAuthentification(repo cookie.CookieRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie(cookie.CookieName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		u, err := repo.GetUserIDByCookieValue(cookieValue)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.Set("UserID", u)
	}
}
