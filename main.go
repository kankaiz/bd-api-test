package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// QueryOption ...
type QueryOption struct {
	BDPick bool
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
	m := getJSON(url)
	fmt.Printf("%v", m["hasProfile"])
}

func main() {
	var o QueryOption
	fetchResult("", o)

}
