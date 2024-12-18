package email

import (
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"mail/api-service/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewEmailRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)

	assert.NotNil(t, router)
	assert.Equal(t, mockEmailUseCase, router.EmailUseCase)
}

func TestEmailRouter_ConfigureEmailRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	router := NewEmailRouter(mockEmailUseCase)
	muxRouter := mux.NewRouter()

	router.ConfigureEmailRouter(muxRouter)

	expectedRoutes := map[string][]string{
		"/email/inbox":       {"GET", "OPTIONS"},
		"/email/sent":        {"GET", "OPTIONS"},
		"/email/{id}":        {"GET", "OPTIONS"},
		"/email/{id}/status": {"PUT", "OPTIONS"},
		"/email/{id}/folder": {"PUT", "OPTIONS"},
		"/getfolder":         {"POST", "OPTIONS"},
		"/draft/send":        {"POST", "OPTIONS"},
		"/status":            {"POST", "OPTIONS"},
		"/getAttachment":     {"POST", "OPTIONS"},
		"/email": {
			"POST", "OPTIONS", "DELETE", "OPTIONS",
		},
		"/folder": {
			"GET", "OPTIONS", "PUT", "OPTIONS", "POST", "OPTIONS", "DELETE", "OPTIONS",
		},
		"/draft": {
			"POST", "OPTIONS", "PUT", "OPTIONS",
		},
		"/attachment": {
			"DELETE", "OPTIONS", "POST", "OPTIONS",
		},
	}

	foundRoutes := make(map[string][]string)
	err := muxRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		if foundRoutes[pathTemplate] == nil {
			foundRoutes[pathTemplate] = []string{}
		}
		foundRoutes[pathTemplate] = append(foundRoutes[pathTemplate], methods...)
		return nil
	})

	assert.NoError(t, err)

	for path, expectedMethods := range expectedRoutes {
		t.Run(path, func(t *testing.T) {
			foundMethods, exists := foundRoutes[path]
			assert.True(t, exists, "Маршрут %s не найден", path)
			assert.ElementsMatch(t, expectedMethods, foundMethods,
				"Методы для маршрута %s не совпадают.\nОжидалось: %v\nПолучено: %v",
				path, expectedMethods, foundMethods)
		})
	}
}

func TestEmailRouter_WithSmtpPop3UseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
	mockSmtpPop3UseCase := mocks.NewMockSmtpPop3Usecase(ctrl)

	router := NewEmailRouter(mockEmailUseCase)
	router.SmtpPop3UseCase = mockSmtpPop3UseCase

	assert.NotNil(t, router)
	assert.Equal(t, mockEmailUseCase, router.EmailUseCase)
	assert.Equal(t, mockSmtpPop3UseCase, router.SmtpPop3UseCase)
}

func TestEmailRouter_Dependencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("проверка зависимостей", func(t *testing.T) {
		mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
		router := NewEmailRouter(mockEmailUseCase)

		assert.Implements(t, (*models.EmailUseCase)(nil), router.EmailUseCase)
	})

	t.Run("проверка опциональных зависимостей", func(t *testing.T) {
		mockEmailUseCase := mocks.NewMockEmailUseCase(ctrl)
		mockSmtpPop3UseCase := mocks.NewMockSmtpPop3Usecase(ctrl)

		router := NewEmailRouter(mockEmailUseCase)
		router.SmtpPop3UseCase = mockSmtpPop3UseCase

		assert.Implements(t, (*models.EmailUseCase)(nil), router.EmailUseCase)
		assert.Implements(t, (*models.SmtpPop3Usecase)(nil), router.SmtpPop3UseCase)
	})
}
