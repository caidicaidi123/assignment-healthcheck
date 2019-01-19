package main

import (
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type Website struct {
    URL string
    Status bool
}

var websites []Website

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/api/healthcheck", addWebsite).Methods("POST")
    r.HandleFunc("/api/healthcheck", deleteWebsite).Methods("DELETE")
    r.HandleFunc("/api/healthcheck", getAllWebsites).Methods("GET")
    
    go updateEveryNSeconds(15*time.Second, updateAllWebsiteStatus)
	http.ListenAndServe(":8080", r)
}

func getResponseTime(url string) time.Duration {
    req, _ := http.NewRequest("GET", url, nil)
    start := time.Now()
    if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
        fmt.Println("Error:", err)
        return 100*time.Second
    }
    fmt.Printf("%v %v\n", time.Since(start), url)
    return time.Since(start)
}

func addWebsite(w http.ResponseWriter, r *http.Request) {
    var site Website
    body, err := ioutil.ReadAll(r.Body)
    if err == nil && body != nil {
        err = json.Unmarshal(body, &site)
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    // check existence of the url
    for _, website := range websites {
        if website.URL == site.URL {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(site)
            return
        }
    }

    updateOneWebsiteStatus(&site)
    websites = append(websites, site)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(site)

    fmt.Printf("%+v\n", websites)
}

func deleteWebsite(w http.ResponseWriter, r *http.Request) {
    var site Website
    body, err := ioutil.ReadAll(r.Body)
    if err == nil && body != nil {
        err = json.Unmarshal(body, &site)
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    deleteIndex := -1
    for i, website := range websites {
        if website.URL == site.URL {
            deleteIndex = i
            break
        }
    }

    // remove the website 
    if deleteIndex >= 0 {
        websites = append(websites[:deleteIndex], websites[deleteIndex+1:]...)
        w.WriteHeader(http.StatusOK)
        return
    }

    // url not found in existing websites
    w.WriteHeader(http.StatusNotFound)
}

func updateOneWebsiteStatus(site *Website) {
    responseTime := getResponseTime((*site).URL)
    if (responseTime.Nanoseconds()/1000000) > 420 {
        (*site).Status = false
    } else {
        (*site).Status = true
    }
}

func updateAllWebsiteStatus(t time.Time) {
    fmt.Printf("\n********\nAuto health check on %v\n", t)
    for i, website := range websites {
        updateOneWebsiteStatus(&website)
        websites[i] = website
    }
    fmt.Printf("New status: %v", websites)
}

func updateEveryNSeconds(d time.Duration, f func(t time.Time)) {
    for t := range time.Tick(d) {
        f(t)
    }
}

func getAllWebsites(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(websites)
}