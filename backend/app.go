package main

import (
	"github.com/gorilla/handlers"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"github.com/gorilla/mux"
)

var websites []Website

// Initialize the app router
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// App struct that holds the health check app
type App struct {
	Router *mux.Router
}

// Run will start the application
func (a *App) Run(addr string) {
	// run auto check on all websites every 300 secs / 5 mins
	go updateEveryNSeconds(300*time.Second, updateAllWebsiteStatus)

	http.ListenAndServe(addr, 
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), 
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), 
		handlers.AllowedOrigins([]string{"*"}))(a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/healthcheck", a.addWebsite).Methods("POST")
    a.Router.HandleFunc("/api/healthcheck", a.deleteWebsite).Methods("DELETE")
    a.Router.HandleFunc("/api/healthcheck", a.getAllWebsites).Methods("GET")
}

// POST /api/healthcheck
func (a *App) addWebsite(w http.ResponseWriter, r *http.Request) {
    var Site Website
    body, err := ioutil.ReadAll(r.Body)
    if err == nil && body != nil {
        err = json.Unmarshal(body, &Site)
    }
    w.Header().Set("Content-Type", "application/json")
    // check existence of the url

    exists := false
    index := -1
    for i, website := range websites {
        if website.URL == Site.URL {
            exists = true
            index = i
            break
        }
    }

    updateOneWebsiteStatus(&Site)
    if !exists {
        websites = append(websites, Site)
    } else {
        websites[index] = Site
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(Site)

    fmt.Printf("%+v\n", websites)
}

// DELETE /api/healthcheck
func (a *App) deleteWebsite(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Delete called")
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

// GET /api/healthcheck
func (a *App) getAllWebsites(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(websites)
}

// Time the roundtrip request of a website
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

func updateOneWebsiteStatus(site *Website) {
    responseTime := getResponseTime((*site).URL)
    if (responseTime.Nanoseconds()/1000000) > 800 {
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
