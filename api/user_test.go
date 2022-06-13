package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

func createRandomUser(t *testing.T) db.User {
	hash, err := util.HashPassword("secret")
	require.NoError(t, err)
	return db.User{
		Username:       util.RandomString(6),
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hash,
	}
}

func requireMatchUserBody(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}

func TestGetUserApi(t *testing.T) {

	user := createRandomUser(t)

	testCases := []struct {
		Name          string
		Username      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name:     "Ok",
			Username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchUserBody(t, recorder.Body, user)
			},
		},
		{
			Name:     "NotFound",
			Username: "NoUserUsername",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), "NoUserUsername").Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			Name:     "InvalidCharacterInUsername",
			Username: "With Gap",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), "With Gap").Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		server := NewServer(store)

		tc.buildStubs(store)

		recorder := httptest.NewRecorder()

		url := "/users/" + tc.Username
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, req)

		tc.checkResponse(t, recorder)

	}

}
