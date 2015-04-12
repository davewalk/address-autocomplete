package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Location struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type Suggest struct {
	Input    string   `json:"input"`
	Output   string   `json:"output"`
	Location Location `json:"payload"`
}

type Address struct {
	Name    string  `json:"name"`
	Suggest Suggest `json:"suggest"`
}

func newAddress(addrStr string, latStr string, lonStr string) (a Address, err error) {
	var lon float64
	lon, err = strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return
	}

	var lat float64
	lat, err = strconv.ParseFloat(latStr, 64)
	if err != nil {
		return
	}

	a = Address{addrStr, Suggest{addrStr, addrStr, Location{lon, lat}}}
	return
}

func main() {
	fmt.Println("Indexing the addresses...")

	f, err := os.Open(os.Getenv("AUTOCOMPLETE_FILE"))
	if err != nil {
		panic(err)
	}

	defer f.Close()
	reader := csv.NewReader(f)
	addresses, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d addresses to index...\n", len(addresses))

	saved := 0

	for idx, addr := range addresses {
		var (
			a         Address
			addrBytes []byte
			req       *http.Request
			resp      *http.Response
		)
		a, err = newAddress(addr[0], addr[1], addr[2])
		if err != nil {
			fmt.Printf("Problem with %s: %s", addr[0], err)
		}
		addrBytes, err = json.Marshal(a)

		if err != nil {
			fmt.Printf("Problem marshaling %s: %s", addr[0], err)
		}
		client := &http.Client{}
		idxStr := strconv.Itoa(idx)
		url := "http://localhost:9200/addresses/address/" + idxStr
		req, err = http.NewRequest("PUT", url, bytes.NewReader(addrBytes))
		if err != nil {
			fmt.Println("Problem creating request for %s: %s", addr[0], err)
		} else {
			saved = saved + 1
		}

		if saved%2500 == 0 {
			fmt.Printf("%d addresses saved...\n", saved)
		}

		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Problem PUTing %s: %s", addr[0], err)
		}
		resp.Body.Close()
	}
}
