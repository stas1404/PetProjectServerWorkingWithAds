package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/app"
	"server/internal/ports/httpgin/cookie"
)

func SetUpCheckAuthentification(a app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie(cookie.CookieName)
		fmt.Println(cookieValue, err)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		u, err := a.GetUserByCookie(c, cookieValue)
		fmt.Println(u, err)
		if err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}
		c.Set("UserID", u.ID)
	}
}
