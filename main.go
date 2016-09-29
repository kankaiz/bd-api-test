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

func getJSON(url string) map[string]interface{} {
	r, err := http.Get(url)
	if err != nil {
		// return err
		panic(err.Error())
	}
	defer r.Body.Close()
	// fmt.Printf(r.Status)
	message := make(map[string]interface{})
	if r.StatusCode == 200 { // OK
		// bodyBytes, _ := ioutil.ReadAll(r.Body)
		json.NewDecoder(r.Body).Decode(&message)
		// fmt.Printf("%v", message["hasProfile"])
	}
	return message
}

func fetchResult(query string, option QueryOption) {
	url := fmt.Sprintf("http://localhost:8000/melbourne/api/search/%s?format=json", query)
	opt, _ := json.Marshal(option)
	if len(string(opt)) > 0 {
		url = url + "&o=" + string(opt)
	}
	fmt.Printf("fetch: %v", url)
	m := getJSON(url)
	fmt.Printf("%v", m["hasProfile"])
}

func main() {
	o := QueryOption{
		Features: []string{"takeaway"}}
	fetchResult("", o)

}
