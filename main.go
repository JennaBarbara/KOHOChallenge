package main

import (
  "fmt"
  "os"
  "bufio"
  "log"
  "encoding/json"

  "./config"
  "./loadFundsApp"
  "./model"
)

func main() {
  //get the program configuration
  config := config.GetConfig()
  a := &loadFundsApp.LoadFundsApp{}
  a.Initialize(config)

  //open the input file
  inputFile, err := os.Open(config.InputFile)
  if err != nil {
        log.Fatalf("failed to open input file")
    }

  //create the output file
  outputFile,err := os.Create(config.OutputFile)
  if err != nil {
      log.Fatalf("failed to create output file")
  }
  outputWriter := bufio.NewWriter(outputFile)

  scanner := bufio.NewScanner(inputFile)
  scanner.Split(bufio.ScanLines)

  //read ionput file line by line and print results
  for scanner.Scan() {
    inputText := scanner.Text()
    req := &model.LoadReq{}
    if err := json.Unmarshal([]byte(inputText), req); err != nil {
      log.Fatalf("failed to read line")
    }
    resp := loadFundsApp.LoadFunds(a, req)
    outputText, err := json.Marshal(resp)
    _, err = fmt.Fprintf(outputWriter, "%s\n", outputText)
    if err != nil {
      log.Fatalf("failed to write to output file")
    }
  }

  //clear memory
  outputWriter.Flush();
  inputFile.Close()
  outputFile.Close()
}
