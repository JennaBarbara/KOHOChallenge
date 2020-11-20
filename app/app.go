package app

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "time"
  "regexp"
  "strconv"
  "errors"

  "../config"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}
//common errors
var (
	ErrBadLoadAmount  = errors.New("Invalid Load Amount")
)

//database table model
type LoadedFunds struct {
  Id      uint  `gorm:"primaryKey"`
  Customer_id   uint
  Load_amount    float64
  Time time.Time
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

	a.DB = db.AutoMigrate(&LoadedFunds{})
  a.Router = mux.NewRouter().StrictSlash(true)
}

func (a *App) HandleRequests() {
    a.Router.HandleFunc("/load", loadFunds).Methods("Post")
    log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func loadFunds(w http.ResponseWriter, r *http.Request) {
  defer r.Body.Close()
  dec := json.NewDecoder(r.Body)
  req := &LoadReq{}

  if err := dec.Decode(req); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  amount, err := amountToNumber(req.Load_amount)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  fmt.Printf("x=%v, type of %T\n",amount, amount)

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

//convert given Load_amount from string to float64 format
func amountToNumber(amount string) (float64, error){
	re := regexp.MustCompile(`\$(\d[\d,]*[\.]?[\d{2}]*)`)
  matches := re.FindStringSubmatch(amount)
  if matches == nil {
    return 0.0, ErrBadLoadAmount
  }
  match := matches[1]
  amountFloat, err := strconv.ParseFloat(match, 64);
  if err != nil {
    return 0.0, ErrBadLoadAmount
  }
  fmt.Printf("x=%v, type of %T\n",amountFloat, amountFloat)
  return amountFloat, nil
}
