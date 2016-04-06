package main

import (
  "fmt"
  "flag"
  "os"
  "strings"
  "encoding/csv"
  "io"
  "io/ioutil"
  "github.com/BurntSushi/toml"
  "crypto/sha1"
)

var (
  configPath = flag.String("config", "./default.conf", "Config file path\n\t\tDefault: ./default.conf")
  inputPath = flag.String("input", "./input.csv", "Input file path\n\t\tDefault: ./input.csv")
)

func usage() {
  fmt.Fprintf(os.Stderr, "\nUsage: %s [flags]\n\n", "go-cleaner2")
  flag.PrintDefaults()
  os.Exit(0)
}

func main() {
  flag.Usage = usage
  flag.Parse()

  // read config file
  configStr, err := ioutil.ReadFile(*configPath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "\nConfig file %s read failed: %s", *configPath, err)
    os.Exit(0)
  }

  type Column struct {
    ColumnNo int
    ColumnType string
  }

  type Config struct {
    Columns map[string]Column
  }

  var config Config
  if _, err := toml.Decode(string(configStr), &config); err != nil {
    fmt.Fprintf(os.Stderr, "\nTOML config string decode error: %s\n\n", err)
    os.Exit(0)
  }

  // parse config
  columnConfig := make(map[int]string)
  for col := range config.Columns {
    colInfo := config.Columns[col]
    columnConfig[colInfo.ColumnNo] = colInfo.ColumnType
  }
  fmt.Println(columnConfig)

  // open and read csv file
  inputFile, err := os.Open(*inputPath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "\nInput file %s open failed: %s\n\n", *inputPath, err)
    os.Exit(0)
  }
  defer inputFile.Close()

  // open .processed file to write
  outputFile, err := os.Create(*inputPath + ".processed")
  if err != nil {
    fmt.Fprintf(os.Stderr, "\nOutput file %s open failed: %s\n\n", *inputPath + ".processed", err)
    os.Exit(0)
  }
  defer outputFile.Close()
  csvWriter := csv.NewWriter(outputFile)

  csvReader := csv.NewReader(inputFile)
  for {
    record, err := csvReader.Read()
    if err == io.EOF {
      // finished
      csvWriter.Flush()
      fmt.Printf("\nInput file %s process done. Output file: %s\n\n", *inputPath, *inputPath + ".processed")
      break;
    } else if err != nil {
      fmt.Fprintf(os.Stderr, "\nInput file %s reading err: %s\n\n", *inputPath, err)
      os.Exit(0)
    }

    // process record
    newRecord := record[:]
    for index, value := range record {
      if value != "" {
        cType, exist := columnConfig[index + 1]
        if exist {
          // ok, clean
          switch cType {
          case "last4x":
            newRecord[index] = value[:(len(value) - 4)] + "XXXX"
          case "last4z":
            newRecord[index] = value[:(len(value) - 4)] + "0000"
          case "allhash":
            shaEncoder := sha1.New()
            shaEncoder.Write([]byte(value))
            newRecord[index] = fmt.Sprintf("%x", shaEncoder.Sum(nil))
          case "allx":
            allx := func(char rune) rune {
              return 'X'
            }
            xArray := strings.Map(allx, value)
            newRecord[index] = xArray
          default:
            newRecord[index] = value[:(len(value) - 4)] + "XXXX"
          } // switch cType
        } // if exist
      } // if value != ""
    }

    // write
    csvWriter.Write(newRecord)
    csvWriter.Flush()
  }
}
