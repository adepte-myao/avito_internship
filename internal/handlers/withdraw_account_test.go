package handlers_test

import (
	"testing"
)

func TestWithdrawAccountHandler(t *testing.T) {

	testCases := []struct {
		name                 string
		inputBody            string
		expextedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Success",
			inputBody:            `{"accountId":1,"value":"100.00"}`,
			expextedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name:                 "Invalid request body",
			inputBody:            `{"accountId":"1","value":"100.00"}`,
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"invalid request body\"}\n",
		},
		{
			name:                 "Account does not exist",
			inputBody:            `{"accountId":1,"value":"100.00"}`,
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"account does not exist\"}\n",
		},
		{
			name:                 "Not enough money",
			inputBody:            `{"accountId":1,"value":"100.00"}`,
			expextedStatusCode:   400,
			expectedResponseBody: "{\"reason\":\"not enough money\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

		})
	}
}
