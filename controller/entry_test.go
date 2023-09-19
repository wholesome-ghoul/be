package controller

import (
	"database/sql"
	"encoding/json"
	"github/be/common"
	"github/be/database"
	"github/be/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type EntryTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *sql.DB
}

var user = model.User{
	Username: "testuser0",
	Password: "testpass0",
}

func (suite *EntryTestSuite) SetupTest() {
	common.Env("../")

	database.Connect()
	if err := database.CreateTables(); err != nil {
		panic(err)
	}

	if _, err := user.Save(); err != nil {
		panic(err)
	}

	var initialEntries = []model.Entry{{
		Content: "test content 0",
		UserID:  user.ID,
	},
		{
			Content: "test content 0",
			UserID:  user.ID,
		},
	}

	for _, entry := range initialEntries {
		if _, err := entry.Save(); err != nil {
			panic(err)
		}
	}

	suite.db = database.Database

	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.router.POST("/api/entry", AddEntry)
	suite.router.GET("/api/entry", GetAllEntries)
}

func (suite *EntryTestSuite) TearDownTest() {
	db := suite.db

	if _, err := db.Exec("DROP TABLE entries"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("DROP TABLE users"); err != nil {
		panic(err)
	}

	db.Close()
}

func (suite *EntryTestSuite) TestAddEntry_Success() {
	body := `{"content": "test content 0"}`
	req, _ := http.NewRequest("POST", "/api/entry", strings.NewReader(body))

	token, err := common.GenerateJWT(user)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusCreated, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		panic(err)
	}

	data := response["data"].(map[string]interface{})
	suite.Equal("test content 0", data["content"])
}

func (suite *EntryTestSuite) TestAddEntry_Invalid_Input() {
	body := `{}`
	req, _ := http.NewRequest("POST", "/api/entry", strings.NewReader(body))

	token, err := common.GenerateJWT(user)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		panic(err)
	}

	data := response["error"].(string)
	suite.Equal("content is required", data)
}

func (suite *EntryTestSuite) TestGetAllEntries_Success() {
	req, _ := http.NewRequest("GET", "/api/entry", nil)

	token, err := common.GenerateJWT(user)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		panic(err)
	}

	data := response["data"].([]interface{})
	suite.Equal(2, len(data))
}

func TestEntrySuite(t *testing.T) {
	suite.Run(t, new(EntryTestSuite))
}
