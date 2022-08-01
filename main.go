package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	startTime := time.Now()
	var count = 100
	rand.Seed(time.Now().Unix())

	genJobs := make(chan int, count)
	putJobs := make(chan User, count)

	var wg = &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		go genWorker(genJobs, putJobs)
		go putWorker(putJobs, wg)
	}
	generateUsers(count, genJobs, wg)
	wg.Wait()
	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func genWorker(genJobs <-chan int, putJobs chan<- User) {
	for {
		id, ok := <-genJobs
		if !ok {
			return
		}
		generateUserInfo(id, putJobs)
	}
}

func putWorker(putJobs <-chan User, wg *sync.WaitGroup) {
	for {
		user, ok := <-putJobs
		if !ok {
			return
		}
		saveUserInfo(user)
		wg.Done()
	}
}

func saveUserInfo(user User) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

func generateUserInfo(id int, putJobs chan<- User) {
	putJobs <- User{
		id:    id,
		email: fmt.Sprintf("user%d@company.com", id),
		logs:  generateLogs(rand.Intn(1000)),
	}
	fmt.Printf("generated user %d\n", id)
	time.Sleep(time.Millisecond * 100)
}

func generateUsers(count int, genJobs chan<- int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		genJobs <- i + 1
	}
	close(genJobs)
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
