package goscrapyd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewScrapyd(t *testing.T) {

	var host = "localhost:32785"

	s := NewScrapyd(host)

	if host != s.Host {
		t.Errorf("Expected %q, got %q", host, s.Host)
	}
}

func TestBuildUrl(t *testing.T) {

	var host = "http://localhost:32785"

	s := NewScrapyd(host)

	var expected_url = fmt.Sprintf("%v/%v", host, "addversion.json")

	if path, _ := s.buildUrl(add_version_endpoint); strings.Compare(expected_url, path) != 0 {
		t.Errorf("Expected %q, got %q", expected_url, path)
	}
}

func TestHttpGet(t *testing.T) {

	// Instantiate a testing server
	mux, server := testServer()
	defer server.Close()

	// Configures the server endpoint and provides a response.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "foobar")
	})

	su, _ := url.Parse(server.URL)
	s := NewScrapyd(fmt.Sprintf("%v://%v", su.Scheme, su.Host))

	params := make(map[string]string)

	u, _ := url.Parse(server.URL)
	resp, err := s.get(u.String(), params)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Unable to request the server, got a %q response code, want %q", resp.StatusCode, "200")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("reading reponse body: %v, want %q", err, b)
	}

	if string(b) != "foobar" {
		t.Errorf("request body mismatch: got %q, want %q", string(b), "foobar")
	}
}

func TestLitsJobs(t *testing.T) {

	// Instantiate a testing server
	mux, server := testServer()
	defer server.Close()

	// Generates a fake body response.
	var sr = &ScrapydResponse{Status: "foo"}

	// Gets a json value from the body response.
	srjson, err := json.Marshal(sr)
	if err != nil {
		fmt.Println(err)
	}

	// Configures the server endpoint and provides a response.
	mux.HandleFunc("/listjobs.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(srjson))
	})

	// Sets the endpoint to the testing server.
	server.URL = fmt.Sprintf("%s/%s", server.URL, "foo")

	su, _ := url.Parse(server.URL)
	s := NewScrapyd(fmt.Sprintf("%v://%v", su.Scheme, su.Host))

	resp, _ := s.ListJobs("foo")

	if strings.Compare(sr.Status, resp.Status) != 0 {
		t.Errorf("request body mismatch: got %q, want %q", resp, *sr)
	}
}

// testServer returns an ServeMux, and Server. The caller must close the test
// server.
func testServer() (*http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return mux, server
}
