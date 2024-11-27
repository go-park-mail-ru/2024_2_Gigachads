package auth

import (
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestConfigureAuthRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUseCase := mocks.NewMockAuthUseCase(ctrl)
	router := mux.NewRouter()
	authRouter := NewAuthRouter(mockAuthUseCase)

	authRouter.ConfigureAuthRouter(router)

	routes := []struct {
		path   string
		method string
	}{
		{"/signup", http.MethodPost},
		{"/login", http.MethodPost},
		{"/logout", http.MethodDelete},
	}

	for _, route := range routes {
		match := &mux.RouteMatch{}
		req, err := http.NewRequest(route.method, route.path, nil)
		assert.NoError(t, err)

		matched := router.Match(req, match)
		if !matched {
			t.Errorf("Route %s %s not found", route.method, route.path)
		}
		assert.True(t, matched)
	}
}

func TestNewAuthRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUseCase := mocks.NewMockAuthUseCase(ctrl)
	router := NewAuthRouter(mockAuthUseCase)

	assert.NotNil(t, router)
	assert.Equal(t, mockAuthUseCase, router.AuthUseCase)
}
