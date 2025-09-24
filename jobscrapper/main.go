package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

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
	pages := getPages(baseURL)
	fmt.Println("Total Pages:", pages)

	// 페이지 단위 동시 실행 제한
	const MaxPageWorkers = 10
	pageSem := make(chan struct{}, MaxPageWorkers)

	c := make(chan []extractedJob)
	var wg sync.WaitGroup

	for i := 1; i <= pages; i++ {
		pageSem <- struct{}{} // 페이지 슬롯 확보
		wg.Add(1)
		go func(page int) {
			defer func() { <-pageSem }() // 슬롯 반환
			defer wg.Done()

			logRuntimeStatus(fmt.Sprintf("Page %d start", page)) // 시작
			getPage(page, c)
			logRuntimeStatus(fmt.Sprintf("Page %d end", page)) // 종료
		}(i)
	}

	// 결과 수집용 고루틴
	go func() {
		wg.Wait()
		close(c)
	}()

	for extractedJobs := range c {
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
}

// ---------- 나머지 함수 ----------

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

// getPage 안에서 공고 단위 고루틴 제한 적용
func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob

	url := baseURL + "&pageIndex=" + strconv.Itoa(page)
	fmt.Println("Requesting:", url)

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tb_bbs.list tbody tr")
	c := make(chan extractedJob)

	for i := 0; i < searchCards.Length(); i++ {
		s := searchCards.Eq(i)
		go func(s *goquery.Selection) {
			extractJob(s, c)
		}(s)
	}

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

// 고루틴 시작/끝마다 호출할 헬퍼 함수
func logRuntimeStatus(tag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("CPU cores:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	fmt.Printf("[%s] Goroutines: %d | Alloc: %d KB | TotalAlloc: %d KB | Sys: %d KB | NumGC: %d\n",
		tag, runtime.NumGoroutine(),
		m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}
