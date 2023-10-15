package gapi

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	mockdb "github.com/tonisco/simple-bank-go/db/mock"
	db "github.com/tonisco/simple-bank-go/db/sqlc"
	"github.com/tonisco/simple-bank-go/pb"
	"github.com/tonisco/simple-bank-go/token"
	"github.com/tonisco/simple-bank-go/util"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateUserRoleAPI(t *testing.T) {
	user, _ := randomUser(t)
	user1, _ := randomUser(t)

	user.Role = util.BankerRole

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRoleRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserRoleResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRoleRequest{
				Username: user1.Username,
				Role:     util.BankerRole,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserRoleParams{
					Username: user1.Username,
					Role:     util.BankerRole,
				}
				updatedUser := db.User{
					Username:          user1.Username,
					HashedPassword:    user1.HashedPassword,
					FullName:          user1.FullName,
					Email:             user1.Email,
					PasswordChangedAt: user1.PasswordChangedAt,
					CreatedAt:         user1.CreatedAt,
					IsEmailVerified:   user1.IsEmailVerified,
					Role:              util.BankerRole,
				}
				store.EXPECT().
					UpdateUserRole(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserRoleResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()
				require.Equal(t, user1.Username, updatedUser.Username)
				require.Equal(t, util.BankerRole, updatedUser.Role)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.UpdateUserRoleRequest{
				Username: user.Username,
				Role:     util.BankerRole,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserRole(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserRoleResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InvalidRole",
			req: &pb.UpdateUserRoleRequest{
				Username: user.Username,
				Role:     "super",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserRoleResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.UpdateUserRoleRequest{
				Username: user.Username,
				Role:     util.DepositorRole,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserRoleResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.UpdateUserRoleRequest{
				Username: user.Username,
				Role:     util.BankerRole,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserRoleResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			server := newTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUserRole(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
