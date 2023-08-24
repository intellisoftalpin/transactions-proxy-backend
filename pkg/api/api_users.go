package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"

	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	"github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/repository"
	utils "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/utils"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
)

type UsersAPI struct {
	usersRepo *repository.UsersRepo
}

func NewUsersAPI(db *sql.DB, config *models.Config, sessions *models.Sessions) *UsersAPI {
	return &UsersAPI{
		usersRepo: repository.NewUsersRepo(db, config, sessions),
	}
}

// LoginUser godoc
// @Summary Login User.
// @Description Login user for 24 hours and create new user in database if needed.
// @Tags users
// @Accept json
// @Produce json
// @Param User body models.CardanoUser true "Struct for requested cardano user with hash of the wallet address"
// @Success 200 {object} models.UserSessionResponse
// @Router /api/v1/user/login [POST]
func (api *UsersAPI) LoginUser(c echo.Context) error {
	// Get data (user_hash string) from POST request
	var cardanoUser models.CardanoUser
	defer c.Request().Body.Close()

	if err := json.NewDecoder(c.Request().Body).Decode(&cardanoUser); err != nil {
		log.Printf("Failed reading the request body: %s", err)
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	if v := validate.Struct(cardanoUser); !v.Validate() {
		return utils.PrepareErrorResponse(c, v.Errors.String(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	// Put user`s wallet address hash from request into logic to create 24 hours session
	userResponse, err := api.usersRepo.LoginUserSession(cardanoUser.UserHash, cardanoUser.UserRuntime)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, userResponse)
}

// UsersSessionsEmptify - function for wrapping internal session emptifying function
func (api *UsersAPI) UsersSessionsEmptify() {
	api.usersRepo.SessionsEmptify()
}
