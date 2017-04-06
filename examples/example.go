package main

import (
	"github.com/juliend/go-scrapyd-rest-api"
_	"fmt"
	"fmt"
)


const (
	scrapyServer string = "http://127.0.0.1:32785"
)

func main () {
	s := goscrapyd.NewScrapyd(scrapyServer)

//	j, _ := s.ListJobs("scraper")

	var params = make(map[string]string)
	params["crawler"] = "http://example.com"

	//s.Schedule("scraper", "scraper", params)

	res, _ := s.Status()

	fmt.Printf("%v\n", res)

}