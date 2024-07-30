package httpgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/ads"
	"server/internal/app"
	errors2 "server/internal/errors"
	"server/internal/ports/httpgin/cookie"
	"server/internal/restriction"
	"server/internal/user"
	"strconv"
	"time"
)

func SetUpGetAdByID(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.Status(http.StatusMisdirectedRequest)
			c.Writer.Write([]byte("<p>" + err.Error() + "</p>"))
			return
		}
		ad, err := a.GetAd(c, id)
		if err != nil {
			defer c.Writer.Write([]byte("<p>" + err.Error() + "</p>"))
			if err.Error() == errors2.NewErrUnexistingAd(id).Error() {
				c.Status(http.StatusNotFound)
				return
			}
			c.Status(http.StatusInsufficientStorage)
			return
		}
		c.JSON(http.StatusOK, ad)
	}
}

func SetUpGetAdCorresponding(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var res restriction.Restriction
		err := c.ShouldBind(&res)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.Writer.Write([]byte("<p>" + err.Error() + "</p>"))
			return
		}
		ads := a.GetAdsCorresponding(c, res)
		c.JSON(http.StatusOK, ads)
	}
}

func SetUpCreateUser(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var u user.User
		err := c.BindJSON(&u)
		if err != nil {
			WriteError(http.StatusMethodNotAllowed, err, c)
			return
		}
		if err = u.Check(); err != nil {
			WriteError(http.StatusMethodNotAllowed, err, c)
			return
		}
		cookie := http.Cookie{
			Name:    cookie.CookieName,
			Value:   cookie.GenerateCookieValue(u),
			Expires: time.Now().Add(time.Hour),
			Path:    "/",
		}
		u, err = a.CreateUser(c, u.Nickname, u.Email, u.Password, cookie)
		if err != nil {
			WriteError(http.StatusInsufficientStorage, err, c)
			return
		}
		http.SetCookie(c.Writer, &cookie)
		c.JSON(http.StatusCreated, u)
	}
}

func SetUpAuthorization(a app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		var reqUser user.User
		err := c.BindJSON(&reqUser)
		if err != nil {
			WriteError(http.StatusMethodNotAllowed, err, c)
			return
		}
		if err = reqUser.Check(); err != nil {
			WriteError(http.StatusMethodNotAllowed, err, c)
			return
		}
		exUser, err := a.GetUserByID(c, reqUser.ID)
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
		}
		if !user.AreSame(exUser, reqUser) {
			c.AbortWithStatus(http.StatusForbidden)
		}
		cookie := http.Cookie{
			Name:     cookie.CookieName,
			Value:    cookie.GenerateCookieValue(exUser),
			Expires:  time.Now().Add(time.Hour),
			Path:     "/",
			HttpOnly: true,
		}
		a.CreateCookie(c, cookie, exUser.ID)
		http.SetCookie(c.Writer, &cookie)
	}
}

func SetUpCreateAd(app app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		var ad ads.Ad
		err := c.BindJSON(&ad)
		if err != nil {
			WriteError(http.StatusBadRequest, err, c)
			return
		}
		fmt.Println(ad, err)
		ad, err = app.CreateAd(c, ad.Title, ad.Text, c.GetInt64("UserID")) //c.Value("UserID").(int64))
		if err != nil {
			WriteError(http.StatusBadRequest, err, c)
			return
		}
		c.JSON(http.StatusCreated, ad)
	}
}

func SetUpModifyAd(a app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WriteError(http.StatusMisdirectedRequest, err, c)
			return
		}
		ad, err := a.GetAd(c, id)
		if err != nil {
			if err.Error() == errors2.NewErrUnexistingAd(id).Error() {
				WriteError(http.StatusNotFound, err, c)
				return
			}
			WriteError(http.StatusInsufficientStorage, err, c)
			return
		}
		if ad.AuthorID != c.GetInt64("UserID") {
			WriteError(http.StatusForbidden, errors2.NewErrWrongUserID(c.GetInt64("UserID")), c)
			return
		}
		err = c.BindJSON(&ad)
		if err != nil {
			WriteError(http.StatusBadRequest, err, c)
			return
		}
		fmt.Println("New add: ", ad)
		ad, err = a.ChangeAdStatus(c, ad.ID, ad.AuthorID, ad.Published)
		if err != nil {
			WriteError(http.StatusBadRequest, err, c)
			return
		}
		ad, err = a.UpdateAd(c, ad.ID, ad.AuthorID, ad.Title, ad.Text)
		if err != nil {
			WriteError(http.StatusBadRequest, err, c)
			return
		}
		c.JSON(http.StatusOK, ad)
	}
}

func WriteError(code int, err error, c *gin.Context) {
	c.Status(code)
	c.Writer.Write([]byte("<p>" + err.Error() + "</p>"))
}
