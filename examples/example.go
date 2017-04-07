package main

import (
	"fmt"
	_ "fmt"
	"github.com/juliend/go-scrapyd-rest-api"
	"time"
	"log"
)

const (
	scrapyServer string = "http://127.0.0.1:32787"
)

func main1() {
	s := goscrapyd.NewScrapyd(scrapyServer)

	//j, res, err := s.ListJobs("scraper")

	var params = make(map[string]string)
	params["crawler_url"] = "http://example.com/"

	j, res, err := s.Schedule("scraper", "scraper", params)

	//res, _ := s.Status()

	fmt.Printf("--> data : %+v\n", j)

	fmt.Printf("--> resp : %+v\n", res)

	fmt.Printf("--> error : %+v\n", err)

}

var (
	pauseDuration = 5 * time.Second
)

func retrieveJobs(s *goscrapyd.Scrapyd, c chan<- []goscrapyd.ScrapydJob) {
	for {
		data, _, err := s.ListJobs("scraper")
		if err != nil {
			log.Fatal(err)
		}
		c <- data.Pending
		time.Sleep(pauseDuration)
	}

}

func schedulJobs(s *goscrapyd.Scrapyd, c chan []goscrapyd.ScrapydJob) {

	var params = make(map[string]string)
	params["crawler_url"] = "http://example.com/"

	jobs := <-c
	if len(jobs) == 0 {
		fmt.Println("jobs = 0. Adding new items")
		for i := 1;  i<=5; i++ {
			s.Schedule("scraper", "scraper", params)
			fmt.Println("new jobs added #", i)
		}
	}

}

func main() {
	s := goscrapyd.NewScrapyd(scrapyServer)

	c := make(chan []goscrapyd.ScrapydJob)
	go retrieveJobs(s, c)
	for {
		schedulJobs(s, c)
	}

}