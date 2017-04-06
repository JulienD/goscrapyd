package goscrapyd

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"bytes"
)

const (
	add_version_endpoint = "add_version"
	cancel_endpoint = "cancel"
	delete_project_endpoint= "delete_project"
	delete_version_endpoint = "delete_version"
	list_jobs_endpoint = "list_jobs"
	list_projects_endpoint = "list_projects"
	list_spiders_endpoint = "list_spiders"
	list_versions_endpoint = "list_versions"
	schedule_endpoint = "schedule"
	daemon_status = "daemonstatus"

	finished_status = "finished"
	pending_status = "pending"
	running_status = "running"
)

var (
	endpoints = map[string]string {
		add_version_endpoint: "addversion.json",
		cancel_endpoint: "cancel.json",
		delete_project_endpoint: "delproject.json",
		delete_version_endpoint: "delversion.json",
		list_jobs_endpoint: "listjobs.json",
		list_projects_endpoint: "listprojects.json",
		list_spiders_endpoint: "listspiders.json",
		list_versions_endpoint: "listversions.json",
		schedule_endpoint: "schedule.json",
		daemon_status: "daemonstatus.json",
	}

	job_states = []string{
		finished_status,
		pending_status,
		running_status,
	}
)


type ScrapydResponse struct {
	Status   string `json:"status"`
	NodeName string `json:"node_name"`
	Jobid    string `json:"jobid"`
	Message  string `json:"message"`
	Pending  []struct {
		ID     string `json:"id"`
		Spider string `json:"spider"`
	} `json:"pending"`
	Running []struct {
		ID        string `json:"id"`
		Spider    string `json:"spider"`
		StartTime string `json:"start_time"`
	} `json:"running"`
	Finished []struct {
		ID        string `json:"id"`
		Spider    string `json:"spider"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	} `json:"finished"`
}

type Scrapyd struct {
	Host    string
}

func NewScrapyd(host string) *Scrapyd {
	var scrapyd = Scrapyd{host}
	return &scrapyd
}
/**
 * Builds the absolute URL using the target and desired endpoint.
 */
func (s *Scrapyd) buildUrl(endpoint string) (string, error) {
	if val, ok := endpoints[endpoint]; ok {
		s := fmt.Sprintf("%v/%v", s.Host, val)

		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		return u.String(), nil

	} else {
		return "", fmt.Errorf("Unknown endpoint %q", val)
	}
}


// Get do a http GET query to the provided url.
func (s *Scrapyd) get(u string, params map[string]string) (*http.Response, error) {

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		// handle err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return resp, nil
}

// Post do a http POST query to the provided url.
func (s *Scrapyd) post(u string, data map[string]string) (*http.Response, error) {

	body := url.Values{}
	if len(data) > 0  {
		for k, v := range data {
			body.Add(k, v)
		}
	}

	req, err := http.NewRequest("POST", u, bytes.NewBufferString(body.Encode())) // <-- URL-encoded payload
	if err != nil {
		fmt.Errorf("%q", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return resp, nil
}

func (s *Scrapyd) handleResponse(r *http.Response, i interface{}) (interface{}, error) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		panic(err)
	}

	return i, nil
}

// ListJobs returns all known jobs given a project.
func (s *Scrapyd) ListJobs(project string) (ScrapydResponse, error) {

	u, err := s.buildUrl(list_jobs_endpoint)
	if err != nil {
		// @TODO
		// handle err
	}

	params := make(map[string]string)
	params["project"] = project

	res, _ := s.get(u, params)

	var sr ScrapydResponse
	s.handleResponse(res, &sr)
	return sr, nil

}

/**
 * Schedules a spider from a specific project to run.
 */
func (s *Scrapyd) Schedule(project string, spider string, settings map[string]string) (*ScrapydResponse, error) {

	params := make(map[string]string)
	params["project"] = project
	params["spider"] = spider

	if len(settings) > 0  {
		for k, v := range settings {
			params[k] = v
		}
	}

	u, err := s.buildUrl(schedule_endpoint)
	if err != nil {
		// @TODO
		// handle err
	}

	res, _ := s.post(u, params)
	var sr ScrapydResponse
	s.handleResponse(res, &sr)

	if sr.Status != "ok" {
		fmt.Printf("Error when requesting the server %+v\n", sr)
	}

	return &sr, nil

}

