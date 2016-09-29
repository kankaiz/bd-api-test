package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// QueryOption ...
type QueryOption struct {
	Budget     int64    `json:"budget"`
	Categories []string `json:"categories"`
	Features   []string `json:"features"`
	Suburbs    []string `json:"suburbs"`
	BDPick     bool     `json:"BDPick"`
	Newly      bool     `json:"newly"`
}

// QueryResult ...
type QueryResult struct {
	Count int64 `json:"count"`
}

func getJSON(url string) QueryResult {
	r, err := http.Get(url)
	if err != nil {
		// return err
		panic(err.Error())
	}
	defer r.Body.Close()
	var message QueryResult
	if r.StatusCode == 200 { // OK
		// bodyBytes, _ := ioutil.ReadAll(r.Body)
		json.NewDecoder(r.Body).Decode(&message)
		// fmt.Printf("%v", message["hasProfile"])
	}
	return message
}

func fetchResult(query string, option QueryOption) QueryResult {
	url := fmt.Sprintf("http://localhost:8000/melbourne/api/search/%s?format=json", query)
	opt, _ := json.Marshal(option)
	if len(string(opt)) > 0 {
		url = url + "&o=" + string(opt)
	}
	fmt.Printf("fetch: %v", url)
	m := getJSON(url)
	return m
}

func testSth() {
	var (
		query  string
		opt    QueryOption
		result QueryResult
	)
	opt = QueryOption{
		Features: []string{"takeaway"}}
	result = fetchResult(query, opt)
	// count = result["count"] //type assertion convert interface{} into int
	if result.Count < 20 {
		fmt.Printf("count should be greater than 20 while it is %v", result.Count)
	} else {
		fmt.Printf("Pass")
	}
}

func main() {
	testSth()

}
