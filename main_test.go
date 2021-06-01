package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdd(t *testing.T) {
	// a.Initialize(
	// 	os.Getenv("USER"),
	// 	os.Getenv("NAME"),
	// 	os.Getenv("PASSWORD"))
}

func TestRootEndpoint(t *testing.T) {
	request, _ := http.NewRequest("GET", "/tasks/list", nil)
	response := httptest.NewRecorder()
	RouteHandler().ServeHTTP(response, request)
	fmt.Println("response -->", response)
	//assert.Equal(t, )
}

