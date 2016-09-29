package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	Filterfeatures []string `json:"filterfeatures"`
	FilterSuburbs  []string `json:"filterSuburbs"`
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
	// opt, _ := json.Marshal(option)
	// if len(string(opt)) > 0 {
	eptOpt := QueryOption{
		Budget:     0,
		Categories: nil,
		Features:   nil,
		Suburbs:    nil,
		BDPick:     false,
		Newly:      false}
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
	fmt.Println("Testing direct matching logic")
	query = "hinoki japanese pantry"
	result = fetchResult(query, opt)
	if len(result.Filterfeatures) > 0 {
		fmt.Println("Should not parse a direct match")
	} else {
		fmt.Println("Pass")
	}
}

// testSuburbAndFeatureMatch will parse query "fitzroy sushi"
// as suburb "fitzroy" + feature "sushi"
func testSuburbAndFeatureMatch() {
	query = "fitzroy sushi"
	result = fetchResult(query, opt)
	// sort.Strings(result.Filterfeatures)
	// i := sort.SearchStrings(result.Filterfeatures, "sushi")
	if !stringInSlice("sushi", result.Filterfeatures) {
		fmt.Println("Should do features match")
	} else if !stringInSlice("fitzroy", result.FilterSuburbs) {
		fmt.Println("Should do suburbs match")
	} else {
		fmt.Println("Pass")
	}
}

func main() {
	//testSth()
	// testDirectMatch()
	testSuburbAndFeatureMatch()
}

func stringInSlice(a string, list []string) bool {
	sort.Strings(list)
	i := sort.SearchStrings(list, a)
	if i >= 0 {
		return true
	}
	return false
}
