//
// Firethorn.
//

package firethorn

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	listenAddr = flag.String("http", ":8080", "HTTP listen address")
	config     = flag.String("config", "firethorn.conf", "firethorn redis configuration")
)

func main() {
	flag.Parse()

	// seed the PNRG used by firethorn to distribute requests
	rand.Seed(time.Now().Unix())

	http.HandleFunc("/get", GetCounter)
	http.HandleFunc("/set", SetCounter)

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func GetCounter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getcounter")
}

func SetCounter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "setcounter")
}
