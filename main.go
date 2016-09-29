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
	Count          int64    `json:"count"`
	Filterfeatures []string `json:"filterfeatures"`
}

var (
	baseURL = "http://localhost:8000"
	query   string
	opt     QueryOption
	result  QueryResult
)

func getJSON(url string) QueryResult {
	r, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer r.Body.Close()
	var message QueryResult
	if r.StatusCode == 200 {
		json.NewDecoder(r.Body).Decode(&message)
	}
	return message
}

func fetchResult(query string, option QueryOption) QueryResult {
	url := fmt.Sprintf("%s/melbourne/api/search/%s?format=json", baseURL, query)
	opt, _ := json.Marshal(option)
	if len(string(opt)) > 0 {
		url = url + "&o=" + string(opt)
	}
	fmt.Printf("fetch: %v", url)
	m := getJSON(url)
	return m
}

func testSth() {
	opt = QueryOption{
		Features: []string{"takeaway"}}
	result = fetchResult(query, opt)
	if result.Count < 20 {
		fmt.Printf("count should be greater than 20 while it is %v", result.Count)
	} else {
		fmt.Printf("Pass")
	}
}

// testDirectMatch matches venue name "hinoki japanese pantry"
// to prevent search enginee parsing name as "hinoki pantry" + feature "Japanese"
func testDirectMatch() {
	fmt.Println("Testing direct matching logic")
	query = "hinoki japanese pantry"
	result = fetchResult(query, opt)
	if len(result.Filterfeatures) > 0 {
		fmt.Println("Should not parse a direct match")
	} else {
		fmt.Printf("Pass")
	}
}

func main() {
	//testSth()
	testDirectMatch()

}
