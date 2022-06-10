package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/izaakdale/goBank2/db/mock"
	db "github.com/izaakdale/goBank2/db/sqlc"
	"github.com/izaakdale/goBank2/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// create the stub that returns the account when account id is requested
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(account, nil).Times(1)

	// start http server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireMatchAccountBody(t, recorder.Body, account)
}
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(100, 1000),
		Owner:    util.RandomName(),
		Balance:  float64(util.RandomBalance()),
		Currency: util.RandomCurrency(),
	}
}

func requireMatchAccountBody(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
