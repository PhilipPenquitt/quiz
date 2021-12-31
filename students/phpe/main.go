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

// - wir wollen eine Datei via Flag öffnen --> erledigt
// - diese Auslesen und die Fragen speichern --> erledigt
// - die Fragen durchgehen und prüfen ob die Antworten übereinstimmen --> erledigt
// Die Fragen sollen optional via Schalter per Zufall ausgegeben werden --> erledigt
// Nach Ablauf des Timers soll ein Score aller erfolgreicher Fragen angezeigt werden
var randomize bool
var filepath string

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

func abfragen(quizfragen []Quiz) {
	reader := bufio.NewReader(os.Stdin)
	for _, quizfrage := range quizfragen {
		fmt.Println("Question:", quizfrage.Question)
		antwort, _ := reader.ReadString('\n')
		antwort = strings.Replace(antwort, "\n", "", -1)
		// Remove \r from Answer when on windows
		if runtime.GOOS == "windows" {
			antwort = strings.Replace(antwort, "\r", "", -1)
		}
		// We don't want this in part 2 of the Tasks
		// if quizfrage.Answer == antwort {
		// fmt.Println("correct")
		// } else {
		// fmt.Println("Thats wrong")
		// }
	}
}

func init() {
	flag.BoolVar(&randomize, "randomize", false, "randomize ofer of question")
	flag.StringVar(&filepath, "file", "problems.csv", "Path to file")
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
	abfragen(quizfragen)
}
