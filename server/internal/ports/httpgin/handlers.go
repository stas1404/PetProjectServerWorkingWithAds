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

// ShowAd godoc
// @Summary      Show an ad
// @Description  get ad by ID
// @Tags         ads
// @Produce      json
// @Param        id   path      int  true  "Ad ID"
// @Success      200  {object}  ports.ResponseAd
// @Failure      400  {object} string
// @Failure      404  {object} string
// @Router       /ads/{id} [get]
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

// ShowAds godoc
// @Summary      Show ads
// @Description  get ads corresponding some restrictions
// @Tags         ads
// @Produce      json
// @Param        published   query     []bool  true  "Ad Published"
// @Param        title   query     []string  true  "Ad title"
// @Param        author   query     []int64  true  "Ad author ID"
// @Param        created   query     time.Time  true  "Ad creation time" time_format:"2006-01-02" time_utc:"1"
// @Success      200  {object}  ports.ResponseAd
// @Failure      400  {object} string
// @Router       /ads [get]
func SetUpGetAdCorresponding(a app.App) func(c *gin.Context) {
	return func(c *gin.Context) {
		var res restriction.Restriction
		err := c.ShouldBind(&res)
		if err != nil {
			GetStatusAndAbort(err, c)
			return
		}
		ads := a.GetAdsCorresponding(c, res)
		c.JSON(http.StatusOK, ads)
	}
}

// CreateUser godoc
// @Summary      Create User
// @Description  Create User
// @Tags         users
// @Accept json
// @Produce      json
// @Param nickname body string true "User nickname"
// @Param email body string true "User email"
// @Param password body string true "User password"
// @Success      200  {object}  ports.ResponseUser
// @Failure      400  {object} string
// @Router       /users [post]
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

// AuthorizeUser godoc
// @Summary      Authorize User
// @Description  Authorize User
// @Tags         users
// @Accept json
// @Produce      json
// @Param id body int64 true "User ID"
// @Param nickname body string true "User nickname"
// @Param email body string true "User email"
// @Param password body string true "User password"
// @Success      200
// @Failure      400  {object} string
// @Router       /users/authorization [post]
func SetUpAuthorization(a app.App, cookies cookie.CookieRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var reqUser ports.ResponseUser
		err := c.BindJSON(&reqUser)
		if err != nil {
			fmt.Println("Marshalling")
			GetStatusAndAbort(err, c)
			return
		}
		if err = reqUser.Check(); err != nil {
			fmt.Println("Check")
			GetStatusAndAbort(err, c)
			return
		}
		exUser, err := a.GetUserByID(c, reqUser.ID)
		if err != nil {
			fmt.Println("Get User By id")
			GetStatusAndAbort(err, c)
			return
		}
		if !user.AreSame(exUser, user.User(reqUser)) {
			fmt.Println("ne saMe")
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

// CreateAD godoc
// @Summary      Create ad
// @Description  Create ad
// @Tags         ads
// @Accept json
// @Produce      json
// @Param text body string true "text of ad"
// @Param title body string true "title of ad"
// @Security ApiKeyAuth
// @Success      201  {object}  ports.ResponseAd
// @Failure      400  {object} string
// @Failure      401  {object} string
// @Router       /ads [post]
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

// ModifyAD godoc
// @Summary      Modify ad
// @Description  Modify existing ad by passing new title and text
// @Tags         ads
// @Accept json
// @Produce      json
// @Param text body string true "text of ad"
// @Param title body string true "title of ad"
// @Security ApiKeyAuth
// @Success      200  {object}  ports.ResponseAd
// @Failure      400  {object} string
// @Failure      401  {object} string
// @Router       /ads/{id}/edit [put]
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
			GetStatusAndAbort(errors2.NewPermissionDenied(c.GetInt64("UserID")), c)
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

// PublishAD godoc
// @Summary      Publish ad
// @Description  Change ad status to Published
// @Tags         ads
// @Produce      json
// @Security ApiKeyAuth
// @Success      200  {object}  ports.ResponseAd
// @Failure      400  {object} string
// @Failure      401  {object} string
// @Failure      404  {object} string
// @Router       /ads/{id}/publish [put]
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

// UnPublishAD godoc
// @Summary      UnPublish ad
// @Description  Change ad status to UnPublished
// @Tags         ads
// @Produce      json
// @Security ApiKeyAuth
// @Success      200  {object}  ports.ResponseAd
// @Failure      400  {object} string
// @Failure      401  {object} string
// @Failure      404  {object} string
// @Router       /ads/{id}/unpublish [put]
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

// EditUser godoc
// @Summary      Edit User
// @Description  Edit User's password, nickname and email
// @Tags         users
// @Accept json
// @Produce      json
// @Param nickname body string true "User nickname"
// @Param email body string true "User email"
// @Param password body string true "User password"
// @Success      200  {object}  ports.ResponseUser
// @Failure      400  {object} string
// @Failure      401  {object} string
// @Failure      404  {object} string
// @Router       /users/profile/edit [post]
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

func GetStatusAndAbort(err any, c *gin.Context) {
	code := http.StatusNotFound
	switch err.(type) {
	case errors2.ErrBadUser:
		code = http.StatusBadRequest
	case errors2.PermissionDenied:
		code = http.StatusForbidden
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
	c.Writer.Write([]byte("<p>" + err.(error).Error() + "</p>"))
	c.Abort()
	//c.AbortWithError(code, err.(error).Error())
}
