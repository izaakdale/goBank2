package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/izaakdale/goBank2/db/mock"
	db "github.com/izaakdale/goBank2/db/sqlc"
	"github.com/izaakdale/goBank2/token"
	"github.com/izaakdale/goBank2/util"
	"github.com/stretchr/testify/require"
)

func TestTransferApi(t *testing.T) {

	user1, _ := createRandomUser(t)
	user2, _ := createRandomUser(t)
	user3, _ := createRandomUser(t)

	account1 := createRandomAccount(user1.Username)
	account2 := createRandomAccount(user2.Username)
	account3 := createRandomAccount(user3.Username)

	account1.Currency = util.USD
	account2.Currency = util.USD
	account3.Currency = util.GBP

	amount := float64(10)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker *token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				req := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(req)).Times(1)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker *token.Maker) {
				addAuthHeader(t, request, *tokenMaker, AuthorizationTypeBearer, user1.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "To Account Incorrect currency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          10,
				"currency":        util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker *token.Maker) {
				addAuthHeader(t, request, *tokenMaker, AuthorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)

				req := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account3.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(req)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "From Account Incorrect currency",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account1.ID,
				"amount":          10,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(0)

				req := db.TransferTxParams{
					FromAccountID: account3.ID,
					ToAccountID:   account1.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(req)).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker *token.Maker) {
				addAuthHeader(t, request, *tokenMaker, AuthorizationTypeBearer, user3.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "From account not found",
			body: gin.H{
				"from_account_id": int64(1001),
				"to_account_id":   account2.ID,
				"amount":          10,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(int64(1001))).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)

				req := db.TransferTxParams{
					FromAccountID: int64(1001),
					ToAccountID:   account2.ID,
					Amount:        10,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(req)).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker *token.Maker) {
				addAuthHeader(t, request, *tokenMaker, AuthorizationTypeBearer, user1.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "To account not found",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   int64(1001),
				"amount":          10,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(int64(1001))).Times(1).Return(db.Account{}, sql.ErrNoRows)

				req := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   int64(1001),
					Amount:        10,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(req)).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker *token.Maker) {
				addAuthHeader(t, request, *tokenMaker, AuthorizationTypeBearer, user1.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "NoAuth",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0).Return(account2, nil)

				req := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(req)).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker *token.Maker) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			// start http server and send request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, &server.tokenMaker)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
