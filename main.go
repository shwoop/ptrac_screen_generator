package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const urlPattern = "https://purpletrac.polestar-testing.com/api/v1/" +
	"registration?api_key=%s&username=%s"

func gatherIMOS() ([]string, error) {
	content, err := ioutil.ReadFile("imos")
	if err != nil {
		return nil, err
	}
	imos := strings.Split(string(content), "\n")
	return imos, nil
}

func pickAShip(imos []string, numImos int) string {
	i := rand.Intn(numImos)
	return imos[i]
}

func registerAShip(sigs chan int, url, ship string) {
	body, _ := json.Marshal(map[string]string{
		"registered_name":  ship,
		"client_reference": "altest autogen",
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(fmt.Sprintf("Error posting request to %s: %s", url, err))
		sigs <- 1
	}
	resp.Body.Close()
	if resp.StatusCode != 201 {
		fmt.Println("Invalid response: " + resp.Status)
		sigs <- 1
	}

	fmt.Println("Registered ship " + ship)
}

type Config struct {
	workers  int
	username string
	apiKey   string
	minTime  int
	maxTime  int
}

func parseArguments() Config {
	conf := Config{}

	flag.IntVar(&conf.workers, "w", 1, "Number of parallel workers")
	flag.StringVar(
		&conf.username,
		"u",
		"alistair.ferguson@polestarglobal.com",
		"API user",
	)
	flag.StringVar(
		&conf.apiKey,
		"k",
		"0d5ad56b6de8cfb14f81232aaae3d2543f9313ed",
		"API key for user",
	)
	flag.IntVar(
		&conf.minTime,
		"min",
		60,
		"Minumum time for workers to wait between jobs (seconds)",
	)
	flag.IntVar(
		&conf.maxTime,
		"max",
		180,
		"Minumum time for workers to wait between jobs (seconds)",
	)

	flag.Parse()

	return conf
}

func main() {
	conf := parseArguments()

	imos, err := gatherIMOS()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	numImos := len(imos)
	if numImos == 0 {
		fmt.Println("Found no IMOS")
		os.Exit(1)
	}

	url := fmt.Sprintf(
		urlPattern,
		url.QueryEscape(conf.apiKey),
		url.QueryEscape(conf.username),
	)
	sigs := make(chan int, 1)

	for i := 0; i < conf.workers; i++ {
		go func() {
			for {
				time.Sleep(
					time.Duration(conf.minTime+rand.Intn(conf.maxTime)) *
						time.Second,
				)
				registerAShip(sigs, url, pickAShip(imos, numImos))
			}
		}()
	}

	os.Exit(<-sigs)
}
