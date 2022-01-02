package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

var randomize bool
var filepath string
var cScore = make(chan string)
var quit = make(chan bool)
var timelimit int

type Quiz struct {
	Question string
	Answer   string
}

func read_csv(filepath string) (quizfragen []Quiz) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	for _, row := range records {
		quiz := Quiz{
			Question: row[0],
			Answer:   strings.TrimSpace(row[1]),
		}
		quizfragen = append(quizfragen, quiz)
	}
	defer file.Close()
	return
}

func request(quizfrage string) {
	reader := bufio.NewReader(os.Stdin)
	antwort, _ := reader.ReadString('\n')
	antwort = strings.Replace(antwort, "\n", "", -1)
	// Remove \r from Answer when on windows
	if runtime.GOOS == "windows" {
		antwort = strings.Replace(antwort, "\r", "", -1)
	}
	cScore <- antwort
	return
}

func init() {
	flag.BoolVar(&randomize, "randomize", false, "randomize ofer of question")
	flag.StringVar(&filepath, "file", "problems.csv", "Path to file")
	flag.IntVar(&timelimit, "timer", 10, "Time for the whole Quiz")
}

func timer(timelimit int) {
	time.Sleep(time.Duration(timelimit) * time.Second)
	quit <- true
	return
}

func main() {
	flag.Parse()
	var quizfragen []Quiz
	if filepath != "problems.csv" {
		quizfragen = read_csv(filepath)
	} else {
		quizfragen = read_csv("problems.csv")
	}

	if randomize {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quizfragen), func(i, j int) {
			quizfragen[i], quizfragen[j] = quizfragen[j], quizfragen[i]
		})
	}

	score := 0
	overall := len(quizfragen)
	for _, quizfrage := range quizfragen {
		fmt.Println("Question:", quizfrage.Question)
		go timer(timelimit)
		go request(quizfrage.Answer)

		select {
		case answer := <-cScore:
			if answer == quizfrage.Answer {
				score += 1
			}
		case <-quit:
			fmt.Printf("\nGame Over, your Score:%d/%d\n", score, overall)
			os.Exit(0)
		}
	}
}
