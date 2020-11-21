package app

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"

  "../config"
  "./model"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
    log.Fatalf("Could not connect database - %s", err)
	}
	a.DB = model.DBMigrate(db)
  a.Router = mux.NewRouter().StrictSlash(true)
}

func (a *App) HandleRequests() {
    //a.Router.HandleFunc("/load", loadFunds).Methods("Post")
    a.Post("/load", a.handleRequest(loadFunds))
    log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func loadFunds(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
  defer r.Body.Close()
  dec := json.NewDecoder(r.Body)
  req := &model.LoadReq{}

  if err := dec.Decode(req); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  resp := &model.LoadResp{
    Id: req.Id,
    Customer_id: req.Customer_id,
  }

  w.Header().Set("Content-Type", "application/json")
  loadFund, err := model.LoadReqToLoadedFunds(req)
  if err != nil {
    writeBadRequestResponse(w, resp, err)
    return
  }

  if err := model.VelocityLimitsCheck(db, loadFund); err != nil {
     writeBadRequestResponse(w, resp, err)
    return
  }

  if err := model.InsertLoadedFunds(db, loadFund); err != nil {
    writeBadRequestResponse(w, resp, err)
    return
  }
  enc := json.NewEncoder(w)
  resp.Accepted = true
	if err := enc.Encode(resp); err != nil {
		log.Printf("can't encode %v - %s", resp, err)
	}
}

func writeBadRequestResponse(w http.ResponseWriter, resp *model.LoadResp, err error) {
  resp.Accepted = false
  resp.Error = err.Error()
  w.WriteHeader(http.StatusBadRequest)
  enc := json.NewEncoder(w)
  if err := enc.Encode(resp); err != nil {
    log.Printf("can't encode %v - %s", resp, err)
  }
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
