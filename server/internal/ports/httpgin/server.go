package httpgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/app"
	"server/internal/ports/httpgin/middleware"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	Authentification := middleware.SetUpCheckAuthentification(a)
	s := Server{port: port, app: gin.New()}
	s.app.Use(gin.Recovery())
	s.app.GET("/ads/:id", SetUpGetAdByID(a))
	s.app.GET("/ads/", SetUpGetAdCorresponding(a))
	s.app.POST("/users/", SetUpCreateUser(a))
	s.app.POST("/ads/", Authentification, SetUpCreateAd(a))
	s.app.POST("/users/authorization/", SetUpAuthorization(a))
	s.app.PUT("/ads/:id", Authentification, SetUpModifyAd(a))
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
