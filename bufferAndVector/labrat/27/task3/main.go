package main
import "fmt"

func main() {
    ch := make(chan int)
    go func() {
        // Хотим отправить 3 числа
        ch <- 1
        ch <- 2
        ch <- 3
        close(ch) 
    }()

    // Пытаемся прочитать больше, чем отправили
	
    fmt.Println(<-ch)
    fmt.Println(<-ch)
    fmt.Println(<-ch)
	val , ok :=<-ch
	if (!ok){
		fmt.Println("канал закрыт")
		return;
	}
    fmt.Println(val," ",ok)
	//выводит 0, потому что мы туда ничего не отправляли, null/nil/NULL
}


/* код с какой-то ошибкой
package main
import "fmt"

func main() {
    ch := make(chan int)
    go func() {
        // Хотим отправить 3 числа
        ch <- 1
        ch <- 2
        ch <- 3
        close(ch) 
    }()

    // Пытаемся прочитать больше, чем отправили
    fmt.Println(<-ch)
    fmt.Println(<-ch)
    fmt.Println(<-ch)
    fmt.Println(<-ch) // <-- Что тут произойдет и почему? Исправьте, используя v, ok
}

*/