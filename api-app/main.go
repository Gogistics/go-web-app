package main

import (
  "log"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/Gogistics/go-web-app/api-app/types"
)

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function
func newRouter() *mux.Router {
  r := mux.NewRouter()
  r.HandleFunc("/api/v1/hello", handler).Methods("GET")
  return r
}

func main() {
  // The router is now formed by calling the `newRouter` constructor function
  // that we defined above. The rest of the code stays the same
  r := newRouter()
  err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", r)
  if err != nil {
    log.Fatal("ListenAndServeTLS: ", err)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
  profile := types.Profile{"Alan", []string{"workout", "programming", "driving"}}
  jProfile, err := json.Marshal(profile)

  if err != nil {
    // handle err
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "applicaiton/json; charset=utf-8")
  w.Write(jProfile) 
}
