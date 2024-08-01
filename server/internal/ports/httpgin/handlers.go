package httpgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/stas1404/validator"
	"net/http"
	"server/internal/ads"
	"server/internal/app"
	errors2 "server/internal/errors"
	"server/internal/errors/cookie_errors"
	"server/internal/ports"
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
			GetStatusAndAbort(err, c)
			return
		}
		ad, err := a.GetAd(c, id)
		if err != nil {
			GetStatusAndAbort(err, c)
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
			GetStatusAndAbort(err, c)
			return
		}
		fmt.Println(res)
		ads := a.GetAdsCorresponding(c, res)
		c.JSON(http.StatusOK, ads)
	}
}

func SetUpCreateUser(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var u ports.ResponseUser
		err := c.BindJSON(&u)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		if err = u.Check(); err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		us, err := a.CreateUser(c, u.Nickname, u.Email, u.Password)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		c.JSON(http.StatusCreated, us)
	}
}

func SetUpAuthorization(a app.App, cookies cookie.CookieRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var reqUser ports.ResponseUser
		err := c.BindJSON(&reqUser)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		if err = reqUser.Check(); err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		exUser, err := a.GetUserByID(c, reqUser.ID)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		fmt.Println(exUser, reqUser)
		if !user.AreSame(exUser, user.User(reqUser)) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		co := http.Cookie{
			Name:     cookie.CookieName,
			Value:    cookie.GenerateCookieValue(exUser),
			Expires:  time.Now().Add(time.Hour),
			Path:     "/",
			HttpOnly: true,
		}
		cookies.AddCookie(co, exUser.ID)
		http.SetCookie(c.Writer, &co)
	}
}

func SetUpCreateAd(app app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		var ad ads.Ad
		err := c.BindJSON(&ad)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		ad, err = app.CreateAd(c, ad.Title, ad.Text, c.GetInt64("UserID"))
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		c.JSON(http.StatusCreated, ad)
	}
}

func SetUpModifyAd(a app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		ad, err := a.GetAd(c, id)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		if ad.AuthorID != c.GetInt64("UserID") {
			GetStatusAndAbort(errors2.NewErrWrongUserID(c.GetInt64("UserID")), c)
			return
		}
		err = c.BindJSON(&ad)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		ad, err = a.UpdateAd(c, ad.ID, ad.AuthorID, ad.Title, ad.Text)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		c.JSON(http.StatusOK, ad)
	}
}

func SetUpPublishAd(a app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		UserID := c.GetInt64("UserID")
		AdID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			GetStatusAndAbort(err, c)
		}
		ad, err := a.ChangeAdStatus(c, AdID, UserID, true)
		if err != nil {
			GetStatusAndAbort(err, c)
		}
		c.JSON(http.StatusOK, ad)
	}
}

func SetUpUnPublishAd(a app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		UserID := c.GetInt64("UserID")
		AdID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			GetStatusAndAbort(err, c)
		}
		ad, err := a.ChangeAdStatus(c, AdID, UserID, false)
		if err != nil {
			GetStatusAndAbort(err, c)
		}
		c.JSON(http.StatusOK, ad)
	}
}

func SetUpEditUser(a app.App) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.GetInt64("UserID")
		us, err := a.GetUserByID(c, id)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		err = c.BindJSON(&us)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		us.ID = id
		u := ports.ResponseUser(us)
		if err = u.Check(); err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		us, err = a.UpdateUser(c, id, u.Nickname, u.Email, u.Password)
		if err != nil {
			GetStatusAndAbort(err, c)
		}
		c.JSON(http.StatusOK, us)
	}
}

func GetStatusAndAbort(err error, c *gin.Context) {
	code := http.StatusNotFound
	switch err.(type) {
	case errors2.ErrBadUser:
		code = http.StatusBadRequest
	case errors2.ErrWrongUserID:
		code = http.StatusBadRequest
	case errors2.PermissionDenied:
		code = http.StatusBadRequest
	case errors2.ErrUnexistingUser:
		code = http.StatusNotFound
	case errors2.ErrUnexistingAd:
		code = http.StatusNotFound
	case cookie_errors.ErrUnexistingCookie:
		code = http.StatusUnauthorized
	case cookie_errors.ExpiredCookie:
		code = http.StatusUnauthorized
	case validator.ValidationErrors:
		code = http.StatusBadRequest
	}
	c.Status(code)
	c.Writer.Write([]byte("<p>" + err.Error() + "</p>"))
	c.AbortWithError(code, err)
}
