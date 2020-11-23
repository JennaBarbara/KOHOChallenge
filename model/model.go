package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

//common errors
var (
	ErrBadLoadAmount = errors.New("Invalid load_amount")
	ErrBadId         = errors.New("Invalid id")
	ErrBadCustomerId = errors.New("Invalid customer_id")
)

//load funds input payload
type LoadReq struct {
	Id          string    `json:"id"`
	Customer_id string    `json:"customer_id"`
	Load_amount string    `json:"load_amount"`
	Time        time.Time `json:"time"`
}

//load funds output payload
type LoadResp struct {
	Id          string `json:"id"`
	Customer_id string `json:"customer_id"`
	Accepted    bool   `json:"accepted"`
}

//database table model
type LoadedFunds struct {
	Id          int
	Customer_id int
	Load_amount float64
	Time        time.Time
	Accepted    bool
}

//result structure for Database queries that retrieve sums of currency
type Result struct {
	Total float64
}

// link db object to the LoadedFunds struct
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&LoadedFunds{})
	return db
}

//convert given Load_amount from string to float64 format
func amountToNumber(amount string) (float64, error) {
	re := regexp.MustCompile(`\$(\d[\d,]*[\.]?[\d{2}]*)`)
	matches := re.FindStringSubmatch(amount)
	if matches == nil {
		return 0.0, ErrBadLoadAmount
	}
	match := matches[1]
	amountFloat, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0.0, ErrBadLoadAmount
	}
	return amountFloat, nil
}

//convert the request format LoadReq to the DB format LoadedFunds
func LoadReqToLoadedFunds(req *LoadReq) (*LoadedFunds, error) {
	amount, err := amountToNumber(req.Load_amount)
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, ErrBadId
	}
	customerId, err := strconv.Atoi(req.Customer_id)
	if err != nil {
		return nil, ErrBadCustomerId
	}
	loadFunds := &LoadedFunds{
		Id:          id,
		Customer_id: customerId,
		Load_amount: amount,
		Time:        req.Time,
		Accepted:    false,
	}
	return loadFunds, nil
}

//get sum of all funds loaded on given day
func SumFundsLoadedThisWeek(db *gorm.DB, loadedFunds *LoadedFunds) float64 {
	var result Result
	year, week := loadedFunds.Time.ISOWeek()
	yearweek := fmt.Sprintf("%04d%02d", year, week)
	db.Model(&LoadedFunds{}).Select("sum(Load_amount) as total").Where("Customer_id = ? AND YEARWEEK(Time,3) = ? AND Accepted = true",
		loadedFunds.Customer_id, yearweek).Scan(&result)
	return result.Total
}

//get sum of all funds loaded in a given week
func SumFundsLoadedToday(db *gorm.DB, loadedFunds *LoadedFunds) float64 {
	var result Result
	db.Model(&LoadedFunds{}).Select("sum(Load_amount) as total").Where("Customer_id = ? AND DATE(Time) = ? AND Accepted = true",
		loadedFunds.Customer_id,
		loadedFunds.Time.Format("2006-01-02")).Scan(&result)
	return result.Total
}

//get count of loads done today
func CountFundsLoadedToday(db *gorm.DB, loadedFunds *LoadedFunds) int64 {
	var result int64
	db.Model(&LoadedFunds{}).Where("Customer_id = ? AND DATE(Time) = ? AND Accepted = true",
		loadedFunds.Customer_id, loadedFunds.Time.Format("2006-01-02")).Count(&result)
	return result
}

//check if a record with a given Id and Customer_id has already been recorded
func GetExistingRecord(db *gorm.DB, loadedFunds *LoadedFunds) bool {
	result := &LoadedFunds{}
	if db.Where("Id = ? AND Customer_id = ? ",
		loadedFunds.Id, loadedFunds.Customer_id).First(&result).RecordNotFound() {
		return false
	}
	return true
}

//add a record of a load Funds being performed to the DB
func InsertLoadedFund(db *gorm.DB, loadedFunds *LoadedFunds) error {
	result := db.Create(&loadedFunds)
	return result.Error
}

//Delete an existing load funds record
func DeleteLoadedFund(db *gorm.DB, loadedFunds *LoadedFunds) error {
	result := db.Where("Id = ? AND Customer_id = ? AND Load_amount = ? AND Time = ?", loadedFunds.Id, loadedFunds.Customer_id, loadedFunds.Load_amount, loadedFunds.Time).Delete(&LoadedFunds{})
	return result.Error
}
