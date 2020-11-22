package loadFundsApp

import (
  "fmt"
  "log"
  "errors"

  "../config"
  "../model"
  "github.com/jinzhu/gorm"
)


type LoadFundsApp struct {
	DB     *gorm.DB
  VL  *config.VelocityLimits
}

//common errors
var (
  ErrExceedsDailyAmountLimit = errors.New("Requested load_amount exceeds daily limit for customer")
  ErrExceedsWeeklyAmountLimit = errors.New("Requested load_amount exceeds weekly limit for customer")
  ErrExceedsDailyLoadLimit = errors.New("Request exceeds daily load limit for customer")
)

// Initialize initializes the app with predefined configuration
func (a *LoadFundsApp) Initialize(config *config.Config) {
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
  a.VL = config.VL
}

func LoadFunds(a *LoadFundsApp, req *model.LoadReq) *model.LoadResp  {
  resp := &model.LoadResp{
    Id: req.Id,
    Customer_id: req.Customer_id,
  }

  loadFund, err := model.LoadReqToLoadedFunds(req)
  if err != nil {
    writeBadRequestResponse( resp, err)
    return resp
  }

  if err := velocityLimitsCheck(a, loadFund); err != nil {
    writeBadRequestResponse(resp, err)
    return resp
  }

  if err := model.InsertLoadedFund(a.DB, loadFund); err != nil {
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
func velocityLimitsCheck(a *LoadFundsApp,  loadFund *model.LoadedFunds) error {
  countToday := model.CountFundsLoadedToday(a.DB, loadFund)
  if countToday >= a.VL.DailyLoadLimit {
    err := ErrExceedsDailyLoadLimit
    return err
  }

  sumToday := model.SumFundsLoadedToday(a.DB, loadFund) + loadFund.Load_amount
  if sumToday > a.VL.DailyAmountLimit {
    err := ErrExceedsDailyAmountLimit
    return err
  }

  sumThisWeek := model.SumFundsLoadedThisWeek(a.DB, loadFund) + loadFund.Load_amount
  if sumThisWeek > a.VL.WeeklyAmountLimit {
    err := ErrExceedsWeeklyAmountLimit
    return err
  }

  return nil
}
