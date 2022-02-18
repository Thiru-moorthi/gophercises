package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func quiz(csv_file *string) {

	qa_file, err := os.Open(*csv_file)

	if err != nil {
		fmt.Println(err)
	}
	defer qa_file.Close()

	csv_reader := csv.NewReader(qa_file)

	var counter, score uint

	for {
		qa, err := csv_reader.Read()
		if err == io.EOF {
			fmt.Println("You scored", score, "of", counter)
			return
		}

		if err != nil {
			fmt.Println("Error Reading CSV file!!")
		}

		fmt.Printf("Problem #%d: %s = ", counter, qa[0])

		var ans string

		fmt.Scanln(&ans)

		if ans == qa[1] {
			score++
		}
		counter++
	}
}

func main() {

	csv_file := flag.String("csv", "problems.csv", "a csv file with qestions and answer")
	flag.Parse()
	quiz(csv_file)

}
