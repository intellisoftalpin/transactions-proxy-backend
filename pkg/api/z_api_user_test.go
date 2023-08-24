package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bykovme/goconfig"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	tlpsdb "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/db"
	"github.com/intellisoftalpin/transactions-proxy-backend/models"
)

const confPath = "etc/transactions-proxy-backend-server/test_config.json"

var (
	loadedConfig     *models.Config
	e                *echo.Echo
	testHandler      *UsersAPI
	testUserJSON_1_1 string = `{"user_hash":"user1","user_runtime":0}`
	testUserJSON_1_2 string = `{"user_hash":"user1","user_runtime":1}`
	testUserJSON_2   string = `{"user_hash":"user2","user_runtime":1}`
	testUserJSON_3   string = `{"user_hash":"user3","user_runtime":2}`
)

func setUpNewContextForPostUserLoginRequest(e *echo.Echo, testUser string) (context echo.Context, response *httptest.ResponseRecorder) {

	// Setup request and reply to create new context
	request := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(testUser))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	context = e.NewContext(request, response)

	return context, response
}

func setUpResponse(status string, authorizationKey string) (response models.UserSessionResponse) {
	response.Status = status
	response.AuthorizationKey = authorizationKey

	return response
}

func TestParallelLoginUsers(t *testing.T) {

	// Setup new echo
	e = echo.New()

	// Get user homepath to load config
	usrHomePath, err := goconfig.GetUserHomePath()
	if err != nil {
		t.Errorf("Unexpected error while getting user home path: %s", err)
	}

	if err = goconfig.LoadConfig(usrHomePath+confPath, &loadedConfig); err != nil {
		t.Errorf("Unexpected error while loading config: %s", err)
	}

	// Setup db connection
	db, err := tlpsdb.SetupDB(loadedConfig.DB)
	if err != nil {
		t.Errorf("Unexpected error while database setup: %s", err)
	}

	// Close db connection on finish
	defer db.Close()

	var sessions *models.Sessions = &models.Sessions{SessionsMap: make(map[string]*models.Session)}

	testHandler = NewUsersAPI(db, loadedConfig, sessions)

	t.Run("UserLogin1", func(t *testing.T) {
		context, response := setUpNewContextForPostUserLoginRequest(e, testUserJSON_1_1)

		testRes := setUpResponse(consts.CSuccessStatus, "0")

		if assert.NoError(t, testHandler.LoginUser(context)) {
			var tempResponse models.UserSessionResponse

			err := json.Unmarshal(response.Body.Bytes(), &tempResponse)

			assert.Equal(t, err, nil)
			assert.Equal(t, testRes.Status, tempResponse.Status)
			assert.Equal(t, testRes.AuthorizationKey, tempResponse.AuthorizationKey)
		}

		t.Parallel()
	})

	t.Run("UserLogin1", func(t *testing.T) {
		context, response := setUpNewContextForPostUserLoginRequest(e, testUserJSON_1_2)

		testRes := setUpResponse(consts.CSuccessStatus, "1")

		if assert.NoError(t, testHandler.LoginUser(context)) {
			var tempResponse models.UserSessionResponse

			err := json.Unmarshal(response.Body.Bytes(), &tempResponse)

			assert.Equal(t, err, nil)
			assert.Equal(t, testRes.Status, tempResponse.Status)
			assert.Equal(t, testRes.AuthorizationKey, tempResponse.AuthorizationKey)
		}

		t.Parallel()
	})

	t.Run("UserLogin2", func(t *testing.T) {
		context, response := setUpNewContextForPostUserLoginRequest(e, testUserJSON_2)

		testRes := setUpResponse(consts.CSuccessStatus, "2")

		if assert.NoError(t, testHandler.LoginUser(context)) {
			var tempResponse models.UserSessionResponse

			err := json.Unmarshal(response.Body.Bytes(), &tempResponse)

			assert.Equal(t, err, nil)
			assert.Equal(t, testRes.Status, tempResponse.Status)
			assert.Equal(t, testRes.AuthorizationKey, tempResponse.AuthorizationKey)
		}

		t.Parallel()
	})

	t.Run("UserLogin3", func(t *testing.T) {
		context, response := setUpNewContextForPostUserLoginRequest(e, testUserJSON_3)

		testRes := setUpResponse(consts.CSuccessStatus, "3")

		if assert.NoError(t, testHandler.LoginUser(context)) {
			var tempResponse models.UserSessionResponse

			err := json.Unmarshal(response.Body.Bytes(), &tempResponse)

			assert.Equal(t, err, nil)
			assert.Equal(t, testRes.Status, tempResponse.Status)
			assert.Equal(t, testRes.AuthorizationKey, tempResponse.AuthorizationKey)
		}

		t.Parallel()
	})
}
