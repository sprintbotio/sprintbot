package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, playground")
	t := time.Now()
	//
	//Mon Jan _2 15:04:05 2006
	//Fri Dec 28 10:10:00 2018
	timeF := fmt.Sprintf("%s %s %d %s:00 %d",t.Weekday().String()[0:3],t.Month().String()[0:3],t.Day(), "10:10", t.Year())
	fmt.Println(timeF)
	//t1 , err := time.ParseInLocation(time.ANSIC,"Fri Dec _28 10:10:00 2018" , l1)
	//if err != nil{
	////	panic(err)
	//}
	//t2 , err := time.ParseInLocation(time.RFC3339,timeF, l2)
	//if err != nil{
	////	panic(err)
	//}
	l, e := time.LoadLocation("Antarctica/Palmer")
	if e!= nil{
		panic(e)
	}
	p := fmt.Println
	time.Local = l

	t3, _ := time.Parse(
		time.RFC3339,
		fmt.Sprintf("%d-%02d-%02dT10:10:00+00:00", t.Year(),t.Month(),t.Day()))

	p(t3.Unix())

	//fmt.Println(t2.Unix(), t1.Unix())
}