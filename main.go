package main

import "fmt"

func waitForDBs(numDBs int, dbChan chan struct{}) {
	<-dbChan
}

func getDBsChannel(numDBs int) (chan struct{}, *int) {
	count := 0
	ch := make(chan struct{})

	go func() {
		for i := 0; i < numDBs; i++ {
			ch <- struct{}{}
			fmt.Printf("Database %v is online\n", i+1)
			count++
		}
	}()

	return ch, &count
}

func main() {
	numDBs := 5
	dbChan, count := getDBsChannel(numDBs)

	for i := 0; i < numDBs; i++ {
		waitForDBs(numDBs, dbChan)
	}

	fmt.Printf("All %v databases are online\n", *count)
}
