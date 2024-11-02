package resolvers

import (
	"app/lib"
	models "app/models/generated"
	"app/services"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TestUserResolverSuite struct {
	WithDBSuite
}

var (
	testUserGraphQLServer *handler.Server
)

func (s *TestUserResolverSuite) SetupTest() {
	s.SetDBCon()

	authService := services.NewAuthService(DBCon)

	// NOTE: テスト対象のサーバのハンドラを設定
	testUserGraphQLServer = lib.GetGraphQLServer(authService)
}

func (s *TestUserResolverSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestUserResolverSuite) TestSignUp() {
	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            signUp(input: {
                Name: "test name 1",
                Email: "test@example.com",
                Password: "password"
            }) {
                id,
                name,
                email,
                nameAndEmail
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testUserGraphQLServer.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]interface{})
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Contains(s.T(), responseBody["data"], "signUp")

	// NOTE: ユーザが作成されていることを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	assert.True(s.T(), isExistUser)
}

func (s *TestUserResolverSuite) TestSignUp_ValidationError() {
	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            signUp(input: {
                Name: "test name 1",
                Email: "",
                Password: "password"
            }) {
                id,
                name,
                email,
                nameAndEmail
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testUserGraphQLServer.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(400), responseBody["errors"][0]["extensions"]["code"])
	assert.Contains(s.T(), responseBody["errors"][0]["extensions"]["error"], "Email")
	log.Println(responseBody)

	// NOTE: ユーザが作成されていないことを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistUser)
}

func TestUserResolver(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestUserResolverSuite))
}