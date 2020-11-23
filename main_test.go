package main

import (
  "testing"
  "reflect"
  "time"

  "./loadFundsApp"
  "./model"
  "./config"
)


func TestAll(t *testing.T){
  //get the program configuration
  config := config.GetConfig()
  //initialize the app
  a := &loadFundsApp.LoadFundsApp{}
  a.Initialize(config)

  //perform Tests
  testSuccessfulLoad(t, a)
  testDuplicateRequest(t, a)
  testDailyAmountLimitValidation(t, a)
   testWeeklyAmountLimitValidation(t, a)


}

func testSuccessfulLoad(t *testing.T, a *loadFundsApp.LoadFundsApp){
  reqTime, _ := time.Parse(time.RFC3339, "2020-01-03T00:00:00Z")
  req := &model.LoadReq{ Id: "1", Customer_id: "1", Load_amount: "$100.00", Time: reqTime}
  expectedResp := &model.LoadResp{ Id: "1", Customer_id: "1", Accepted: true }
  actualResp := loadFundsApp.LoadFunds(a, req)

  loadedFunds, _ := model.LoadReqToLoadedFunds(req)
  model.DeleteLoadedFund(a.DB, loadedFunds)
  if !reflect.DeepEqual(expectedResp, actualResp) {
    	t.Fatalf("testSuccessfulLoad failed!")
  }
}

func testDuplicateRequest(t *testing.T, a *loadFundsApp.LoadFundsApp){
  reqTime, _ := time.Parse(time.RFC3339, "2020-01-03T00:00:00Z")
  req := &model.LoadReq{ Id: "1", Customer_id: "1", Load_amount: "$100.00", Time: reqTime}
  loadFundsApp.LoadFunds(a, req)
  actualResp := loadFundsApp.LoadFunds(a, req)

  loadedFunds, _ := model.LoadReqToLoadedFunds(req)
  model.DeleteLoadedFund(a.DB, loadedFunds)
  if  actualResp != nil {
    	t.Fatalf("testDuplicateRequest failed!")
  }
}

func testDailyAmountLimitValidation(t *testing.T, a *loadFundsApp.LoadFundsApp){
  reqTime, _ := time.Parse(time.RFC3339, "2020-01-03T00:00:00Z")
  req := &model.LoadReq{ Id: "1", Customer_id: "1", Load_amount: "$5100.00", Time: reqTime}
  expectedResp := &model.LoadResp{ Id: "1", Customer_id: "1", Accepted: false }
  actualResp := loadFundsApp.LoadFunds(a, req)

  if !reflect.DeepEqual(expectedResp, actualResp) {
    	t.Fatalf("testDailyAmountLimitValidation failed!")
  }
  loadedFunds,_ := model.LoadReqToLoadedFunds(req)
  model.DeleteLoadedFund(a.DB, loadedFunds);
}

func testWeeklyAmountLimitValidation(t *testing.T, a *loadFundsApp.LoadFundsApp){

  req1Time, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
  req2Time, _ := time.Parse(time.RFC3339, "2020-01-02T00:00:00Z")
  req3Time, _ := time.Parse(time.RFC3339, "2020-01-03T00:00:00Z")
  req4Time, _ := time.Parse(time.RFC3339, "2020-01-04T00:00:00Z")
  req5Time, _ := time.Parse(time.RFC3339, "2020-01-05T00:00:00Z")

  req1 := &model.LoadReq{ Id: "1", Customer_id: "1", Load_amount: "$4000.00", Time: req1Time}
  req2 := &model.LoadReq{ Id: "2", Customer_id: "1", Load_amount: "$4000.00", Time: req2Time}
  req3 := &model.LoadReq{ Id: "3", Customer_id: "1", Load_amount: "$4000.00", Time: req3Time}
  req4 := &model.LoadReq{ Id: "4", Customer_id: "1", Load_amount: "$4000.00", Time: req4Time}
  req5 := &model.LoadReq{ Id: "5", Customer_id: "1", Load_amount: "$4100.00", Time: req5Time}

  expectedResp := &model.LoadResp{ Id: "5", Customer_id: "1", Accepted: false }
  loadFundsApp.LoadFunds(a, req1)
  loadFundsApp.LoadFunds(a, req2)
  loadFundsApp.LoadFunds(a, req3)
  loadFundsApp.LoadFunds(a, req4)
  actualResp := loadFundsApp.LoadFunds(a, req5)

//delete possible modifications to table
  loadedFunds1,_ := model.LoadReqToLoadedFunds(req1)
  model.DeleteLoadedFund(a.DB, loadedFunds1);
  loadedFunds2,_ := model.LoadReqToLoadedFunds(req2)
  model.DeleteLoadedFund(a.DB, loadedFunds2);
  loadedFunds3,_ := model.LoadReqToLoadedFunds(req3)
  model.DeleteLoadedFund(a.DB, loadedFunds3);
  loadedFunds4,_ := model.LoadReqToLoadedFunds(req4)
  model.DeleteLoadedFund(a.DB, loadedFunds4);
  loadedFunds5,_ := model.LoadReqToLoadedFunds(req5)
  model.DeleteLoadedFund(a.DB, loadedFunds5);

  if !reflect.DeepEqual(expectedResp, actualResp) {
      t.Fatalf("testWeeklyAmountLimitValidation failed!")
  }

}

func testDailyLoadLimitValidation(t *testing.T, a *loadFundsApp.LoadFundsApp){
  reqTime, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")

  req1 := &model.LoadReq{ Id: "1", Customer_id: "1", Load_amount: "$100.00", Time: reqTime}
  req2 := &model.LoadReq{ Id: "2", Customer_id: "1", Load_amount: "$100.00", Time: reqTime}
  req3 := &model.LoadReq{ Id: "3", Customer_id: "1", Load_amount: "$100.00", Time: reqTime}
  req4 := &model.LoadReq{ Id: "4", Customer_id: "1", Load_amount: "$100.00", Time: reqTime}

  expectedResp := &model.LoadResp{ Id: "4", Customer_id: "1", Accepted: false }
  loadFundsApp.LoadFunds(a, req1)
  loadFundsApp.LoadFunds(a, req2)
  loadFundsApp.LoadFunds(a, req3)
  actualResp := loadFundsApp.LoadFunds(a, req4)

//delete possible modifications to table
  loadedFunds1,_ := model.LoadReqToLoadedFunds(req1)
  model.DeleteLoadedFund(a.DB, loadedFunds1);
  loadedFunds2,_ := model.LoadReqToLoadedFunds(req2)
  model.DeleteLoadedFund(a.DB, loadedFunds2);
  loadedFunds3,_ := model.LoadReqToLoadedFunds(req3)
  model.DeleteLoadedFund(a.DB, loadedFunds3);
  loadedFunds4,_ := model.LoadReqToLoadedFunds(req4)
  model.DeleteLoadedFund(a.DB, loadedFunds4);

  if !reflect.DeepEqual(expectedResp, actualResp) {
      t.Fatalf("testDailyLoadLimitValidation failed!")
  }
}
