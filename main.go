package main

import (
//  "fmt"
  // "log"
  // "net/http"
  // "encoding/json"
  // "time"

  "./config"
  "./app"
  //
//  "database/sql"
)



func main() {
  config := config.GetConfig()
  app := &app.App{}
  app.Initialize(config)
  app.HandleRequests()
}
