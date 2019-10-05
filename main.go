package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	_ = csvFilename
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll() //csvファイルのデータが小さいので、メモリを大量に消費する恐れはないので読み込む
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problemloop:
	for i, p := range problems {
		fmt.Printf("第%d問: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer) //時間が過ぎても止め続ける
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\n%d問中%d問正解です。 \n", len(problems), correct)
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			} else {
				fmt.Println("不正解です。")
			}
		}
	}
	fmt.Printf("\n%d問中%d正解です。 \n", len(problems), correct)

}

func parseLines(lines [][]string) []problem {
	fmt.Println(len(lines))
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]), //スペースがcsvに入っているときのため
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
