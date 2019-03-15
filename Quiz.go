package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type QuestionAnswer struct {
	question string
	answer   string
}

var wg sync.WaitGroup

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Press return to start the quiz !")
	_, _ = reader.ReadString('\n')

	fmt.Println("Starting timer for 15 secs...")
	wg.Add(1)
	go startTimer()

	file, err := os.Open("quiz.csv")
	if err != nil {
		log.Fatal(err)
	}

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	_ = file.Close()

	var questionAnswerList []QuestionAnswer
	for _, line := range lines {
		questionAnswerList = append(questionAnswerList, QuestionAnswer{question: line[0], answer: line[1]})
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(questionAnswerList), func(i, j int) {
		temp := questionAnswerList[i]
		questionAnswerList[i] = questionAnswerList[j]
		questionAnswerList[j] = temp
	})
	score := 0
	go func() {
		for i, q := range questionAnswerList {
			fmt.Printf("Question %d: %s: ", i+1, q.question)
			ans, _ := reader.ReadString('\n')
			ans = strings.Trim(ans, "\n")
			if ans == q.answer {
				score++
			}
		}
		str := "Quiz over. Let's see how you did"
		fmt.Printf("%s\r", str)
		time.Sleep(500 * time.Millisecond)
		for i := 0; i < 3; i++ {
			str += "."
			fmt.Printf("%s\r", str)
			time.Sleep(500 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("\nYour score is %d/13\n", score)
}

func startTimer() {
	defer wg.Done()
	for i := 0; i < 15; i++ {
		time.Sleep(time.Second * 1)
	}
	fmt.Println("\n\nTime's up sucker !!!")
}
