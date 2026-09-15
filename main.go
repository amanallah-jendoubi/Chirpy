package main

import (
	"time"
)

type email struct {
	body string
	date time.Time
}

func checkEmailAge(emails [3]email) [3]bool {
	isOldChan := make(chan bool)

	go sendIsOld(isOldChan, emails)

	isOld := [3]bool{}
	isOld[0] = <-isOldChan
	isOld[1] = <-isOldChan
	isOld[2] = <-isOldChan
	return isOld
}

func sendIsOld(isOldChan chan<- bool, emails [3]email) {
	for _, e := range emails {
		if e.date.Before(time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)) {
			isOldChan <- true
			continue
		}
		isOldChan <- false
	}
}

func main() {
	emails := [3]email{
		{body: "Hello", date: time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC)},
		{body: "World", date: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC)},
		{body: "!", date: time.Date(2022, 6, 15, 0, 0, 0, 0, time.UTC)},
	}

	isOld := checkEmailAge(emails)
	for i, old := range isOld {
		if old {
			println("Email", i, "is old.")
		} else {
			println("Email", i, "is not old.")
		}
	}
}
