package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/mrohadi/simplebank/db/mock"
	db "github.com/mrohadi/simplebank/db/sqlc"
	"github.com/mrohadi/simplebank/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCase := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccount(t *testing.T) {
	account := randomAccount()
	reqBody := randomReqBody(&account)
	param := db.CreateAccountParams{
		Owner:    reqBody.Owner,
		Currency: reqBody.Currency,
		Balance:  0,
	}

	testCase := []struct {
		name          string
		param         db.CreateAccountParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "Created",
			param: param,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(param)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotAccount db.Account
				err = json.Unmarshal(data, &gotAccount)
				require.NoError(t, err)
				require.Equal(t, reqBody.Currency, gotAccount.Currency)
				require.Equal(t, reqBody.Owner, gotAccount.Owner)
			},
		},
		{
			name:  "InternalServerError",
			param: param,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(param)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "BadRequest",
			param: db.CreateAccountParams{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			//build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts"

			body, err := json.Marshal(tc.param)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)
			server.Start("")
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccount(t *testing.T) {
	var listAccount []db.Account
	for i := 0; i < 5; i++ {
		listAccount = append(listAccount, randomAccount())
	}

	queryParam := listAccountRequest{PageID: 1, PageSize: 5}
	arg := db.ListAccountsParams{
		Limit:  queryParam.PageSize,
		Offset: (queryParam.PageID - 1) * queryParam.PageSize,
	}
	testCase := []struct {
		name          string
		param         listAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "BadRequest",
			param: listAccountRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "InternalServerError",
			param: queryParam,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "OK",
			param: queryParam,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(listAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotAccounts []db.Account
				err = json.Unmarshal(data, &gotAccounts)
				require.NoError(t, err)
				require.Equal(t, int(queryParam.PageSize), len(gotAccounts))
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stub
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tc.param.PageID, tc.param.PageSize)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomReqBody(a *db.Account) createAccountParams {
	return createAccountParams{
		Owner:    a.Owner,
		Currency: a.Currency,
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
