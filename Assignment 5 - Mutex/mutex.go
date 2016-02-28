package main

import ( 
	"fmt"
	"sync"
	"time"
)

//Initialize global variables
var masterData int = 0
var mutExLock = &sync.Mutex{}

func main() {
		
		//Fire up a thousand read and write threads
		for i:=0; i < 1000; i++ {
			go readMasterData(i)
			go writeMasterData(i)
			time.Sleep(1 * time.Second)
		}
		
}

func readMasterData(i int) {
	//Acquire Lock  
	mutExLock.Lock()
	fmt.Printf("Masterdata currently being read by ReadThread(%d), and the value is: %d \n",i,masterData)
	mutExLock.Unlock()
}

func writeMasterData(i int) {
	mutExLock.Lock()
	masterData += i
	fmt.Printf("Masterdata currently being written by WriteThread(%d), and the value is: %d \n",i,masterData)
	time.Sleep(2 * time.Second)
	mutExLock.Unlock()
}
