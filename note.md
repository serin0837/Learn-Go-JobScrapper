- go lang std library 
- http package

- it's not that fast to go through of urls


- creat map 
    ```go
    var results map[string]string
	results["gello"] = "hello"
    ```
    above throw error (panic) because map is nil

    I can not write in uninitialize map 

    I have to put {}
    ```go
    var results = map[string]string{}
	results["gello"] = "hello"
    ```
    or can use make function then can initialize empty map
    ```go
    var results = make(map[string]string)
    ```

    ```go
    func main(){
	var results = make(map[string]string)

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
		result := "OK"
		err := hitURL(url)
		if err != nil{
			result = "FAILED"
		}
		results[url] = result
	}
	fmt.Println(results)
       }
    ```


can do with python as well 
and there was 429 error which is too many request

we don't want go check things in order 
we wnat to chack things all together


- Go concurrency 
```go
func sexyCount(person string){
	for i := 0; i <10 ; i++{
		fmt.Println(person, "is sexy", i)
		time.Sleep(time.Second)
	}
}

	sexyCount("nico")
	sexyCount("flynn")
    every second print 
    //nico is sexy 0
    // nico is sexy 1
    // nico is sexy 2
    // nico is sexy 3
    // nico is sexy 4
    // nico is sexy 5
    //...
```
now it take 20 secound to finish main func

- goroutines!!! just add go keyword
```go
	go sexyCount("nico")
	sexyCount("flynn")
//nico is sexy 5
// flynn is sexy 6
// nico is sexy 6
// nico is sexy 7
// flynn is sexy 7
// flynn is sexy 8
// nico is sexy 8
// nico is sexy 9
// flynn is sexy 9
//...
```
- when I add both go, program just finish and nothing happen 
- becuase main fuc have nothing to do 

- how to adpt to url checker?
- need to find way to talk to main func



- Channel
gorouitne communicate
channel is like a pipe
```go
//boolean is value of goroutine's func//what type of information we want to send to main func
func main(){
    c := make(chan bool)
    peopler := [2]string{"nico", "flynn"}
    for _, person := range people {
        go isSexy(person, c)
    }
    //receive message throuh channel
    result := <-c
    fmt.Println(result)
    fmt.Println(<-c)
}

func isSexy(person string, c chan bool){
    time.Sleep(time.Second * 5)
    //to the channel c send true
    c <- true 
}
```

- when you recieve something from channel, main function will wait until go get one message 
- no return keyword
- chan : type of chan 
- waiting for message is blocking 
- so then we can use loop and it's happening concurrently

- rule of channel
    - main func finish, then go routine die
    - I have to specified the type of data that going to send and recive
    - channle <- message : send mesage
    - <- channel : receiving a message, but this it blocking 


- now change main code!
before code
```go
package main

import (
	"fmt"
	"net/http"
	"errors"
	"time"
)

var errRequestFailed = errors.New("Request failed")


func main(){
	var results = make(map[string]string)

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
		result := "OK"
		err := hitURL(url)
		if err != nil{
			result = "FAILED"
		}
		results[url] = result
	}
	//fmt.Println(results)
	//format better
	for url, result := range results{
		fmt.Println(url, result)
	}

}

func hitURL(url string) error{
	fmt.Println("Cheking:",url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400{
		fmt.Println(err, resp.StatusCode)
		return errRequestFailed
	}
	return nil
}

```

- we going to send chaneel to hitURL

```go
package main

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
```

really fast!!!!!!!!
- take slowes url to finish time, and not waiting each url


- go query
intall go get github.com/PuerkitoBio/goquery