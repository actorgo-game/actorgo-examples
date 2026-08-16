package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestControllerIndex(t *testing.T) {
	engine := gin.New()
	controller := new(Test1Controller)
	controller.PreInit(nil, engine)
	controller.Init()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "this is index") {
		t.Fatalf("unexpected response body %q", response.Body.String())
	}
}
