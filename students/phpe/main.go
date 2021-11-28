package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// - wir wollen eine Datei via Flag öffnen
// - diese Auslesen und die Fragen speichern --> erledigt
// - die Fragen durchgehen und prüfen ob die Antworten übereinstimmen --> erledigt

type Quiz struct {
	Frage   string
	Antwort string
}

func read_csv() (quizfragen []Quiz) {
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	for _, row := range records {
		quiz := Quiz{
			Frage:   row[0],
			Antwort: strings.TrimSpace(row[1]),
		}
		quizfragen = append(quizfragen, quiz)
	}
	defer file.Close()
	return
}

func abfragen(quizfragen []Quiz) {
	reader := bufio.NewReader(os.Stdin)
	for _, quizfrage := range quizfragen {
		fmt.Println("Lösung zu:", quizfrage.Frage)
		antwort, _ := reader.ReadString('\n')
		antwort = strings.Replace(antwort, "\n", "", -1)
		// Remove \r from Answer when on windows
		if runtime.GOOS == "windows" {
			antwort = strings.Replace(antwort, "\r", "", -1)
		}
		if quizfrage.Antwort == antwort {
			fmt.Println("stimmt")
		} else {
			fmt.Println("Das war falsch")
		}
	}
}

func main() {
	quizfragen := read_csv()
	abfragen(quizfragen)
}
