package server

import (
	"github.com/gin-gonic/gin"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"github.com/zalando/gin-oauth2/google"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewServer(connectionString string, jwtSecret string, authCredentialsPath string, authSecret []byte, authRedirectUrl string) (*Server, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Server{database: db, jwtSecret: jwtSecret, credentialsPath: authCredentialsPath, oAuthSecret: authSecret, redirectUrl: authRedirectUrl}, nil
}

type Server struct {
	database        *gorm.DB
	credentialsPath string
	redirectUrl     string
	oAuthSecret     []byte
	jwtSecret       string

	// Handlers
	authHandler       *routes.AuthHandler
	jwtHandler        *routes.JwtHandler
	getCreatorHandler *routes.GetCreatorHandler
}

func (s *Server) Configure(router *gin.Engine) error {
	if err := s.database.AutoMigrate(&domain.Game{}, &domain.Quiz{}, &domain.Creator{}, &domain.Player{}); err != nil {
		return err
	}

	// Authentication setup
	google.Setup(s.redirectUrl, s.credentialsPath, []string{"openid"}, s.oAuthSecret)
	router.Use(google.Session("qq"))

	s.configureServices()
	s.configureRoutes(router)

	return nil
}

func (s *Server) configureServices() {
	creatorService := &services.CreatorService{Database: s.database}
	jwtService := &services.JwtService{SecretKey: s.jwtSecret, Issuer: "QQ"}
	s.authHandler = &routes.AuthHandler{CreatorService: creatorService, JwtService: jwtService}
	s.jwtHandler = &routes.JwtHandler{JwtService: jwtService}
}

func (s *Server) configureRoutes(router *gin.Engine) {
	router.GET("/login", google.LoginHandler)

	// Specific route for handling the login
	googleRoute := router.Group("/auth")
	googleRoute.Use(google.Auth())
	googleRoute.GET("/", s.authHandler.Handle)

	// Guarded routes with JWT
	apiRoutes := router.Group("/api")
	apiRoutes.Use(s.jwtHandler.AuthorizeJWT())
	apiRoutes.GET("/creators/:id", s.getCreatorHandler.GetWithID)
}
