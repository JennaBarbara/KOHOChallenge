package app

import (
  "fmt"
  "log"

  "../config"
  "../model"
  "github.com/jinzhu/gorm"
)

type App struct {
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

  if err := model.VelocityLimitsCheck(db, loadFund); err != nil {
     writeBadRequestResponse(resp, err)
    return resp
  }

  if err := model.InsertLoadedFunds(db, loadFund); err != nil {
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
