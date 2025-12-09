package main

import (
	"sync"
	"fmt"
	"time"
)

type Cashe struct{
	mu sync.RWMutex
	store map[string]string
}

func (c *Cashe) Set(key,val string){
	c.mu.Lock()
	fmt.Println("read: ",key)
	c.store[key]=val
}

func (c *Cashe) Get(key string) string{
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.store[key]
}

func main(){
	cashe:= Cashe{store: make(map[string]string)}

	go func(){
		cashe.Set("user","Admin")
	}()

	time.Sleep(10* time.Millisecond)

	for i:=0;i<3;i++{
		go func(id int){
			val:=cashe.Get("user")
			fmt.Printf("reader %d get %s\n",id,val)
		}(i)
	}

	time.Sleep(3*time.Second)
}
