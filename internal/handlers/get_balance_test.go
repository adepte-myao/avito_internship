package handlers_test

import (
	"testing"
)

func TestGetBalanceHandler(t *testing.T) {

	testCases := []struct {
		name                 string
		inputBody            string
		expextedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Success without fractional part",
			inputBody:            `{"accountId":1}`,
			expextedStatusCode:   200,
			expectedResponseBody: "{\"accountId\":1,\"balance\":\"100\"}\n",
		},
		{
			name:                 "Success with fractional part",
			inputBody:            `{"accountId":1}`,
			expextedStatusCode:   200,
			expectedResponseBody: "{\"accountId\":1,\"balance\":\"100.01\"}\n",
		},
		{
			name:                 "Invalid request body",
			inputBody:            `{"accountId":"1"}`,
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"invalid request body\"}\n",
		},
		{
			name:                 "Account does not exist",
			inputBody:            `{"accountId":1}`,
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"account with given ID does not exist\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

		})
	}
}
