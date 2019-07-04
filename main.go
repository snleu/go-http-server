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

var m sync.RWMutex

// struct to store requests sent in the last 60 seconds
type DB struct {
	Requests []time.Time `json:"db"`
}

//
func getPersistedRequests(filename string) DB {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print(err)
	}

	var db DB
	err = json.Unmarshal(plan, &db)
	if err != nil {
		// log.error(err)
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

func getRequestCount() int {
	defer m.RUnlock()
	m.RLock()

	db := getPersistedRequests("./db.json")

	return countRequests(db.Requests, time.Now())
}

func persistRequests() error {
	defer m.Unlock()
	m.Lock()

	db := getPersistedRequests("./db.json")

	// append new request to database
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

func incrementRequestCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := persistRequests(); err != nil {
			log.Print(err)
		}
		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "count: %d", getRequestCount())
}

func main() {
	http.Handle("/count", incrementRequestCount(http.HandlerFunc(home)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
