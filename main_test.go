package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestNewApp(t *testing.T) {
	// Act
	app := NewApp()

	// Assert
	if app.Server.Addr != ":9876" {
		t.Fatal("port doesn't match")
	}

	if app.Router == nil {
		t.Fatal("routes should not be nil")
	}
}

func TestAppRun(t *testing.T) {
	// Arrange
	app := NewApp()
	app.Server.Addr = ""

	// Tranfer logger output to buf
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	// Act
	app.Run()

	// Assert
	t.Log(buf.String())
	re := regexp.MustCompile(".* Cannot start server on 9876")
	if !re.Match(buf.Bytes()) {
		t.Fatal("Server should not start.")
	}
}

func TestSetRoutes(t *testing.T) {
	// Arrange
	app := &App{
		Router: http.NewServeMux(),
		Server: &http.Server{},
	}

	// Act
	app.setRoutes()

	// Assert
	if app.Server.Handler == nil || app.Router == nil {
		t.Fatal("Handler should not be nil")
	}
}

func TestIndexRouter(t *testing.T) {
	// Arrange
	ts := httptest.NewServer(IndexRouter())
	defer ts.Close()

	// Act
	resp, _ := http.Get(ts.URL)
	b, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Assert
	if string(b) != "Index" {
		t.Fatal("response body doen't match")
	}
}

func TestLogMiddleWare(t *testing.T) {
	// Arrange
	fakeHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test")
	}
	ts := httptest.NewServer(LogMiddleware(fakeHandlerFunc))
	defer ts.Close()

	// Tranfer logger output to buf
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	// Act
	http.Get(fmt.Sprintf("%s/", ts.URL))

	// Assert
	t.Logf("buf: %s\n", buf.String())
	re := regexp.MustCompile("http: [\\d]{4}/[\\d]{2}/[\\d]{2} [\\d]{2}:[\\d]{2}:[\\d]{2} (.)*")
	if !re.Match(buf.Bytes()) {
		t.Fatal("Log format doesn't match")
	}
}

func TestLogMiddleWareFailed(t *testing.T) {
	// Arrange
	fakeHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		panic("failed!")

	}
	ts := httptest.NewServer(LogMiddleware(fakeHandlerFunc))
	defer ts.Close()

	// Tranfer logger output to buf
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	// Act
	http.Get(fmt.Sprintf("%s/", ts.URL))

	// Assert
	t.Logf("buf: %s\n", buf.String())
	re := regexp.MustCompile("http: [\\d]{4}/[\\d]{2}/[\\d]{2} [\\d]{2}:[\\d]{2}:[\\d]{2} (.*) \\[error: (.)*\\]")
	if !re.Match(buf.Bytes()) {
		t.Fatal("Log format doesn't match")
	}
}
