package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

/*
It is possible that multiple requests are sent at the same time. If this happens,
we want to be sure that read and write access to the request data is synchronized.
This has been implemented using a mutex.
*/
var m sync.RWMutex

// struct to store requests sent in the last 60 seconds
type DB struct {
	Requests []time.Time `json:"db"`
}

// helper to read the persisted requests from the json file
func getPersistedRequests(filename string) DB {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print(err)
	}

	var db DB
	err = json.Unmarshal(plan, &db)
	if err != nil {
		log.Print(err)
	}

	return db
}

// count the requests sent within the last 60s
func countRequests(times []time.Time, now time.Time) int {
	x := len(times)
	for x > 0 {
		x -= 1
		if now.Sub(times[x]).Seconds() > 60 {
			times = times[x+1:]
			break
		}
	}
	return len(times)
}

// loads requests from json file and returns number of requests sent within last 60s
func getRequestCount() int {
	defer m.RUnlock()
	m.RLock()

	db := getPersistedRequests("./db.json")

	return countRequests(db.Requests, time.Now())
}

// saves the requests to json file
func persistRequests() error {
	defer m.Unlock()
	m.Lock()

	db := getPersistedRequests("./db.json")
	db.Requests = append(db.Requests, time.Now())

	encodeDb, err := json.Marshal(db)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./db.json", encodeDb, 0644)
	if err != nil {
		return err
	}

	return nil
}

// middleware: when request is received, persist requests
func incrementRequestCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := persistRequests(); err != nil {
			log.Print(err)
		}
		next.ServeHTTP(w, r)
	})
}

// request returns 200 and request count from the last 60s
func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "count: %d", getRequestCount())
}

func main() {
	http.Handle("/count", incrementRequestCount(http.HandlerFunc(home)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
