package model
//
import (
  "fmt"
	"time"
  "regexp"
  "errors"

  "strconv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//common errors
var (
	ErrBadLoadAmount  = errors.New("Invalid load_amount")
  ErrBadId = errors.New("Invalid id")
  ErrBadCustomerId = errors.New("Invalid customer_id")
  ErrExceedsDailyAmountLimit = errors.New("Requested load_amount exceeds daily limit for customer")
  ErrExceedsWeeklyAmountLimit = errors.New("Requested load_amount exceeds weekly limit for customer")
  ErrExceedsDailyLoadLimit = errors.New("Request exceeds daily load limit for customer")
)

//load funds input payload
type LoadReq struct {
    Id      string `json:"id"`
    Customer_id   string `json:"customer_id"`
    Load_amount    string `json:"load_amount"`
    Time time.Time `json:"time"`
}

//load funds output payload
type LoadResp struct {
    Id      string `json:"id"`
    Customer_id   string `json:"customer_id"`
    Accepted    bool `json:"accepted"`
    Error string `json:"error,omitempty"`
}

 //database table model
type LoadedFunds struct {
  Id      int  `gorm:"primaryKey"`
  Customer_id   int
  Load_amount    float64
  Time time.Time
}

//result structure for Database queries that retrieve sums of currency
type Result struct {
  Total float64
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&LoadedFunds{})
	return db
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
func LoadReqToLoadedFunds(req *LoadReq) (*LoadedFunds, error) {
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
  return loadFunds, nil
}

//perform velocity checks described in the requirements
func VelocityLimitsCheck(db *gorm.DB,  loadFund *LoadedFunds) error {
  countToday := CountFundsLoadedToday(db, loadFund)
  if countToday >= 3 {
    err := ErrExceedsDailyLoadLimit
    return err
  }

  sumToday := SumFundsLoadedToday(db, loadFund) + loadFund.Load_amount
  if sumToday > 5000.0 {
    err := ErrExceedsDailyAmountLimit
    return err
  }

  sumThisWeek := SumFundsLoadedThisWeek(db, loadFund) + loadFund.Load_amount
  if sumThisWeek > 20000.0 {
    err := ErrExceedsWeeklyAmountLimit
    return err
  }

  return nil
}

//get sum of all funds loaded on given day
func  SumFundsLoadedThisWeek(db *gorm.DB,  loadedFunds *LoadedFunds) float64 {
  var result Result
  year,week := loadedFunds.Time.ISOWeek();
  yearweek := fmt.Sprintf("%04d%02d", year, week)
  fmt.Printf("yearweek - %s\n", yearweek)
  db.Model(&LoadedFunds{}).Select("sum(Load_amount) as total").Where("Customer_id = ? AND YEARWEEK(Time,3) = ?", loadedFunds.Customer_id, yearweek).Scan(&result)
  return result.Total
}

//get sum of all funds loaded in a given week
func  SumFundsLoadedToday(db *gorm.DB,  loadedFunds *LoadedFunds) float64 {
  var result Result
  db.Model(&LoadedFunds{}).Select("sum(Load_amount) as total").Where("Customer_id = ? AND DATE(Time) = ?", loadedFunds.Customer_id, loadedFunds.Time.Format("2006-01-02")).Scan(&result)
  return result.Total
}

//get count of loads done today
func  CountFundsLoadedToday(db *gorm.DB,  loadedFunds *LoadedFunds) int64 {
  var result int64
  db.Model(&LoadedFunds{}).Where("Customer_id = ? AND DATE(Time) = ?", loadedFunds.Customer_id, loadedFunds.Time.Format("2006-01-02")).Count(&result)
  return result
}

//add a record of a load Funds being performed to the DB
func InsertLoadedFunds(db *gorm.DB, loadedFunds *LoadedFunds) error {
  result := db.Create(&loadedFunds)
  return result.Error
}
