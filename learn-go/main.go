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
	id       string
	title    string
	location string
}

const searchWord string = "flutter"
const baseUrl string = "https://www.saramin.co.kr/zf_user/search/recruit?recruitPageCount=100&searchword=" + searchWord
const fileName string = "jobs/" + searchWord + ".csv"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	// 페이지 결과를 받을 채널
	pageChannel := make(chan []extractedJob)

	fmt.Printf("총 %d 페이지를 스크래핑합니다.\n", totalPages)

	for i := 1; i <= totalPages; i++ {
		// getPage 함수를 고루틴으로 실행 (1단계 병렬화)
		go getPage(i, pageChannel)
	}

	// 모든 페이지로부터 결과를 수집
	for range totalPages {
		extractedJobs := <-pageChannel
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("완료, 총", len(jobs), "개의 직업을 추출했습니다.")
}

func writeJobs(jobs []extractedJob) {
	if _, err := os.Stat("jobs"); os.IsNotExist(err) {
		os.Mkdir("jobs", 0755)
	}

	file, err := os.Create(fileName)
	checkErr(err)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location"}
	headerErr := w.Write(headers)
	checkErr(headerErr)

	for _, job := range jobs {
		jobLink := job.id
		csvLine := []string{jobLink, job.title, job.location}
		csvError := w.Write(csvLine)
		checkErr(csvError)
	}
}

func getPage(page int, pageChannel chan<- []extractedJob) {
	pageUrl := baseUrl + "&recruitPage=" + strconv.Itoa(page)
	fmt.Println("Requesting:", pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkResponseCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// extractJobs는 이제 내부적으로 고루틴을 사용
	jobs := extractJobs(doc)
	pageChannel <- jobs
}

// extractJobs는 이제 페이지 내의 각 공고를 병렬로 처리하는 오케스트레이터가 됩니다.
func extractJobs(doc *goquery.Document) []extractedJob {
	var jobs []extractedJob
	// 각 채용공고 카드를 개별 고루틴에서 처리하기 위한 채널
	jobChannel := make(chan extractedJob)
	cards := doc.Find(".item_recruit")

	cards.Each(func(i int, card *goquery.Selection) {
		// 각 카드에 대해 extractJob 함수를 고루틴으로 실행 (2단계 병렬화)
		go extractJob(card, jobChannel)
	})

	// 모든 고루틴으로부터 결과를 수집
	for range cards.Length() {
		job := <-jobChannel
		jobs = append(jobs, job)
	}

	return jobs
}

// extractJob은 단일 채용 공고 카드를 파싱하여 결과를 채널로 보냅니다.
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("value")
	title := cleanString(card.Find(".job_tit a").Text())
	location := cleanString(card.Find(".job_condition span a").Text())

	// 파싱된 결과를 채널로 전송
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	var pages int
	firstPageUrl := baseUrl + "&recruitPage=1"
	res, err := http.Get(firstPageUrl)

	checkErr(err)
	checkResponseCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination a").Each(func(i int, s *goquery.Selection) {
		pageNumber, err := strconv.Atoi(s.Text())
		if err == nil && pageNumber > pages {
			pages = pageNumber
		}
	})

	if pages == 0 {
		return 1
	}
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkResponseCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("Request failed with Status:", res.StatusCode)
	}
}
