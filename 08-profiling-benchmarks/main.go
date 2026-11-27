package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"strings"
)

// ----- Benchmarking -----

func BuildStringV1(strs []string) string {
	var res string

	for _, s := range strs {
		res += s
	}

	return res
}

func BuildStringV2(strs []string) string {
	var b strings.Builder

	for _, s := range strs {
		b.WriteString(s)
	}

	return b.String()
}

// ----- Profiling -----

var garbage []*int

func LoaderHandler(w http.ResponseWriter, r *http.Request) {

	n := rand.IntN(100)
	garbage = append(garbage, &n)

	fmt.Println("litter is dropped")
	w.WriteHeader(200)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))

	mux.HandleFunc("/load", LoaderHandler)

	log.Fatal(http.ListenAndServe("127.0.0.1:6060", mux))
}
