package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seferin-x/token"
	"github.com/stretchr/testify/require"
)

type Server struct {
	router     *gin.Engine
	tokenMaker *token.TokenMaker
}

func NewServer(t *testing.T) *Server {
	tm, err := token.NewTokenMaker("12345678901234567890123456789012")
	require.NoError(t, err)
	require.NotNil(t, tm)

	r := gin.Default()

	return &Server{
		router:     r,
		tokenMaker: tm,
	}
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.TokenMaker,
) {
	payloadIn := make(map[string]interface{})

	payloadIn["test"] = "test"

	token, payloadOut, err := tokenMaker.CreateToken(payloadIn)
	require.NoError(t, err)
	require.NotEmpty(t, payloadOut)

	authorizationHeader := token
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "WrongTokenKey",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {

				tm, err := token.NewTokenMaker("11111111112222222222333333333312")
				require.NoError(t, err)
				require.NotNil(t, tm)
				addAuthorization(t, request, *tm)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(t)
			authPath := "/"
			server.router.GET(
				authPath,
				ginAuthMiddleware(*server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, *server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
