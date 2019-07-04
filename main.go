package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DB struct {
	Requests []time.Time `json:"db"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func cut(times []time.Time, now time.Time) ([]time.Time, int) {
	x := len(times)
	for x > 0 {
		x -= 1
		if now.Sub(times[x]).Seconds() > 60 {
			times = times[x+1:]
			break
		}
	}
	return times, len(times)
}

func main() {
	// load previous timeseries request data on server start up
	plan, err := ioutil.ReadFile("./db.json")
	check(err)
	var db DB
	err = json.Unmarshal(plan, &db)
	check(err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		db.Requests = append(db.Requests, now)
		var count int
		db.Requests, count = cut(db.Requests, now)
		fmt.Println(db)
		enCutDb, err := json.Marshal(db)

		fmt.Println(string(enCutDb))
		check(err)
		err = ioutil.WriteFile("./db.json", enCutDb, 0644)
		check(err)
		w.WriteHeader(200)
		fmt.Fprintf(w, "count: %d, db.requests: %s", count, enCutDb)
	})

	http.ListenAndServe(":8080", nil)

}
