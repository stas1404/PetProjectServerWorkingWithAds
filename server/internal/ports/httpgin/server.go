package httpgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/app"
	"server/internal/ports/httpgin/cookie"
	"server/internal/ports/httpgin/middleware"
)

type Server struct {
	port    string
	app     *gin.Engine
	cookies cookie.CookieRepository
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}
	repo := cookie.NewRepository()
	Authentification := middleware.SetUpCheckAuthentification(repo)
	s.app.Use(gin.Logger())
	s.app.Use(gin.Recovery())
	s.app.GET("/ads/:id", SetUpGetAdByID(a))
	s.app.GET("/ads/", SetUpGetAdCorresponding(a))
	s.app.POST("/users/", SetUpCreateUser(a))
	s.app.POST("/ads/", Authentification, SetUpCreateAd(a))
	s.app.POST("/users/authorization/", SetUpAuthorization(a, repo))
	s.app.PUT("/ads/:id/edit", Authentification, SetUpModifyAd(a))
	s.app.PUT("/ads/:id/edit/publish", Authentification, SetUpPublishAd(a))
	s.app.PUT("/ads/:id/edit/unpublish", Authentification, SetUpUnPublishAd(a))
	s.app.PUT("/users/profile/edit", Authentification, SetUpEditUser(a))
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
