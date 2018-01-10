package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"labix.org/v2/mgo"
	"os"
	"strings"
	"time"
	//        "labix.org/v2/mgo/bson"
)

type Ball struct {
	Issue    int
	Red      []string
	Green    []string
	Datetime string
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("发生错误, ERROR:", err)
		}
	}()

	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ball").C("record")

	result := Ball{}
	c.Find(nil).Sort("-issue").One(&result)
	maxIssue := result.Issue
	if maxIssue == 0 {
		maxIssue = 2016000
	}
	year := time.Now().Year()
	maxYear := maxIssue / 1000
	fmt.Println("currentYear:", year)
	fmt.Println("maxYear:", maxYear)
	if year > maxYear {
		maxIssue = (maxYear + 1) * 1000
	}
	datetime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("datetime:", datetime)
	fmt.Println("maxIssue: ", result.Issue)

	issue := maxIssue + 1
	url := "http://caipiao.163.com/award/ssq/" + fmt.Sprintf("%d", issue) + ".html"
	fmt.Println("url:", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	var redball []string
	var blueball []string
	doc.Find("#zj_area span.red_ball").Each(func(i int, s *goquery.Selection) {
		str := strings.TrimSpace(s.Text())
		if str != "" {
			redball = append(redball, str)
		}
	})
	doc.Find("#zj_area span.blue_ball").Each(func(i int, s *goquery.Selection) {
		str := strings.TrimSpace(s.Text())
		if str != "" {
			blueball = append(blueball, s.Text())
		}
	})

	fmt.Println(redball)
	fmt.Println(blueball)
	len := len(redball)
	fmt.Println("len:", len)
	if len <= 0 {
		fmt.Println("not data!")
		os.Exit(1)
	}

	err = c.Insert(&Ball{issue, redball, blueball, datetime})
	if err != nil {
		panic(err)
	}

}
