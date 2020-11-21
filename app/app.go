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
	ErrBadLoadAmount  = errors.New("Invalid load_amount")
  ErrBadId = errors.New("Invalid id")
  ErrBadCustomerId = errors.New("Invalid customer_id")
)

// //database table model
type LoadedFunds struct {
  Id      int  `gorm:"primaryKey"`
  Customer_id   int
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
    //a.Router.HandleFunc("/load", loadFunds).Methods("Post")
    a.Post("/load", a.handleRequest(loadFunds))
    log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func loadFunds(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
  defer r.Body.Close()
  dec := json.NewDecoder(r.Body)
  req := &LoadReq{}

  if err := dec.Decode(req); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  if loadFunds, err := loadReqToLoadedFunds(req *LoadReq); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  countToday := CountFundsLoadedToday(db, req)
  fmt.Printf("count=%v, type of %T\n",countToday, countToday)

  sumToday := SumFundsLoadedToday(db, req)
  fmt.Printf("sum=%v, type of %T\n",sumToday, sumToday)

  sumThisWeek := SumFundsLoadedThisWeek(db, req)
  fmt.Printf("sum=%v, type of %T\n",sumThisWeek, sumThisWeek)


  InsertLoadedFunds(db, loadFunds)
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

//convert the request format LoadReq to the DB format LoadedFunds
func loadReqToLoadedFunds(req *LoadReq) (*LoadedFunds, error) {
  amount, err := amountToNumber(req.Load_amount)
  if err != nil {
    return nil, err;
  }
  id, err := strconv.Atoi(req.Id)
  if err != nil {
    return nil, ErrBadId
  }
  customerId, err := strconv.Atoi(req.Customer_id)
  if err != nil {
    return nil, ErrBadCustomerId
  }
  fmt.Printf("x=%v, type of %T\n",amount, amount)
  loadFunds := &LoadedFunds{
    Id: id,
    Customer_id: customerId,
    Load_amount: amount,
    Time: req.Time,
  }
  return loadFunds
}

type Result struct {
  Total float64
}

//get sum of all funds loaded on given day
func  SumFundsLoadedThisWeek(db *gorm.DB, req *LoadReq) float64 {
  var result Result
  year,week := req.Time.ISOWeek();
  yearweek := fmt.Sprintf("%04d%02d", year, week)
  fmt.Printf("yearweek - %s\n", yearweek)
  db.Model(&LoadedFunds{}).Select("sum(Load_amount) as total").Where("Customer_id = ? AND YEARWEEK(Time,3) = ?", req.Customer_id, yearweek).Scan(&result)
  return result.Total
}

//get sum of all funds loaded in a given week
func  SumFundsLoadedToday(db *gorm.DB, req *LoadReq) float64 {
  var result Result
  db.Model(&LoadedFunds{}).Select("sum(Load_amount) as total").Where("Customer_id = ? AND DATE(Time) = ?", req.Customer_id, req.Time.Format("2006-01-02")).Scan(&result)
  return result.Total
}

//get count of loads done today
func  CountFundsLoadedToday(db *gorm.DB, req *LoadReq) int64 {
  var result int64
  db.Model(&LoadedFunds{}).Where("Customer_id = ? AND DATE(Time) = ?", req.Customer_id, req.Time.Format("2006-01-02")).Count(&result)
  return result
}

//add a record of a load Funds being performed to the DB
func InsertLoadedFunds(db *gorm.DB, loadedFunds *LoadedFunds) {
    db.Create(&loadedFunds);
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
