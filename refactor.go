package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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
	rand.Seed(time.Now().Unix())
	startTime := time.Now()
	const usersCount, workersCount = 100, 20

	users := make(chan User, usersCount)
	save := make(chan string)

	for i := 0; i < workersCount; i++ {
		go worker(i+1, users, save)
	}

	generateUsers(usersCount, users)
	for i := 0; i < usersCount; i++ {
		fmt.Printf(<-save)
	}

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func worker(id int, users <-chan User, save chan<- string) {
	for user := range users {
		save <- saveUserInfo(user)
	}
}

func saveUserInfo(user User) string {
	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second) // функция долго думать
	_, err = file.WriteString(user.getActivityInfo())
	if err != nil {
		return fmt.Sprintf("Unable to write file for UID %d\n", user.id)
	}
	return fmt.Sprintf("Wrote file for UID %d\n", user.id)
}

func generateUsers(count int, users chan<- User) {
	// генерируя пользователя, посылаем его в канал из кого воркеры будут брать данные для сохранения
	for i := 0; i < count; i++ {
		// генерируем в рутиах
		go func(i int, users chan<- User) {
			users <- User{
				id:    i + 1,
				email: fmt.Sprintf("user%d@company.com", i+1),
				logs:  generateLogs(rand.Intn(1000)),
			}
			fmt.Printf("generated user %d\n", i+1)
			time.Sleep(time.Millisecond * 100)
		}(i, users)
	}
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
