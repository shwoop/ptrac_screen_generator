package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
var num_imos = len(imos)

const workers = 10
const url = "https://purpletrac.polestar-testing.com/api/v1/registration?api_key=0d5ad56b6de8cfb14f81232aaae3d2543f9313ed&username=alistair.ferguson%40polestarglobal.com"

func pick_a_ship() string {
	i := rand.Intn(num_imos)
	return imos[i]
}

func main() {
	sigs := make(chan int, 1)

	for i := 0; i < workers; i++ {
		go func() {
			for {
				time.Sleep(time.Duration(rand.Intn(60)) * time.Second)
				ship := pick_a_ship()
				body, _ := json.Marshal(map[string]string{
					"registered_name":  ship,
					"client_reference": "altest autogen",
				})

				resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
				if err != nil {
					log.Fatalln(err)
					sigs <- 1
				}
				resp.Body.Close()
				if resp.StatusCode != 201 {
					log.Fatalln(resp.Status)
					sigs <- 1
				}

				fmt.Println("Registered ship " + ship)
			}
		}()
	}

	<-sigs
}
