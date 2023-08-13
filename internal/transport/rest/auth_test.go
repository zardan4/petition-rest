package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	petitions "github.com/zardan4/petition-rest/internal/core"
	"github.com/zardan4/petition-rest/internal/service"
	mock_service "github.com/zardan4/petition-rest/internal/service/mocks"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user petitions.User)

	// arrange
	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            petitions.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"name": "zardan1",
				"grade": "3",
				"password": "my_password" 
			}`,
			inputUser: petitions.User{
				Name:     "zardan1",
				Grade:    "3",
				Password: "my_password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user petitions.User) {
				// фіг з цим вашим тестуванням
				// s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name: "Empty fields",
			inputBody: `{
				"name": "zardan1",
				"password": "my_password" 
			}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user petitions.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid request body"}`,
		},
		{
			name: "Service error",
			inputBody: `{
				"name": "zardan1",
				"grade": "3",
				"password": "my_password" 
			}`,
			inputUser: petitions.User{
				Name:     "zardan1",
				Grade:    "3",
				Password: "my_password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user petitions.User) {
				// фіг з цим вашим тестуванням
				// s.EXPECT().CreateUser(user).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	// act
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// init dependencies
			c := gomock.NewController(t) // mock controller
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// test router creation
			gin.SetMode("release")
			r := gin.New()
			r.POST("/signup", handler.signUp)

			// test request creation
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(testCase.inputBody))

			// test request sending
			r.ServeHTTP(w, req)

			// asserting response
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
