package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domainErrors "github.com/Dorrrke/shortener-url-cource/internal/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func MiddlewareForTest(ctx *gin.Context) {
	ctx.Set("userID", "test-user")
	ctx.Next()
}

func TestSaveURL(t *testing.T) {
	linksHandler := LinkHandler{domain: "test.org"}
	r := gin.New()
	r.POST("/save-url", MiddlewareForTest, linksHandler.SaveURL)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		statusCode int
		body       string
		shortLink  string
		err        error
	}
	type test struct {
		name    string
		method  string
		fullUrl string
		reqJson string
		want    want
	}

	tests := []test{
		{
			name:    "Case 1: Default call",
			method:  http.MethodPost,
			fullUrl: "https://vk.com/test",
			reqJson: `{"url":"https://vk.com/test"}`,
			want: want{
				statusCode: http.StatusOK,
				body:       `{"short_link":"http://test.org/link/shortlink"}`,
				shortLink:  "shortlink",
				err:        nil,
			},
		},
		{
			name:    "Case 2: Bad request",
			method:  http.MethodPost,
			reqJson: "url:123",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       `{"error":"invalid character 'u' looking for beginning of value"}`,
			},
		},
		{
			name:    "Case 3: Link Alredy Exists",
			method:  http.MethodPost,
			fullUrl: "https://vk.com/test",
			reqJson: `{"url":"https://vk.com/test"}`,
			want: want{
				statusCode: http.StatusConflict,
				body:       `{"error":"Вы уже сокращали эту ссылку"}`,
				shortLink:  "",
				err:        domainErrors.ErrLinkAlreadyExists,
			},
		},
		{
			name:    "Case 4: Unknown error",
			method:  http.MethodPost,
			fullUrl: "https://vk.com/test",
			reqJson: `{"url":"https://vk.com/test"}`,
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       `{"error":"Unknown error"}`,
				shortLink:  "",
				err:        errors.New("Unknown error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := NewMockLinkService(t)
			mockService.
				On("SaveURL", tc.fullUrl, "test-user").
				Return(tc.want.shortLink, tc.want.err).Maybe()

			linksHandler.service = mockService

			req := resty.New().R()

			req.Method = tc.method
			req.URL = httpSrv.URL + "/save-url"

			req.Body = []byte(tc.reqJson)

			response, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tc.want.statusCode, response.StatusCode())
			assert.Equal(t, tc.want.body, string(response.Body()))
		})
	}
}

func TestGet(t *testing.T) {
	linksHandler := LinkHandler{domain: "test.org"}
	type want struct {
		statusCode int
		link       string
		err        error
	}
	type test struct {
		name    string
		method  string
		request string
		want    want
	}

	tests := []test{
		{
			name:    "Case 1: Default Call",
			method:  http.MethodGet,
			request: "/link/short-link",
			want: want{
				statusCode: http.StatusFound,
				link:       "http://vk.com/test",
			},
		},
		{
			name:    "Case 2: Link not found",
			method:  http.MethodGet,
			request: "/link/new-link",
			want: want{
				statusCode: http.StatusNotFound,
				err:        domainErrors.ErrLinkNotFound,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := NewMockLinkService(t)
			mockService.
				On("Get", strings.TrimPrefix(tc.request, "/link/")).
				Return(tc.want.link, tc.want.err).Maybe()

			linksHandler.service = mockService

			router := gin.New()
			router.GET("/link/:id", linksHandler.Get)

			req := httptest.NewRequest(
				tc.method,
				tc.request,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.want.statusCode, rec.Code)
			if tc.want.err != nil {
				assert.Contains(t, rec.Body.String(), tc.want.err.Error())
			} else {
				assert.Equal(t, tc.want.link, rec.Header().Get("Location"))
			}
		})
	}

}
