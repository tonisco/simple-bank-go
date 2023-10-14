package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	mockdb "github.com/tonisco/simple-bank-go/db/mock"
	db "github.com/tonisco/simple-bank-go/db/sqlc"
	"github.com/tonisco/simple-bank-go/token"
	"go.uber.org/mock/gomock"
)

func TestRenewAccessToken(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	type Params struct {
		accessUsername   string
		refreshUsername  string
		accessDuration   time.Duration
		refreshDuration  time.Duration
		refreshIsBlocked bool
	}

	testCases := []struct {
		name          string
		getBody       func(token string) gin.H
		params        Params
		setAuth       func(request *http.Request, token string)
		buildStubs    func(store *mockdb.MockStore, session db.Session)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoToken",
			getBody: func(refreshToken string) gin.H {
				return gin.H{}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidToken",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": "abc",
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFoundSession",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "SessionError",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BlockedSession",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: true,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "MismatchUsername",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				newSession := session
				newSession.Username = user2.Username
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(newSession, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "MismatchToken",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				newSession := session
				newSession.RefreshToken = user2.Username
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(newSession, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredSession",
			getBody: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			params: Params{
				accessUsername:   user1.Username,
				refreshUsername:  user1.Username,
				accessDuration:   time.Minute,
				refreshDuration:  time.Minute,
				refreshIsBlocked: false,
			},
			setAuth: func(request *http.Request, token string) {
				authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorizationHeader)
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session) {
				newSession := session
				newSession.ExpiresAt = time.Now().Add(-time.Minute)
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(newSession, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			store := mockdb.NewMockStore(ctr)
			server := newTestServer(t, store)

			accessToken, _ := randomAccessToken(
				t,
				server.tokenMaker,
				authorizationTypeBearer,
				tc.params.accessUsername,
				tc.params.accessDuration,
			)
			refreshToken, _, session := randomRefreshToken(
				t,
				server.tokenMaker,
				authorizationTypeBearer,
				tc.params.refreshUsername,
				tc.params.refreshDuration,
				tc.params.refreshIsBlocked,
			)

			tc.buildStubs(store, session)

			data, err := json.Marshal(tc.getBody(refreshToken))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/tokens/renew_access", bytes.NewReader(data))
			tc.setAuth(request, accessToken)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomAccessToken(t *testing.T, tokenMaker token.Maker, authorizationType, username string, duration time.Duration) (string, *token.Payload) {
	accessToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	return accessToken, payload
}

func randomRefreshToken(
	t *testing.T,
	tokenMaker token.Maker,
	authorizationType,
	username string,
	duration time.Duration,
	isBlocked bool,
) (string, *token.Payload, db.Session) {
	RefreshToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	session := db.Session{
		ID:           payload.ID,
		Username:     payload.Username,
		RefreshToken: RefreshToken,
		UserAgent:    gomock.Any().String(),
		ClientIp:     gomock.Any().String(),
		IsBlocked:    isBlocked,
		ExpiresAt:    time.Now().Add(duration),
		CreatedAt:    time.Now(),
	}

	return RefreshToken, payload, session
}
