package resolvers

import (
	"app/lib"
	models "app/models/generated"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TestTodoResolverSuite struct {
	WithDBSuite
}

var (
	testTodoGraphQLServerHandler http.Handler
)

func (s *TestTodoResolverSuite) SetupTest() {
	s.SetDBCon()

	// NOTE: テスト対象のサーバのハンドラを設定
	testTodoGraphQLServerHandler = lib.GetGraphQLHttpHandler(DBCon)
}

func (s *TestTodoResolverSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestTodoResolverSuite) TestCreateTodo_Unauthorized() {
	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            createTodo(input: {
                title: "test title 1",
                content: "",
            }) {
                id,
                title,
                content,
                createdAt,
				updatedAt
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(401), responseBody["errors"][0]["extensions"]["code"])

	// NOTE: Todoリストが作成されていないことを確認
	isExistTodo, _ := models.Todos(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistTodo)
}

func (s *TestTodoResolverSuite) TestCreateTodo() {
	s.signIn()

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            createTodo(input: {
                title: "test title 1",
                content: "",
            }) {
                id,
                title,
                content,
                createdAt,
				updatedAt
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]interface{})
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Contains(s.T(), responseBody["data"], "createTodo")

	// NOTE: Todoリストが作成されていることを確認
	isExistTodo, _ := models.Todos(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.True(s.T(), isExistTodo)
}

func (s *TestTodoResolverSuite) TestCreateTodo_ValidationError() {
	s.signIn()

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            createTodo(input: {
                title: "",
                content: "",
            }) {
                id,
                title,
                content,
                createdAt,
				updatedAt
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(400), responseBody["errors"][0]["extensions"]["code"])

	// NOTE: Todoリストが作成されていないことを確認
	isExistTodo, _ := models.Todos(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistTodo)
}

func TestTodoResolver(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestTodoResolverSuite))
}
