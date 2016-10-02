package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
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
	HasProfile     bool     `json:"hasProfile"`
	FilterFeatures []string `json:"filterfeatures"`
	FilterSuburbs  []string `json:"filterSuburbs"`
}

var (
	baseURL = "https://www.broadsheet.com.au" //"http://localhost:8000"
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
	eptOpt := QueryOption{
		Budget:     0,
		Categories: nil,
		Features:   nil,
		Suburbs:    nil,
		BDPick:     false,
		Newly:      false}
	// if option is not null then append the request URL
	if !reflect.DeepEqual(eptOpt, option) {
		opt, _ := json.Marshal(option)
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
	fmt.Println("Testing direct matching logic...")
	query = "hinoki japanese pantry"
	result = fetchResult(query, opt)
	if len(result.FilterFeatures) > 0 {
		fmt.Println("Should not parse a direct match")
	} else {
		fmt.Println("Pass")
	}
}

// testSuburbAndFeatureMatch will parse query "fitzroy sushi"
// as suburb "fitzroy" + feature "sushi"
func testSuburbAndFeatureMatch() {
	fmt.Println("Testing suburb and feature matching logic...")
	query = "fitzroy sushi"
	result = fetchResult(query, opt)
	if !stringInSlice("sushi", result.FilterFeatures) {
		fmt.Println("Should do features match")
	} else if !stringInSlice("fitzroy", result.FilterSuburbs) {
		fmt.Println("Should do suburbs match")
	} else {
		fmt.Println("Pass")
	}
}

// testAlias will test whether an alias is correctly parsed
// attention! This alias may be edited, so the test case might fail
func testAlias() {
	fmt.Println("Testing alias matching logic...")
	query = "read"
	result = fetchResult(query, opt)
	if !stringInSlice("books/records", result.FilterFeatures) {
		fmt.Println("Should do alias match")
	} else {
		fmt.Println("Pass")
	}
}

// testBoostProfile will test whether profiles are prioritised
func testBoostProfile() {
	fmt.Println("Testing profile boosting logic...")
	result = fetchResult(query, opt)
	if !result.HasProfile {
		fmt.Println("Profiles should be prioritised")
	} else {
		fmt.Println("Pass")
	}
}

func main() {
	url := os.Getenv("URL")
	if url != "" {
		baseURL = url
	}
	// testSth()
	testBoostProfile()
	testDirectMatch()
	testSuburbAndFeatureMatch()
	testAlias()
}

// general func to see if an element is in slice
func stringInSlice(a string, list []string) bool {
	sort.Strings(list)
	i := sort.SearchStrings(list, a)
	if i >= 0 {
		return true
	}
	return false
}
