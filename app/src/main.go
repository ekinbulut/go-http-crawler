package main

import (
	"fmt"
	"time"

	"flag"

	service "github.com/ekinbulut/go-http-crawler/app/srv"
)

// flag -c=https://lego.storeturkey.com.tr/technic?ps=4
// flag -c=https://lego.storeturkey.com.tr/technic?ps=4 -o=.\lego.html

var site string
var outputFile string
var interval int

type App struct {
	Name        string
	Version     string
	crawler     *service.Crawler
	fileWriter  *service.FileWriter
	diffChecker *service.DiffChecker
}

func NewApp() *App {
	return &App{
		Name:        "go-http-crawler",
		Version:     "1.0.0",
		crawler:     service.NewCrawler(site),
		fileWriter:  service.NewFileWriter(outputFile),
		diffChecker: service.NewDiffChecker(),
	}
}

func (a *App) Run() {

	a.printAppInfo()
	parseFlags()
	// print flags
	fmt.Println("site:", site)
	fmt.Println("outputFile:", outputFile)

	// execute in given interval
	if interval > 0 {
		for {
			a.execute()
			fmt.Println("sleeping...")
			time.Sleep(time.Duration(interval) * time.Second)
		}
	} else {
		a.execute()
	}

}

func (app *App) execute() {
	// print progress
	fmt.Println("crawling...")

	resp, err := app.crawlsite(site)
	if err != nil {
		fmt.Println(err)
	}

	b := app.fileWriter.Exists()
	if b {
		// read file
		old, err := app.fileWriter.Read()
		if err != nil {
			fmt.Println(err)
		}
		// check diff
		b, err := app.diffChecker.Check(old, resp)
		if err != nil {
			fmt.Println(err)
		}
		if b {
			fmt.Println("no changes")
		} else {
			fmt.Println("changes found")
			err := app.createOutput(resp)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		err := app.createOutput(resp)
		if err != nil {
			fmt.Println(err)
		}
	}

	// print progress
	fmt.Println("done")

}

// print App info
func (a *App) printAppInfo() {
	fmt.Printf("%s %s\n", a.Name, a.Version)
}

func (app *App) crawlsite(site string) (string, error) {

	resp, err := app.crawler.Crawl()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return resp, nil
}

func (app *App) createOutput(resp string) error {
	f := app.fileWriter
	err := f.Write(resp)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = f.Rename(outputFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	app := NewApp()
	app.Run()

}

// parse flags
func parseFlags() {
	flag.StringVar(&site, "u", "", "u=https://sample.com")
	flag.StringVar(&site, "url", "", "url=https://sample.com")
	flag.StringVar(&outputFile, "o", "", "o=output.html")
	flag.IntVar(&interval, "i", 0, "i=1")

	flag.Parse()
}
