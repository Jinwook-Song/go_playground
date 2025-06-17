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
	// 결과를 받을 채널 생성
	c := make(chan []extractedJob)

	fmt.Printf("총 %d 페이지를 스크래핑합니다.\n", totalPages)

	for i := 1; i <= totalPages; i++ {
		// getPage 함수를 고루틴으로 실행
		go getPage(i, c)
	}

	// 모든 고루틴으로부터 결과를 수집
	for range totalPages {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("완료, 총", len(jobs), "개의 직업을 추출했습니다.")
}

func writeJobs(jobs []extractedJob) {
	// jobs 디렉토리가 없으면 생성
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
		// 링크를 포함한 ID로 변경
		jobLink := "https://www.saramin.co.kr/zf_user/jobs/relay/view?rec_idx=" + job.id
		csvLine := []string{jobLink, job.title, job.location}
		csvError := w.Write(csvLine)
		checkErr(csvError)
	}
}

// getPage는 이제 채널을 인자로 받아 결과를 채널로 보냅니다.
func getPage(page int, c chan<- []extractedJob) {
	pageUrl := baseUrl + "&recruitPage=" + strconv.Itoa(page)
	fmt.Println("Requesting:", pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkResponseCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	jobs := extractJobs(doc)
	// 함수가 끝날 때 결과를 채널로 보냄
	c <- jobs
}

func extractJobs(doc *goquery.Document) []extractedJob {
	jobs := []extractedJob{}
	doc.Find(".item_recruit").Each(func(i int, s *goquery.Selection) {
		// 'rec_idx' 속성으로 ID를 가져오도록 수정
		id, _ := s.Attr("rec_idx")
		title := cleanString(s.Find(".job_tit a").Text())
		location := cleanString(s.Find(".job_condition span a").Text())

		job := extractedJob{
			id:       id,
			title:    title,
			location: location,
		}
		jobs = append(jobs, job)
	})
	return jobs
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// getPages 함수 로직을 더 견고하게 수정
func getPages() int {
	var pages int
	// URL에 recruitPage=1을 추가하여 첫 페이지만 확실히 가져옴
	firstPageUrl := baseUrl + "&recruitPage=1"
	res, err := http.Get(firstPageUrl)

	checkErr(err)
	checkResponseCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// 페이지네이션 영역에서 마지막 페이지 번호를 찾음
	doc.Find(".pagination a").Each(func(i int, s *goquery.Selection) {
		// 링크의 텍스트를 정수로 변환 시도
		pageNumber, err := strconv.Atoi(s.Text())
		// 에러가 없고, 현재 페이지 번호가 pages 변수보다 크면 업데이트
		if err == nil && pageNumber > pages {
			pages = pageNumber
		}
	})

	// 페이지네이션이 없는 경우 (결과가 1페이지 이하인 경우)
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
