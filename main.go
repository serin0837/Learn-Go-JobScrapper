package main

import ("github.com/serin0837/learngo/scrapper"
"github.com/labstack/echo"
"strings"
//"fmt"
"os"
)


func handleHome(c echo.Context) error{
	//send home html
	return c.File("home.html")
}

//use clean string
func handleScrape(c echo.Context)error{
	//when user download file I want to remove that file in backend
	defer os.Remove("jobs.csv")
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	//fmt.Println(c.FormValue("term"))//need name in input 
	scrapper.Scrape(term)
	//retun file 
	return c.Attachment("jobs.csv","jobs.csv")
}

func main(){
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	scrapper.Scrape("term")
	e.Logger.Fatal(e.Start(":1323"))
}
