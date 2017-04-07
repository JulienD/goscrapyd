package goscrapyd

import (
	"bytes"
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
	"net/url"
)

const (
	add_version_endpoint    = "add_version"
	cancel_endpoint         = "cancel"
	delete_project_endpoint = "delete_project"
	delete_version_endpoint = "delete_version"
	list_jobs_endpoint      = "list_jobs"
	list_projects_endpoint  = "list_projects"
	list_spiders_endpoint   = "list_spiders"
	list_versions_endpoint  = "list_versions"
	schedule_endpoint       = "schedule"
	daemon_status           = "daemonstatus"

	finished_status = "finished"
	pending_status  = "pending"
	running_status  = "running"
)

var (
	endpoints = map[string]string{
		add_version_endpoint:    "addversion.json",
		cancel_endpoint:         "cancel.json",
		delete_project_endpoint: "delproject.json",
		delete_version_endpoint: "delversion.json",
		list_jobs_endpoint:      "listjobs.json",
		list_projects_endpoint:  "listprojects.json",
		list_spiders_endpoint:   "listspiders.json",
		list_versions_endpoint:  "listversions.json",
		schedule_endpoint:       "schedule.json",
		daemon_status:           "daemonstatus.json",
	}

	job_states = []string{
		finished_status,
		pending_status,
		running_status,
	}
)

type Scrapyd struct {
	sling *sling.Sling
	Host  string
}

func NewScrapyd(h string) *Scrapyd {
	return &Scrapyd{
		sling: sling.New().Base(h),
		Host:  h,
	}
}

type ScrapydError struct {
	Status   string `json:"status"`
	NodeName string `json:"node_name"`
	Message  string `json:"message"`
	Errors   []struct {
		Resource string `json:"resource"`
		Field    string `json:"field"`
		Code     string `json:"code"`
	} `json:"errors"`
}

func (e ScrapydError) Error() string {
	return fmt.Sprintf("Error scrapyd. Got the following response %+v", e.Message)
}

/**
 * Returns a list with per statuses of all scheduled jobs.
 */

type ScrapydJobListResponse struct {
	Status   string       `json:"status"`
	NodeName string       `json:"node_name"`
	Jobid    string       `json:"jobid"`
	Message  string       `json:"message"`
	Pending  []ScrapydJob `json:"pending"`
	Running  []ScrapydJob `json:"running"`
	Finished []ScrapydJob `json:"finished"`
}

type ScrapydJob struct {
	ID        string `json:"id"`
	Spider    string `json:"spider"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// ListJobs returns all known jobs given a project.
func (s *Scrapyd) ListJobs(project string) (*ScrapydJobListResponse, *http.Response, error) {

	joblist := new(ScrapydJobListResponse)
	scrapydError := new(ScrapydError)

	type Options struct {
		Project string `url:"project"`
	}
	opt := Options{"scraper"}

	resp, err := s.sling.Get(endpoints[list_jobs_endpoint]).QueryStruct(opt).Receive(joblist, scrapydError)

	if err != nil {
		return nil, resp, scrapydError
	}

	if joblist.Status == "error" {
		scrapydError.Status = joblist.Status
		scrapydError.Message = joblist.Message
		return nil, resp, scrapydError
	}

	return joblist, resp, nil
}

/**
 * Schedules a spider from a specific project to run.
 */

type ScrapydScheduleResponse struct {
	Status   string `json:"status"`
	NodeName string `json:"node_name"`
	Jobid    string `json:"jobid"`
}

func (s *Scrapyd) Schedule(project string, spider string, settings map[string]string) (*ScrapydScheduleResponse, *http.Response, error) {

	schedule := new(ScrapydScheduleResponse)
	scrapydError := new(ScrapydError)

	b := url.Values{}
	b.Add("project", project)
	b.Add("spider", spider)
	if len(settings) > 0 {
		for k, v := range settings {
			b.Add(k, v)
		}
	}

	s.sling.Post(endpoints[schedule_endpoint])
	s.sling.Set("Content-Type", "application/x-www-form-urlencoded")
	s.sling.Body(bytes.NewBufferString(b.Encode()))
	resp, err := s.sling.Receive(schedule, scrapydError)

	if err != nil {
		return nil, resp, scrapydError
	}
	return schedule, resp, err
}
