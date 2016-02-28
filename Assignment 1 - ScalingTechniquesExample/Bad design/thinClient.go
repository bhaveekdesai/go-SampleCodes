package main

import ( "fmt"
	"bufio"
	"net"
	"time"
	"os"
	"log"
)
	
func main() {
	go begin()
	time.Sleep(10 * time.Minute)
}

func begin() {
	//Default server IP. Or take from comman-line argument
	serverIP := "127.0.0.1"
	if len(os.Args) >= 2 {
		serverIP = os.Args[1:][0]
	}
	
	//Initiate Server connection
	connection,error := net.Dial("tcp",serverIP+":8082")
	if error != nil {
		log.Println(error)
	}
	
	//Print initial messages
	fmt.Println("Thin Client started")

	for {
		//Read user entered messages
		readerHandle := bufio.NewReader(os.Stdin)
		fmt.Print("Enter email ID: ")
		outgoingMessage,error := readerHandle.ReadString('\n')
		if error != nil {
			log.Println(error)
		}
	
		//Send message to server & fetch response 
		clockStart := time.Now()
		
		fmt.Fprint(connection, outgoingMessage + "\n")
		incomingResponse,_ := bufio.NewReader(connection).ReadString('\n')
		
		latency := time.Since(clockStart)
		
		//Print response
		fmt.Printf("Response: \n-----%s-----Latency: %s\n",incomingResponse, latency)
		
	}
}
