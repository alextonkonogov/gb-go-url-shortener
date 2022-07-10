package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type Config struct {
	Services Services `yaml:"services"`
}

type Services struct {
	URLShortener URLShortener `yaml:"url_shortener"`
}

type URLShortener struct {
	Environment Environment `yaml:"environment"`
}

type Environment struct {
	DBConnectionString string `yaml:"DB_CONNECTION_STRING"`
	AppIP              string `yaml:"APP_IP"`
	AppPort            string `yaml:"APP_PORT"`
}

var (
	appIP,
	appPort,
	requestURL,
	pgStr,
	short,
	admin string
)

func Test(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../../../docker-compose.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	c := Config{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	appIP = c.Services.URLShortener.Environment.AppIP
	appPort = c.Services.URLShortener.Environment.AppPort
	pgStr = c.Services.URLShortener.Environment.DBConnectionString
	requestURL = "http://" + appIP + ":" + appPort

	testCreateURL(t)
	testReadURL(t)
	testReadAdmin(t)
}

var CreateURLTests = []struct {
	method            string
	bodyStr           string
	contentType       string
	wantedStatus      int
	wantedContentType string
	waited            bool
}{
	{
		method:            "GET",
		bodyStr:           "",
		contentType:       "",
		wantedStatus:      http.StatusOK,
		wantedContentType: "text/html; charset=utf-8",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           "",
		contentType:       "",
		wantedStatus:      http.StatusBadRequest,
		wantedContentType: "application/json; charset=utf-8",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           `{"lorem":"ipsum"}`,
		contentType:       "",
		wantedStatus:      http.StatusBadRequest,
		wantedContentType: "application/json; charset=utf-8",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           `{"long":"lorem ipsum"}`,
		contentType:       "",
		wantedStatus:      http.StatusBadRequest,
		wantedContentType: "application/json; charset=utf-8",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           `{"long":"https://gbcdn.mrgcdn.ru/uploads/asset/3001858/attachment/c3640e219eb26045352728efea6d443e.pdf"}`,
		contentType:       "",
		wantedStatus:      http.StatusBadRequest,
		wantedContentType: "application/json; charset=utf-8",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           `{"long":"https://gbcdn.mrgcdn.ru/uploads/asset/3001858/attachment/c3640e219eb26045352728efea6d443e.pdf"}`,
		contentType:       "application/json",
		wantedStatus:      http.StatusOK,
		wantedContentType: "application/json; charset=utf-8",
		waited:            true,
	},
}

func testCreateURL(t *testing.T) {
	for i, tt := range CreateURLTests {
		fmt.Println("test", i, tt)

		jsonBody := []byte(tt.bodyStr)
		bodyReader := bytes.NewReader(jsonBody)
		req, err := http.NewRequest(
			tt.method,
			requestURL+"/s/create",
			bodyReader,
		)
		req.Header.Set("Content-Type", tt.contentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("failed to do http request: %s", err.Error())
		}

		if status := res.StatusCode; status != tt.wantedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.wantedStatus)
		}

		if res.Header.Get("Content-Type") != tt.wantedContentType {
			t.Errorf("handler returned unexpected content type: got %v want %v",
				res.Header.Get("Content-Type"), tt.wantedContentType)
		}
		if tt.waited {
			body, _ := ioutil.ReadAll(res.Body)
			url := url.URL{}
			_ = json.Unmarshal([]byte(string(body)), &url)
			short = url.Short
			admin = url.Admin
		}

		res.Body.Close()
	}
}

var ReadURLTests = []struct {
	method       string
	pathParam    string
	wantedStatus int
	waited       bool
}{
	{
		method:       "POST",
		pathParam:    "",
		wantedStatus: http.StatusNotFound,
		waited:       false,
	},
	{
		method:       "GET",
		pathParam:    "",
		wantedStatus: http.StatusNotFound,
		waited:       false,
	},
	{
		method:       "GET",
		pathParam:    "lorem+ipsum",
		wantedStatus: http.StatusOK,
		waited:       false,
	},
	{
		method:       "GET",
		pathParam:    "",
		wantedStatus: http.StatusOK,
		waited:       true,
	},
}

func testReadURL(t *testing.T) {
	for i, tt := range ReadURLTests {
		fmt.Println("test", i, tt)

		if tt.waited {
			tt.pathParam = short
		}

		req, err := http.NewRequest(
			tt.method,
			requestURL+"/s/"+tt.pathParam,
			nil,
		)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("failed to do http request: %s", err.Error())
		}

		if status := res.StatusCode; status != tt.wantedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.wantedStatus)
		}
	}
}

var ReadAdminTests = []struct {
	method            string
	bodyStr           string
	contentType       string
	wantedStatus      int
	wantedContentType string
	waited            bool
}{
	{
		method:            "GET",
		bodyStr:           "",
		contentType:       "",
		wantedStatus:      http.StatusMethodNotAllowed,
		wantedContentType: "",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           "",
		contentType:       "",
		wantedStatus:      http.StatusBadRequest,
		wantedContentType: "application/json; charset=utf-8",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           `{"admin":"lorem+ipsum"}`,
		contentType:       "",
		wantedStatus:      http.StatusBadRequest,
		wantedContentType: "application/json; charset=utf-8",
		waited:            false,
	},
	{
		method:            "POST",
		bodyStr:           `{"admin":"%s"}`,
		contentType:       "application/json",
		wantedStatus:      http.StatusOK,
		wantedContentType: "application/json; charset=utf-8",
		waited:            true,
	},
}

func testReadAdmin(t *testing.T) {
	for i, tt := range ReadAdminTests {
		fmt.Println("test", i, tt)

		if tt.waited {
			tt.bodyStr = fmt.Sprintf(tt.bodyStr, admin)
		}

		jsonBody := []byte(tt.bodyStr)
		bodyReader := bytes.NewReader(jsonBody)
		req, err := http.NewRequest(
			tt.method,
			requestURL+"/a",
			bodyReader,
		)
		req.Header.Set("Content-Type", tt.contentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("failed to do http request: %s", err.Error())
		}

		if status := res.StatusCode; status != tt.wantedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.wantedStatus)
		}

		if res.Header.Get("Content-Type") != tt.wantedContentType {
			t.Errorf("handler returned unexpected content type: got %v want %v",
				res.Header.Get("Content-Type"), tt.wantedContentType)
		}

		res.Body.Close()
	}
}
