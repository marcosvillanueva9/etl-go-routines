package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	
	"etl-sample/model"
	"etl-sample/util"
)

var yaml model.Yaml

var inputFilePath = "./default.csv"
var fileType = "csv"
var destPath = "./defaultdest.csv"

var columnMapper = make(map[int]string)
var destColumnMapper = make(map[int]string)

func main() {
	router := gin.Default()

	router.POST("/", RunJob)
	
	router.Run(":8080")
}

func RunJob(c *gin.Context) {
	// Just to see how long it took
	start := time.Now()

	//read yml
	configYaml()

	// Set up the channels.

	extractCh := make(chan []string)
	transformCh := make(chan map[string]string)
	doneCh := make(chan bool)

	go extract(extractCh)
	go transform(extractCh, transformCh)
	go load(transformCh, doneCh)

	<-doneCh
	// Show the user how long the etl process took
	fmt.Println(time.Since(start))
}

func configYaml() {
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.Unmarshal(&yaml)

	inputFilePath = yaml.Source.Connection.Path
	destPath = yaml.Destination.Connection.Path
	fileType = yaml.Source.Connection.Filetype

	// Set up columns.

	columnsOrigin := yaml.Source.Columns.Origin
	columns := strings.Split(columnsOrigin, ",")
	for i, name := range columns {
		columnMapper[i] = name
	}

	columnsDest := yaml.Source.Columns.Destination
	columnsdest := strings.Split(columnsDest, ",")
	for i, name := range columnsdest {
		destColumnMapper[i] = name
	}
}

// Extracts the data from orders.
func extract(extractCh chan []string) {
	// Get all the inform from orders.
	fmt.Println("Extract function started")
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	// Load the csv format.
	switch {
		case fileType == "csv":
			util.CsvReader(f, extractCh)
		case fileType == "json":
			util.JsonReader(f, extractCh)
		case fileType == "txt":
			util.TxtReader(f, extractCh)
		case fileType == "flat":
			util.FlatReader(f, extractCh)
		default:
			log.Fatal("Filetype not supported")
	}
	
	// Close the channel.
	close(extractCh)
}

// Transform the orders that are in the extract channel.
func transform(extractCh chan []string, transformCh chan map[string]string) {
	fmt.Println("Transform function started")
	
	var destColumns = make(map[string]string)
	columnsDest := yaml.Source.Columns.Destination
	columnsdest := strings.Split(columnsDest, ",")
	for _, name := range columnsdest {
		destColumns[name] = ""
	}

	// Set up a wait group.
	var waitGrp sync.WaitGroup
	// Go trough every record in the extract channel.
	for o := range extractCh {
		// Add a wait group.
		waitGrp.Add(1)
		// Start a goroutine that will transform the record.
		// For each record in the extract channel.
		go func(o []string, destination map[string]string, connectors []model.Connectors) {

			for _, connector := range connectors {
			 	switch connector.Connector.Name {
			 	case "trim":
			 		util.Trim(o, destination, connector.Connector.Params)
			 	case "parse":
			 		util.Parse(o, destination, connector.Connector.Params)
			 	case "concat":
			 		util.Concat(o, destination, connector.Connector.Params, columnMapper)
			 	}
			}

			transformCh <- destination
			waitGrp.Done()
		}(o, destColumns, yaml.Transform)
	}

	waitGrp.Wait()

	close(transformCh)
}


// Load the data into a text file.
func load(transformCh chan map[string]string, doneCh chan bool) {
	fmt.Println("Load function started")
	f, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Add a header to the text file.
	_, err = fmt.Fprintf(f, "%s\n", yaml.Source.Columns.Destination)

	if err != nil {
		log.Fatal(err)
	}

	// Set up a wait group.
	var waitGrp sync.WaitGroup
	// Go trough every record in the channel.
	for o := range transformCh {
		waitGrp.Add(1)
		go func(o map[string]string) {
			for _, columnname := range destColumnMapper {
				_, err = fmt.Fprintf(f, "%s,", o[columnname])
			}
			_, err = fmt.Fprintf(f, "%s", "\n")
			if err != nil {
				log.Fatal(err)
			}
			waitGrp.Done()
		}(o)

	}

	waitGrp.Wait()

	// Tell the main function that we are done.
	doneCh <- true
}
