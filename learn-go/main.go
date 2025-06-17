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

const searchWord string = "python"
const baseUrl string = "https://www.saramin.co.kr/zf_user/search/recruit?recruitPageCount=100&recruitPage=1&searchword=" + searchWord
const fileName string = "jobs/" + searchWord + ".csv"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	for i := 1; i <= totalPages; i++ {
		jobs = append(jobs, getPage(i)...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create(fileName)
	checkErr(err)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"id", "title", "location"}
	headerErr := w.Write(headers)
	checkErr(headerErr)

	for _, job := range jobs {
		csvLine := []string{job.id, job.title, job.location}
		csvError := w.Write(csvLine)
		checkErr(csvError)
	}
}

func getPage(page int) []extractedJob {
	pageUrl := baseUrl + "&recruitPage=" + strconv.Itoa(page)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkResponseCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	jobs := extractJobs(doc)
	return jobs
}

func extractJobs(doc *goquery.Document) []extractedJob {
	jobs := []extractedJob{}
	doc.Find(".item_recruit").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("value")
		title := cleanString(s.Find(".area_job").Find(".job_tit").Find("a").Text())
		location := cleanString(s.Find(".area_job").Find(".job_condition>span>a").Text())

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

func getPages() int {
	pages := 1
	res, err := http.Get(baseUrl)

	checkErr(err)
	checkResponseCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, a *goquery.Selection) {
			if !strings.Contains(a.Text(), "다음") {
				pages++
			}
		})
	})

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
