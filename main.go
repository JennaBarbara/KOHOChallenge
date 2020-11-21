package main

import (
  "fmt"
  "os"
  "bufio"
  "log"
  "encoding/json"

  "./config"
  "./app"
  "./model"
)

func main() {
  //get the program configuration
  config := config.GetConfig()

  //open the input file
  inputFile, err := os.Open(config.InputFile)
  if err != nil {
        log.Fatalf("failed to open input file")
    }

  a := &app.App{}
  a.Initialize(config)

  //create the output file
  outputFile,err := os.Create(config.OutputFile)
  if err != nil {
      log.Fatalf("failed to create output file")
  }
  outputWriter := bufio.NewWriter(outputFile)

  scanner := bufio.NewScanner(inputFile)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    inputText := scanner.Text()
    req := &model.LoadReq{}
    if err := json.Unmarshal([]byte(inputText), req); err != nil {
      log.Fatalf("failed to read line %d")
    }
    resp := app.LoadFunds(a.DB, req)
    outputText, err := json.Marshal(resp)
    _, err = fmt.Fprintf(outputWriter, "%s\n", outputText)
    if err != nil {
      log.Fatalf("failed to write to output file")
    }
  }
  outputWriter.Flush();
  inputFile.Close()
  outputFile.Close()
}
