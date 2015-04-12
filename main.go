package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type SuggestRequest struct {
	AddressSuggest AddressSuggest `json:"address-suggest"`
}

type AddressSuggest struct {
	Text       string     `json:"text"`
	Completion Completion `json:"completion"`
}

type Completion struct {
	Field string `json:"field"`
	Size  int    `json:"size":`
}

type SuggestResponse struct {
	Suggestions []Suggestion `json:"address-suggest"`
}
type Suggestion struct {
	Addresses []Address `json:"options"`
}

type Address struct {
	Address string   `json:"text"`
	Loc     Location `json:"payload"`
}

type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Result struct {
	Results []SuggestedAddress `json:"results"`
}

type SuggestedAddress struct {
	Address string  `json:"address"`
	Lon     float64 `json:"lon"`
	Lat     float64 `json:"lat"`
}

type ErrorResult struct {
	StatusCode int
	Message    string
}

func newSuggestRequest(query string, num int) (request *SuggestRequest) {
	request = &SuggestRequest{AddressSuggest{query, Completion{"suggest", num}}}
	return
}

func main() {
	http.HandleFunc("/autocomplete", autocompleteHandler)

	port := ":" + os.Getenv("AUTOCOMPLETE_PORT")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error: %v", err)
	}

	fmt.Println("Server listening on 8000")
}

func autocompleteHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	err, resp := getSuggestions(req)

	if err.Message != "" {
		res.WriteHeader(err.StatusCode)
		errMsg := `{"error": "` + err.Message + `"}`
		res.Write([]byte(errMsg))
		return
	}

	res.WriteHeader(200)
	res.Write(resp)
	return
}

func getSuggestions(req *http.Request) (errRes ErrorResult, resp []byte) {

	var (
		err        error
		reqBytes   []byte
		num        int
		esReq      *http.Request
		esRes      *http.Response
		esResBytes []byte
		sr         *SuggestResponse
		results    []SuggestedAddress
		result     *Result
	)

	params := req.URL.Query()
	if len(params.Get("q")) == 0 {
		errRes.StatusCode = 400
		errRes.Message = "invalid request: 'q' parameter is required"
		return
	}

	num = 10

	if len(params.Get("num")) != 0 {
		if num, err = strconv.Atoi(params.Get("num")); err != nil {
			errRes.StatusCode = 400
			errRes.Message = "invalid request: 'num' parameter is not an integer"
			return
		}
		if num > 25 {
			num = 10
		}
	}

	r := newSuggestRequest(params.Get("q"), num)
	reqBytes, err = json.Marshal(r)
	if err != nil {
		errRes.StatusCode = 500
		errRes.Message = "server error: problem making request to Elasticsearch"
		return
	}

	client := &http.Client{}
	url := "http://localhost:9200/addresses/_suggest"
	esReq, err = http.NewRequest("POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		errRes.StatusCode = 500
		errRes.Message = "server error: problem making request to Elasticsearch"
		return
	}
	esReq.Header.Add("Content-Type", "application/json")

	esRes, err = client.Do(esReq)
	esResBytes, err = ioutil.ReadAll(esRes.Body)
	esRes.Body.Close()
	if err = json.Unmarshal(esResBytes, &sr); err != nil {
		errRes.StatusCode = 500
		errRes.Message = "server error: problem parsing response from Elasticsearch"
		return
	}

	if len(sr.Suggestions) == 0 {
		resp = []byte(nil)
		return
	}

	for _, addr := range sr.Suggestions[0].Addresses {
		r := SuggestedAddress{addr.Address, addr.Loc.Lon, addr.Loc.Lat}
		results = append(results, r)
	}

	result = &Result{results}

	resp, err = json.Marshal(result)
	if err != nil {
		errRes.StatusCode = 500
		errRes.Message = "server error: problem dealing with Elasticsearch results"
		return
	}
	return
}
