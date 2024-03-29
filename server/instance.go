package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/coordinator"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
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

var databaseOpen = postgres.Open

func NewServer(connectionString string, jwtSecret string, oAuthID string, oAuthSecret string, authRedirectUrl string) (*Server, error) {
	db, err := gorm.Open(databaseOpen(connectionString))
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
	jwtService  services.JwtService

	// Handlers
	tokenHandler          *routes.TokenHandler
	gameControlHandler    *routes.GameControlHandler
	creatorHandler        *routes.CreatorHandler
	quizHandler           *routes.QuizHandler
	playerHandler         *routes.PlayerHandler
	publicGameHandler     *routes.PublicGameHandler
	gameConnectionHandler *routes.GameConnectionHandler
}

func (s *Server) Configure(router *gin.Engine) error {
	if err := s.database.AutoMigrate(
		&domain.Game{},
		&domain.Quiz{},
		&domain.Creator{},
		&domain.Player{},
		&domain.QuestionOption{},
		&domain.GameAnswer{},
	); err != nil {
		logrus.WithError(err).Error("Failed to migrate")
		return err
	}

	s.configureServices()
	s.configureRoutes(router)
	s.configureValidator()
	return nil
}

func (s *Server) configureServices() {
	quizService := &services.DBQuizService{Database: s.database}
	creatorService := &services.DBCreatorService{Database: s.database}
	gameService := &services.DBGameService{Database: s.database}
	playerService := &services.DBPlayerService{Database: s.database}

	gameCoordinator := &coordinator.LocalGameCoordinator{GameService: gameService}

	s.jwtService = &services.HMacJwtService{SecretKey: s.jwtSecret, Issuer: "QQ"}

	s.tokenHandler = &routes.TokenHandler{CreatorService: creatorService, JwtService: s.jwtService, AuthConfig: s.oAuthConfig}
	s.quizHandler = &routes.QuizHandler{QuizService: quizService}
	s.creatorHandler = &routes.CreatorHandler{CreatorService: creatorService}
	s.gameControlHandler = &routes.GameControlHandler{GameService: gameService, QuizService: quizService}
	s.playerHandler = &routes.PlayerHandler{PlayerService: playerService, GameService: gameService}
	s.publicGameHandler = &routes.PublicGameHandler{GameService: gameService}
	s.gameConnectionHandler = &routes.GameConnectionHandler{
		GameService:    gameService,
		PlayerService:  playerService,
		CreatorService: creatorService,
		Coordinator:    gameCoordinator,
	}
}

func (s *Server) configureRoutes(router *gin.Engine) {
	router.POST("/api/v1/tokens", s.tokenHandler.CreateToken)

	// Guarded routes with JWT
	apiRoutes := router.Group("/api/v1")
	apiRoutes.Use(s.tokenHandler.JwtGuard())

	apiRoutes.GET("/creators/self", s.creatorHandler.GetWithID)
	apiRoutes.GET("/quizzes", s.quizHandler.Get)
	apiRoutes.GET("/games/:id/players", s.playerHandler.Get)
	apiRoutes.GET("/games/:id/connection", s.gameConnectionHandler.GetCreator)
	apiRoutes.GET("/games/:id", s.gameControlHandler.GetByID)

	apiRoutes.POST("/quizzes", s.quizHandler.Post)
	apiRoutes.POST("/quizzes/:id/games", s.gameControlHandler.Post)

	apiRoutes.PUT("/tokens", s.tokenHandler.Refresh)
	apiRoutes.PUT("/quizzes/:id", s.quizHandler.Put)

	apiRoutes.PATCH("/games/:id", s.gameControlHandler.Patch)

	apiRoutes.DELETE("/quizzes/:id", s.quizHandler.Delete)
	apiRoutes.DELETE("/games/:id", s.gameControlHandler.Delete)

	// Anonymous routes
	publicRoutes := router.Group("/api/v1")
	publicRoutes.GET("/games", s.publicGameHandler.GetByCode)
	publicRoutes.GET("/games/:id/quiz", s.publicGameHandler.GetQuiz)
	publicRoutes.GET("/games/:id/players/:player/connection", s.gameConnectionHandler.Get)
	publicRoutes.POST("/games/:id/players", s.playerHandler.Post)
	publicRoutes.DELETE("/players/:id", s.playerHandler.Delete)

	// Swagger
	url := ginSwagger.URL("/api/swagger/doc.json")
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func (s *Server) configureValidator() {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		val.RegisterStructValidation(inputs.IsValidator, new(inputs.Quiz))
		val.RegisterStructValidation(inputs.IsValidator, new(inputs.MultipleChoiceQuestion))
		return
	}

	logrus.Fatal("Failed to register validation")
}
