package web_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"webook/internal/integration_test/startup"
	"webook/internal/integration_test/util"
	"webook/internal/repository/dao"
	"webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ArticleTestSuite struct {
	suite.Suite

	server    *gin.Engine
	db        *gorm.DB
	user      dao.User
	userAgent string
	userToken string
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}

func (s *ArticleTestSuite) SetupSuite() {
	s.server = startup.InitApiServer()
	s.db = ioc.InitDB()
	util.InitTables(s.db)

	// create user
	util.TruncateTables(s.T(), s.db, "draft_articles", "users")
	s.userAgent = "test-agent"
	var err error
	s.user, s.userToken, err = util.CreateUserAndGetToken(s.T(), s.db, s.userAgent)
	if err != nil {
		assert.FailNow(s.T(), "failed to create user", err)
	}
}

func (s *ArticleTestSuite) TearDownTestSuite() {
	util.TruncateTables(s.T(), s.db, "draft_articles", "users")
}

func (s *ArticleTestSuite) BeforeTest() {
	util.TruncateTables(s.T(), s.db, "draft_articles")
}

func (s *ArticleTestSuite) TestCreateArticle() {
	testCases := []struct {
		name    string
		method  string
		path    string
		token   string
		reqBody string

		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Success Case: create article",
			method:       "POST",
			path:         "/articles",
			token:        s.userToken,
			reqBody:      `{"title": "test", "content": "test"}`,
			wantStatus:   200,
			wantResponse: `{"code":"200","message": "success","data":{"id": 1}}`,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.reqBody))
			req.Header.Set("Authorization", "Bearer "+tc.token)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", s.userAgent)
			w := httptest.NewRecorder()
			s.server.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tc.wantStatus, resp.StatusCode)
			assert.JSONEq(t, tc.wantResponse, string(body))
		})
	}

}
