package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	_ "github.com/survivorbat/qq.maarten.dev/server/swagger" // docs is generated by Swag CLI, you have to import it.
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	swaggerEndpoint = "/swagger/*any"
	swaggerJsonPath = "/swagger/doc.json"
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
	tokenHandler   *routes.TokenHandler
	jwtHandler     *routes.JwtHandler
	creatorHandler *routes.CreatorHandler
	quizHandler    *routes.QuizHandler
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
	quizService := &services.QuizService{Database: s.database}
	creatorService := &services.CreatorService{Database: s.database}
	jwtService := &services.JwtService{SecretKey: s.jwtSecret, Issuer: "QQ"}

	s.tokenHandler = &routes.TokenHandler{CreatorService: creatorService, JwtService: jwtService, AuthConfig: s.oAuthConfig}
	s.jwtHandler = &routes.JwtHandler{JwtService: jwtService}
	s.quizHandler = &routes.QuizHandler{QuizService: quizService}
	s.creatorHandler = &routes.CreatorHandler{CreatorService: creatorService}
}

func (s *Server) configureRoutes(router *gin.Engine) {
	router.POST("/api/v1/tokens", s.tokenHandler.CreateToken)

	// Guarded routes with JWT
	apiRoutes := router.Group("/api/v1")
	apiRoutes.Use(s.jwtHandler.JwtGuard())
	apiRoutes.GET("/creators/self", s.creatorHandler.GetWithID)
	apiRoutes.GET("/quizzes", s.quizHandler.Get)
	apiRoutes.PUT("/tokens", s.jwtHandler.Refresh)

	// Swagger
	url := ginSwagger.URL(swaggerJsonPath)
	router.GET(swaggerEndpoint, ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
