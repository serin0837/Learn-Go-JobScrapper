package main

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"os"
	"log"
	"strconv"
	"strings"
	"encoding/csv"
)

//create struct
type extractedJob struct {
	id string
	location string
	title string
	salary string
	summary string
}
//only start is change when I go to different page
var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main(){
	var jobs []extractedJob
	//1.make channel-main
	c := make(chan []extractedJob)
	totalPages := getPages()
	//fmt.Println(totalPages)//5
	//hit the url 
	for i := 0; i < totalPages; i++{
		//2. go routine-main
		go getPage(i, c)
		

	}
	//5. recive message of mainc-main
	for i := 0 ; i<totalPages;i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}
	//fmt.Println(jobs)
	writeJobs(jobs)
	fmt.Println("Done")
}

//write jobs /save in csv file

func writeJobs(jobs []extractedJob){
	//csv package
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	//write data to that file 
	defer w.Flush()

	//headers
	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	//for loop job from jobs
	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjobs?jk="+job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}


//create get page func
//3. channel-main
func getPage(page int, mainC chan<- []extractedJob) {
	//create empty jobs
	var jobs []extractedJob
	//1. make channel
	c := make(chan extractedJob)

	pageURL := baseURL + "&start=" + strconv.Itoa(page * 50)
	//page*50 is number so we have to use pacakage
	//print its wokring 
	fmt.Println("Resuqesiting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	//goguery document//resbody is byte so we have to close 
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection){
		//each card I want to extract job
		//s is each card (div)
		//each card have id(data-jk)
		//attr return two values, actual value and existence
		// id, _ := s.Attr("data-jk")
		// //print all of id
		// fmt.Println(id)
		// title := cleanString(s.Find(".title>a").Text())
		
		// location := cleanString(s.Find(".sjcl").Text())
		// fmt.Println(id, title, location)
		
		//separate extractjob function //4. add go 
		go extractJob(card, c)// 2. add channel
	
		//jobs = append(jobs, job)
	})
		//5. recevie channel message 
		for i :=0 ; i <searchCards.Length(); i++{
			job := <-c
			jobs = append(jobs, job)
		}
		//4. not retun job we going to send job to main chaneel -main
	mainC <- jobs
}
//crete function taht only extracting job//3. c chan<-extractedJob// get rid of return value
func extractJob(card *goquery.Selection, c chan<- extractedJob){
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text()) 
	//fmt.Println(id, title, location, salary, summary)
	c <- extractedJob{
		id:id, 
		title: title, 
		location: location, 
		salary: salary, 
		summary:summary,
	}
	
}

// create function that trim whitespace 
func cleanString(str string) string{
	//filed change string to array of string and join
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int{
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	//goguery document//resbody is byte so we have to close 
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	//fmt.Println(doc)
	//pagination is class in that webpage
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection){
		//I can see every tag of that class
		//fmt.Println(s.Html())

		// we want count links with a tag
		//fmt.Println(s.Find("a").Length())
		pages = s.Find("a").Length()
	})
	
	return pages
}

func checkErr(err error){
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response){
	if res.StatusCode != 200{
		log.Fatalln("request failed with Status:", res.StatusCode)
	}
}