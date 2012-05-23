//
// Firethorn.
//

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	configFileValue = ""
	configFile      = &configFileValue //flag.String("config", "", "firethorn configuration")
)

func main() {
	flag.Parse()

	if len(*configFile) == 0 {
		log.Fatal("No configuration file specified")
	}

	configContents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading %s: %s", *configFile, err)
	}

	var config Config
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		log.Fatalf("Could not parse configuration %s: %s", *configFile, err)
	}

	// seed the PNRG used by firethorn to distribute requests across pools
	rand.Seed(time.Now().Unix())

	http.HandleFunc("/", GetCounter)
	http.HandleFunc("/set", SetCounter)

	log.Printf("Starting Firethorn on %s", config.ListenAddress)

	log.Fatal(http.ListenAndServe(config.ListenAddress, nil))
}

func GetCounter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getcounter")
}

func SetCounter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "setcounter")
}
