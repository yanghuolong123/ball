package main

import (
	"fmt"
	"github.com/widuu/gojson"
	"labix.org/v2/mgo"
	"os/exec"
	"time"
	//        "labix.org/v2/mgo/bson"
)

type Ball struct {
	Issue int
	Red   []string
	Green []string
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
	maxYear := maxIssue/1000
	if year != maxYear {
		maxIssue = maxYear * 1000
	}
	fmt.Println("maxIssue: ", result.Issue)

	issue := maxIssue + 1
	url := "http://hao123.lecai.com/lottery/draw/ajax_get_detail.php?lottery_type=50&phase=" + fmt.Sprintf("%d", issue)

	cmd := exec.Command("curl", url)
	bytes, err := cmd.Output()
	jsonData := string(bytes)
	fmt.Println("jsonData: ", jsonData)

	redball := gojson.Json(jsonData).Getpath("data", "result", "result").Getkey("data", 1).StringtoArray()
	greenball := gojson.Json(jsonData).Getpath("data", "result", "result").Getkey("data", 2).StringtoArray()

	err = c.Insert(&Ball{issue, redball, greenball})

	fmt.Println("redball:", redball)
	fmt.Println("greenball:", greenball)

}
