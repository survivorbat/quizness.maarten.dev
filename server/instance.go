package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewServer(connectionString string, jwtSecret string, oAuthID string, oAuthSecret string, authRedirectUrl string) (*Server, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Server{
		database:  db,
		jwtSecret: jwtSecret,
		oAuthConfig: &oauth2.Config{
			ClientID:     oAuthID,
			ClientSecret: oAuthSecret,
			RedirectURL:  authRedirectUrl,
			Scopes:       []string{"openid"},
			Endpoint:     google.Endpoint,
		},
	}, nil
}

type Server struct {
	database  *gorm.DB
	jwtSecret string

	// Configs
	oAuthConfig *oauth2.Config

	// Handlers
	authHandler       *routes.TokenHandler
	jwtHandler        *routes.JwtHandler
	getCreatorHandler *routes.CreatorHandler
}

func (s *Server) Configure(router *gin.Engine) error {
	if err := s.database.AutoMigrate(&domain.Game{}, &domain.Quiz{}, &domain.Creator{}, &domain.Player{}); err != nil {
		logrus.WithError(err).Error("Failed to migrate")
		return err
	}

	s.configureServices()
	s.configureRoutes(router)
	return nil
}

func (s *Server) configureServices() {
	creatorService := &services.CreatorService{Database: s.database}
	jwtService := &services.JwtService{SecretKey: s.jwtSecret, Issuer: "QQ"}
	s.authHandler = &routes.TokenHandler{CreatorService: creatorService, JwtService: jwtService, AuthConfig: s.oAuthConfig}
	s.jwtHandler = &routes.JwtHandler{JwtService: jwtService}
}

func (s *Server) configureRoutes(router *gin.Engine) {
	router.POST("/api/v1/tokens", s.authHandler.CreateToken)

	// Guarded routes with JWT
	apiRoutes := router.Group("/api")
	apiRoutes.Use(s.jwtHandler.JwtGuard())
	apiRoutes.GET("/creators/:id", s.getCreatorHandler.GetWithID)
}
