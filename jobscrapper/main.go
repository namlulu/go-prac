package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	company    string
	jobTitle   string
	salary     string
	location   string
	experience string
	startDate  string
	endDate    string
}

var baseURL = "https://job.seoul.go.kr/www/job_offer_info/JobOfferInfo.do?method=selectJobOfferInfo&rcritJssfcCmmnCodeSe=&indutyCmmnCodeSe=&rcritJssfcCmmnCodeSe1=&jobCode11=&jobCode33=&jobCode55=&gradeFlag1=A&sido1=&sido21=&sexdstnAt1=&ageCo11=&ageCo21=&acdmcrCmmnCodeSe11=&hopeWageAmountCo1=&hopeMntslrAmountCo1=&hopeDailyAmountCo1=&hopePymhrAmountCo1=&upayLow1=&upayHigh1=&emplymStleCmmnCodeSe1=&workTmCmmnCodeSe1=&dispatchWorkAt1=&careerCndCmmnCodeSe1=&cmputrPrcuseCmmnCodeSe1=&disableEmpHopeYnA=&workDfkCmmnCodeSe1=&joFeinsrSbscrbNm1=&joReqstNo=&mberSn=&registDtHmFrom=&registDtHmTo=&registDtHmPeriod=&gradeFlag=A&sido=&sido2=&sido3=&sido4=&emplCo=&cmpnyGb=&acdmcrCmmnCodeSe=&careerCndCmmnCodeSe=&hopeAmountCo=&upayLow=&upayHigh=&preferentialGbn=&welfareCmmnCodeSeDiv=0&welfareCmmnCodeSe=&bdCmmnSeDiv=0&bdCmmnSe=&gradeSearch="

func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	pages := getPages(baseURL)

	for i := 1; i <= pages; i++ {
		go getPage(i, c)
	}

	for i := 1; i <= pages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
		fmt.Println("Page", i, "done")
	}

	writeJobs(jobs)
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func extractJob(s *goquery.Selection, c chan<- extractedJob) {
	job := extractedJob{}
	job.company = cleanString(s.Find("td").Eq(0).Text())
	job.jobTitle = cleanString(s.Find("td.multi_subject strong").Text())
	job.salary = cleanString(s.Find("td.multi_subject span").Eq(0).Text())
	job.location = cleanString(s.Find("td.multi_subject span").Eq(1).Text())
	job.experience = cleanString(s.Find("td.multi_subject span").Eq(2).Text())
	job.startDate = cleanString(s.Find("td").Eq(3).Text())
	job.endDate = cleanString(s.Find("td").Eq(4).Text())
	c <- job
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Company", "Job Title", "Salary", "Location", "Experience", "Start Date", "End Date"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.company, job.jobTitle, job.salary, job.location, job.experience, job.startDate, job.endDate}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	url := baseURL + "&pageIndex=" + strconv.Itoa(page)
	fmt.Println("Requesting:", url)
	res, err := http.Get(url)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tb_bbs.list tbody tr")
	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = (s.Find("a").Length() - 4)
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
