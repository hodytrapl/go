package main

import (
	"flag"
	"fmt"
	"time"
)

const (
	messagePerGorutine = 3
	timePerSecond      = 5
)

/*func homework(id int) {
	for i := 1; i <= id; i++ {
		fmt.Printf("[gorutine %d] message %d\n", id, i)
	}
}*/

func homework() {
	for i := 1; i <= messagePerGorutine; i++ {
		fmt.Printf("[gorutine %d] message %d\n", messagePerGorutine, i)
	}
}

func main() {
	/*if len(os.Args) != 2 {
		log.Fatal("Ожидается аргумент: количетсво горутин\n")
	}
	countgoroutline, err := strconv.Atoi(os.Args[1])
	if err != nil || countgoroutline <= 0 {
		log.Fatal("аргумент не положительное число\n")
	}
	for i := 0; i < countgoroutline; i++ {
		go homework(i)
	}*/

	count := flag.Int("n", 10, "флаг который указывает значение по умолчанию 10")
	flag.Parse()

	for i := 0; i < *count; i++ {
		go homework()
	}

	time.Sleep(time.Duration(*count) * timePerSecond)
}
