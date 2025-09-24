package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/namlulu/web/scrapper"
)

func handleHome(c echo.Context) error {
	return c.File("index.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove("jobs.csv")

	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	fmt.Println("Scraping for term:", term)

	scrapper.Scrapper()
	return c.Attachment("jobs.csv", "jobs.csv")
}

func main() {
	// 웹 서버 초기화 및 라우팅 설정
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)

	e.Logger.Fatal(e.Start(":8080"))
}
