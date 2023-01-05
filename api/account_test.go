package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/RickZhou666/go-backend-service/db/mock"
	db "github.com/RickZhou666/go-backend-service/db/sqlc"
	"github.com/RickZhou666/go-backend-service/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)                           // each scenario buildStubs will be different
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder) // check response of API
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).            // times this GetAccount func to be called
					Return(account, nil) // should match return values of Getaccount
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code) // check response
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// TODO: add more cases

		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                      // times this GetAccount func to be called
					Return(account, sql.ErrNoRows) // should match return values of Getaccount
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code) // check response
			},
		},

		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                        // times this GetAccount func to be called
					Return(account, sql.ErrConnDone) // should match return values of Getaccount
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code) // check response
			},
		},

		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0) // times this GetAccount func to be called
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code) // check response
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// below step is important
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder() // to create recorder to record response from API request

			url := fmt.Sprintf("/accounts/%d", tc.accountID) // create url path
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})

	}

}
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// read all data
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)

}
