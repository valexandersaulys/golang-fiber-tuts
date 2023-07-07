package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
)

func TestHelloRouteUno(t *testing.T) {
	// http.Request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Custom-Header", "hi")

	app := CreateApp()
	// http.Response
	resp, _ := app.Test(req)

	// Do something with results:
	if resp.StatusCode == fiber.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // => Hello, World!
	} else {
		err := errors.New("ASDFASDFASD")
		log.Fatal(err)
	}

}

func TestRoutesDos(t *testing.T) {
	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		reqType      string
		reqJSON      map[string]string
		expectedCode int // expected HTTP status code
		expectedBody map[string]string
	}{
		// First test case
		{
			description:  "get HTTP status 200",
			route:        "/",
			reqType:      "GET",
			expectedCode: 200,
		},
		// Second test case
		{
			description:  "get HTTP status 404, when route is not exists",
			route:        "/not-found",
			reqType:      "GET",
			expectedCode: 404,
		},
		{
			description: "Can send and receive a Person struct",
			route:       "/person/",
			reqType:     "POST",
			reqJSON: map[string]string{
				"name":  "vincent",
				"email": "vincent@saulys.me",
			},
			expectedCode: 200,
			expectedBody: map[string]string{
				"name":  "vincent",
				"email": "vincent@saulys.me",
			},
		},
	}

	app := CreateApp()

	for _, test := range tests {
		if test.reqType == "GET" {
			// is nil payload?
			req := httptest.NewRequest(test.reqType, test.route, nil)
			resp, _ := app.Test(req, 100)
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		} else if test.reqType == "POST" {
			jsonStr, _ := json.Marshal(test.reqJSON)
			fmt.Println(string(jsonStr))
			req := httptest.NewRequest(test.reqType, test.route, bytes.NewReader(jsonStr))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8") // REQUIRED
			resp, _ := app.Test(req, 100)
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
			var decodedJson map[string]string
			body, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &decodedJson)
			assert.Equalf(t, test.expectedBody, decodedJson, test.description)
		}
	}
}
