package url-example

import (
	"fmt"
	"net/http"
	"errors"
	//"time"
)

type result struct {
	url string
	status string
}

var errRequestFailed = errors.New("Request failed")


func main(){
	
	//channel result 
	c := make(chan result)

	urls := []string{
	"https://www.airbnb.com/",
	"https://www.google.com/",
	"https://www.amazon.com/",
	"https://www.reddit.com/",
	"https://www.google.com/",
	"https://soundcloud.com/",
	"https://www.facebook.com/",
	"https://www.instagram.com/",
	}
	
	for _,url := range urls {
		go hitURL(url, c)
	}
	//this is same as fmt.Println(<-c) nine time
	for i := 0; i<len(urls); i++{
		fmt.Println(<-c)
	}

}

func hitURL(url string, c chan<- result){
	fmt.Println("Cheking:",url)
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400{
		status = "FAILED"
	}
	c <- result{url:url, status: status}
}

