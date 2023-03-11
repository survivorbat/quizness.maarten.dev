package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/routes/inputs"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type QuizHandler struct {
	QuizService services.IQuizService
}

// Get godoc
//
//	@Summary	Fetch your quizzes
//	@Tags		Quiz
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		[]domain.Quiz	"Your quizzes"
//	@Failure	500	{object}	any				"Internal Server Error"
//	@Router		/api/v1/quizzes [get]
//	@Security	JWT
func (g *QuizHandler) Get(c *gin.Context) {
	authID := c.GetString("user")

	quizzes, err := g.QuizService.GetByCreator(uuid.MustParse(authID))
	if err != nil {
		logrus.WithError(err).Error("Failed to get by creator")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// Post godoc
//
//	@Summary	Create a quiz
//	@Tags		Quiz
//	@Accept		json
//	@Produce	json
//	@Param		input	body		inputs.Quiz	true	"Your quiz"
//	@Success	200		{object}	inputs.Quiz	"Your quiz"
//	@Failure	500		{object}	any			"Internal Server Error"
//	@Router		/api/v1/quizzes [post]
//	@Security	JWT
func (g *QuizHandler) Post(c *gin.Context) {
	authID := c.GetString("user")

	var input *inputs.Quiz
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithError(err).Error("Failed to parse input")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	quiz := input.ToDomain()
	quiz.CreatorID = uuid.MustParse(authID)

	logrus.Infof("Creating %#v", quiz)
	if err := g.QuizService.CreateOrUpdate(quiz); err != nil {
		logrus.WithError(err).Error("Failed to create")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// Put godoc
//
//	@Summary	Update or create a quiz
//	@Tags		Quiz
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string		true	"ID of the quiz"
//	@Param		input	body		inputs.Quiz	true	"Your quiz"
//	@Success	200		{object}	inputs.Quiz	"Your quiz"
//	@Failure	403		{object}	any			"You can only update your own quizzes"
//	@Failure	500		{object}	any			"Internal Server Error"
//	@Router		/api/v1/quizzes/{id} [put]
//	@Security	JWT
func (g *QuizHandler) Put(c *gin.Context) {
	authID := c.GetString("user")
	id := c.Param("id")

	quizID, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	quiz, err := g.QuizService.GetByID(quizID)
	if err == nil && quiz.CreatorID.String() != authID {
		logrus.Errorf("Creator is %s not %s", quiz.CreatorID, authID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var input *inputs.Quiz
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithError(err).Error("Failed to parse input")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	update := input.ToDomain()
	update.CreatorID = uuid.MustParse(authID)
	update.ID = quizID

	logrus.Infof("Overwriting %#v", quiz)
	if err := g.QuizService.CreateOrUpdate(update); err != nil {
		logrus.WithError(err).Error("Failed create or update")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// Delete godoc
//
//	@Summary	Delete a quiz
//	@Tags		Quiz
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"ID of the quiz"
//	@Success	204	{object}	any
//	@Failure	403	{object}	any	"You can only delete your own quizzes"
//	@Failure	500	{object}	any	"Internal Server Error"
//	@Router		/api/v1/quizzes/{id} [delete]
//	@Security	JWT
func (g *QuizHandler) Delete(c *gin.Context) {
	authID := c.GetString("user")
	id := c.Param("id")

	quizID, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	quiz, err := g.QuizService.GetByID(quizID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Prevent users from deleting other people's quizzes
	if quiz.CreatorID.String() != authID {
		logrus.Errorf("Creator is %s not %s", quiz.CreatorID, authID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	logrus.Infof("Deleting %#v", quiz)
	if err := g.QuizService.Delete(quiz.ID); err != nil {
		logrus.WithError(err).Error("Failed to delete")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, quiz)
}
