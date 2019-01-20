package main_test

import (
	"encoding/json"
	"net/http/httptest"
	"bytes"
	"net/http"
    "os"
    "testing"

	"."
)

type Website struct {
    URL string
    Status bool
}

var a main.App

func TestMain(m *testing.M) {
    a = main.App{}
    a.Initialize()


    code := m.Run()


    os.Exit(code)
}

// POST /api/healthcheck test add a website
func TestAddWebSite(t *testing.T) {
	w := []byte(`{"URL":"http://google.com"}`)
	req, _ := http.NewRequest("POST", "/api/healthcheck", bytes.NewBuffer(w))
    response := executeRequest(req)
	
	// t.Logf("Response code: %v \n", response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	
	// check response URL
	if m["URL"] != "http://google.com" {
		t.Errorf("Expected the response url to be http://google.com. Got %v\n", m["URL"])
	}

	// check response status
	if m["Status"] != true {
		t.Errorf("Expect the website status to be true. Got %v\n", m["Status"])
	}
}

// GET /api/healthcheck test get all websites
func TestGetAllWebsites(t *testing.T) {
	w := []byte(`{"URL":"http://baidu.com"}`)
	req,_ := http.NewRequest("POST", "/api/healthcheck", bytes.NewBuffer(w))
	response := executeRequest(req)

	var ws []Website
	req, _ = http.NewRequest("GET", "/api/healthcheck", nil)
	response = executeRequest(req)

	json.Unmarshal(response.Body.Bytes(), &ws)
	
	if len(ws) != 2 {
		t.Errorf("Expect 2 websites in response. Got %v\n", len(ws))
	}
}

// DELETE /api/healthcheck Test delete website by url
func TestDeleteWebsite(t *testing.T) {
	w := []byte(`{"URL":"http://baidu.com"}`)
	req,_ := http.NewRequest("DELETE", "/api/healthcheck", bytes.NewBuffer(w))
	response := executeRequest(req)

	// Expect DELETE response code 200
	if response.Code != http.StatusOK {
		t.Errorf("Expect status 200. Got %v\n", response.Code)
	}

	var ws []Website
	req, _ = http.NewRequest("GET", "/api/healthcheck", nil)
	response = executeRequest(req)

	json.Unmarshal(response.Body.Bytes(), &ws)
	
	// Expect remaining website does not have http://baidu.com
	if ws[0].URL == "http://baidu.com" {
		t.Errorf("Expect http://baidu.com to be removed.\n")
	}
}








func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}