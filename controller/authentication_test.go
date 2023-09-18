package controller

import (
	"encoding/json"
	"github/be/common"
	"github/be/database"
	"github/be/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthenticationTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

func (suite *AuthenticationTestSuite) SetupTest() {
	common.Env("../")

	database.Connect()
	if err := database.Database.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}

	if err := database.Database.AutoMigrate(&model.Entry{}); err != nil {
		panic(err)
	}

	user := model.User{
		Username: "testuser0",
		Password: "testpass0",
	}

	if _, err := user.Save(); err != nil {
		panic(err)
	}

	suite.db = database.Database

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	suite.router.POST("/auth/register", Register)
	suite.router.POST("/auth/login", Login)
}

func (suite *AuthenticationTestSuite) TearDownTest() {
	db, _ := suite.db.DB()

	if err := suite.db.Migrator().DropTable(&model.User{}); err != nil {
		panic(err)
	}

	if err := suite.db.Migrator().DropTable(&model.Entry{}); err != nil {
		panic(err)
	}

	db.Close()
}

func (suite *AuthenticationTestSuite) testRegisterWithInput(body string, expectedStatus int, expectedResponse string) {
	req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), expectedStatus, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		panic(err)
	}

	assert.Contains(suite.T(), rr.Body.String(), expectedResponse)
}

func (suite *AuthenticationTestSuite) TestRegister_Success() {
	suite.testRegisterWithInput(`{"username": "testuser1", "password": "testpassword"}`, http.StatusCreated, "user")
}

func (suite *AuthenticationTestSuite) TestRegister_Invalid_Input() {
	suite.testRegisterWithInput(`{"username": "testuser2", "password": ""}`, http.StatusBadRequest, "invalid input")
	suite.testRegisterWithInput(`{"username": "", "password": "testpassword"}`, http.StatusBadRequest, "invalid input")
}

func (suite *AuthenticationTestSuite) TestRegister_Username_Exists() {
	suite.testRegisterWithInput(`{"username": "testuser0", "password": "testpass0"}`, http.StatusBadRequest, "invalid input")
}

func (suite *AuthenticationTestSuite) testLoginWithInput(body string, expectedStatus int, expectedResponse string) {
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), expectedStatus, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		panic(err)
	}

	assert.Contains(suite.T(), rr.Body.String(), expectedResponse)
}

func (suite *AuthenticationTestSuite) TestLogin_Success() {
	suite.testLoginWithInput(`{"username": "testuser0", "password": "testpass0"}`, http.StatusOK, "jwt")
}

func (suite *AuthenticationTestSuite) TestLogin_Invalid_Username() {
	suite.testLoginWithInput(`{"username": "testuser99", "password": "testpass0"}`, http.StatusBadRequest, "invalid username or password")
}

func (suite *AuthenticationTestSuite) TestLogin_Invalid_Input() {
	suite.testLoginWithInput(`{"username": "testuser0", "password": ""}`, http.StatusBadRequest, "username and password are required")
	suite.testLoginWithInput(`{"username": "", "password": "testpass0"}`, http.StatusBadRequest, "username and password are required")
}

func TestAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}
