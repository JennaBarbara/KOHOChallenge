package main

import (
  "./config"
  "./app"
)

func main() {
  config := config.GetConfig()
  app := &app.App{}
  app.Initialize(config)
  app.HandleRequests()
}
