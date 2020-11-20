package app

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "time"

  "../config"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

type LoadedFunds struct {
  Id      uint  `gorm:"primaryKey"`
  Customer_id   uint
  Load_amount    float64
  Time time.Time
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

	db.AutoMigrate(&LoadedFunds{})
  fmt.Printf("here I am! %s\n", db)
}

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

func (a *App) HandleRequests() {
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
