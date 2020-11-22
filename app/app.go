package app

import (
  "fmt"
  "log"
  "errors"

  "../config"
  "../model"
  "github.com/jinzhu/gorm"
)


type App struct {
	DB     *gorm.DB
}

//common errors
var (
  ErrExceedsDailyAmountLimit = errors.New("Requested load_amount exceeds daily limit for customer")
  ErrExceedsWeeklyAmountLimit = errors.New("Requested load_amount exceeds weekly limit for customer")
  ErrExceedsDailyLoadLimit = errors.New("Request exceeds daily load limit for customer")
)

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
}

func LoadFunds(db *gorm.DB, req *model.LoadReq) *model.LoadResp  {
  resp := &model.LoadResp{
    Id: req.Id,
    Customer_id: req.Customer_id,
  }

  loadFund, err := model.LoadReqToLoadedFunds(req)
  if err != nil {
    writeBadRequestResponse( resp, err)
    return resp
  }

  if err := velocityLimitsCheck(db, loadFund); err != nil {
    writeBadRequestResponse(resp, err)
    return resp
  }

  if err := model.InsertLoadedFund(db, loadFund); err != nil {
    writeBadRequestResponse(resp, err)
    return resp
  }

  resp.Accepted = true
  return resp
}

func writeBadRequestResponse(resp *model.LoadResp, err error) {
  resp.Accepted = false
  log.Printf("Request ID: %s - Error: %s ", resp.Id, err.Error())
}

//perform velocity checks described in the requirements
func velocityLimitsCheck(db *gorm.DB,  loadFund *model.LoadedFunds) error {
  countToday := model.CountFundsLoadedToday(db, loadFund)
  if countToday >= 3 {
    err := ErrExceedsDailyLoadLimit
    return err
  }

  sumToday := model.SumFundsLoadedToday(db, loadFund) + loadFund.Load_amount
  if sumToday > 5000.0 {
    err := ErrExceedsDailyAmountLimit
    return err
  }

  sumThisWeek := model.SumFundsLoadedThisWeek(db, loadFund) + loadFund.Load_amount
  if sumThisWeek > 20000.0 {
    err := ErrExceedsWeeklyAmountLimit
    return err
  }

  return nil
}
