package main

import (
//  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "time"

  "github.com/gorilla/mux"
)

//load funds input payload
type LoadReq struct {
    Id      string `json:"id"`
    Customer_id   string `json:"customer_id"`
    Load_amount    string `json:"load_amount"`
    Time time.Time `json:"time"`
}

type LoadResp struct {
    Id      string `json:"id"`
    Customer_id   string `json:"customer_id"`
    Accepted    bool `json:"accepted"`
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/load", loadFunds).Methods("Post")
    log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func loadFunds(w http.ResponseWriter, r *http.Request) {
  defer r.Body.Close()
  dec := json.NewDecoder(r.Body)
  req := &LoadReq{}

  if err := dec.Decode(req); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  resp := &LoadResp{
    Id: req.Id,
    Customer_id: req.Customer_id,
    Accepted: true,
  }

  w.Header().Set("Content-Type", "application/json")
  enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Printf("can't encode %v - %s", resp, err)
	}
}



func main() {

   handleRequests()
}
