package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

var imos = []string{
	"9590113",
	"9539236",
	"9435064",
	"9621417",
	"9431680",
	"8702630",
	"9615145",
	"9402483",
	"9601235",
	"9565170",
	"9440150",
	"9248409",
	"9351725",
	"9323687",
	"9261011",
	"9353802",
}
var numImos = len(imos)

const workers = 1
const urlPattern = "https://purpletrac.polestar-testing.com/api/v1/registration?api_key=%s&username=%s"

const api_key = "0d5ad56b6de8cfb14f81232aaae3d2543f9313ed"

const user = "alistair.ferguson@polestarglobal.com"

func pickAShip() string {
	i := rand.Intn(numImos)
	return imos[i]
}

func registerAShip(sigs chan int, url, ship string) {
	time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
	body, _ := json.Marshal(map[string]string{
		"registered_name":  ship,
		"client_reference": "altest autogen",
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		sigs <- 1
	}
	resp.Body.Close()
	if resp.StatusCode != 201 {
		fmt.Println(resp.Status)
		sigs <- 1
	}

	fmt.Println("Registered ship " + ship)
	sigs <- 0
}

func main() {
	url := fmt.Sprintf(
		urlPattern,
		url.QueryEscape(api_key),
		url.QueryEscape(user),
	)
	fmt.Println(url)
	sigs := make(chan int, 1)

	for i := 0; i < workers; i++ {
		go func() {
			for {
				registerAShip(sigs, url, pickAShip())
			}
		}()
	}

	os.Exit(<-sigs)
}
