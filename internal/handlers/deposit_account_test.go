package handlers_test

import (
	"testing"
)

func TestDepositAccountHandler(t *testing.T) {

	testCases := []struct {
		name                 string
		inputBody            string
		expextedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Success account exists",
			inputBody:            `{"accountId":1,"value":"100.00"}`,
			expextedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name:                 "Success account does not exist",
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

		})
	}
}
