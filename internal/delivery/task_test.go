package delivery

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"task-queue/internal/models"
	"task-queue/internal/service"
	mock_service "task-queue/internal/service/mocks"
	"testing"
)

func TestDelivery_createTask(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTask, input string)
	testTable := []struct {
		name           string
		inputBody      string
		inputData      string
		mockBehavior   mockBehavior
		expStatusCode  int
		expRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"input_data":"test"}"`,
			inputData: "test",
			mockBehavior: func(s *mock_service.MockTask, input string) {
				s.EXPECT().Create(input).Return(uint32(1), nil)
			},
			expStatusCode:  200,
			expRequestBody: `{"task_id":1}`,
		},
		{
			name:           "wrong input body",
			inputBody:      "qwe",
			inputData:      "",
			mockBehavior:   func(s *mock_service.MockTask, input string) {},
			expStatusCode:  400,
			expRequestBody: `{"message":"Неверное содержание json"}`,
		},
		{
			name:      "wrong input data",
			inputBody: `{"input_data":""}`,
			inputData: "",
			mockBehavior: func(s *mock_service.MockTask, input string) {
				s.EXPECT().Create(input).Return(uint32(0), fmt.Errorf("Неверное содержание ввода: %s", input))
			},
			expStatusCode:  400,
			expRequestBody: `{"message":"Неверное содержание ввода: "}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			c.Finish()

			task := mock_service.NewMockTask(c)
			testCase.mockBehavior(task, testCase.inputData)

			services := &service.Service{Task: task}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/tasks", handler.createTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testCase.expStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expRequestBody)
		})
	}
}

func TestDelivery_getTask(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTask, id uint32)
	testTable := []struct {
		name           string
		idParam        string
		id             uint32
		mockBehavior   mockBehavior
		expStatusCode  int
		expRequestBody string
	}{
		{
			name:    "ok",
			idParam: "1",
			id:      1,
			mockBehavior: func(s *mock_service.MockTask, id uint32) {
				s.EXPECT().Get(id).Return(models.Task{Status: "completed", Result: "Результат: test"}, nil)
			},
			expStatusCode:  200,
			expRequestBody: `{"status":"completed","result":"Результат: test"}`,
		},
		{
			name:           "invalid ID",
			idParam:        "qwe",
			id:             0,
			mockBehavior:   func(s *mock_service.MockTask, id uint32) {},
			expStatusCode:  400,
			expRequestBody: `{"message":"Неверный парамметр id"}`,
		},
		{
			name:    "task not found",
			idParam: "2",
			id:      2,
			mockBehavior: func(s *mock_service.MockTask, id uint32) {
				s.EXPECT().Get(id).Return(models.Task{}, fmt.Errorf("Задача с id %d не найдена", id))
			},
			expStatusCode:  500,
			expRequestBody: fmt.Sprintf(`{"message":"Задача с id %d не найдена"}`, 2),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			task := mock_service.NewMockTask(c)
			testCase.mockBehavior(task, testCase.id)

			services := &service.Service{Task: task}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/tasks/:id", handler.getTask)

			w := httptest.NewRecorder()

			var req *http.Request
			var err error
			if testCase.id > 0 {
				req, err = http.NewRequest("GET", fmt.Sprintf("/tasks/%s", testCase.idParam), nil)
				assert.NoError(t, err)
			} else {
				req, err = http.NewRequest("GET", fmt.Sprintf("/tasks/%s", testCase.idParam), nil)
				assert.NoError(t, err)
			}

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testCase.expStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expRequestBody)
		})
	}
}
