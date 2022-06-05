package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/rs/cors"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type gitpodSetup struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"cluster_name"`
	Secret  string `json:",omitempty"`
	Keyfile string `json:"sa_credentials"`
	SA      string `json:"sa_name"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(`{"hello": "world"}`)
}

var allGitpods []gitpodSetup

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}

var ai autoInc // global instance

func createGitpod(w http.ResponseWriter, r *http.Request) {
	newGitpod := gitpodSetup{ID: ai.ID()}
	log.Warn(newGitpod.ID)
	log.Warn(newGitpod)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Missing GCP vars for cluster creation")
		return
	}

	err = json.Unmarshal(reqBody, &newGitpod)
	if err != nil {
		fmt.Fprintf(w, "gcp_project, sa_credentials, sa_name and project_name are necessary")
		return
	}

	allGitpods = append(allGitpods, newGitpod)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newGitpod)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/create", createGitpod).Methods("POST")
	log.Info("Starting server at port 8080")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
