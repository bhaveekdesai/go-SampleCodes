package main

import ( "fmt"
	"bufio"
	"net"
	"time"
	"os"
	"strings"
	"log"
	"encoding/json"
)
	
type Mailbox struct {
	Messages []MessagesStruct
}

type MessagesStruct struct {
    SentTime string
    SenderName string
    MessageContent string
}

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
	connection,error := net.Dial("tcp",serverIP+":8083")
	if error != nil {
		log.Println("Server is down!")
		os.Exit(0)
	}
	
	//Print initial messages
	fmt.Println("Welcome to Rabbl")
	//fmt.Println("say 'rabbl' to get your messages")
	
	//Initiating command
	outgoingMessage := "rabbl"
	fmt.Fprint(connection, outgoingMessage + "\n")
		
	for {
		
		//Fetch messages from server
		incomingResponse,_ := bufio.NewReader(connection).ReadString('`')
		
		responseMessage := strings.Split(incomingResponse,"`")[0]
	
		//Decode response
		noNewMessages := false
		if strings.Contains(responseMessage, "|EMPTY|") {
			noNewMessages = true
		}
		
		responseMessageParts := strings.Split(responseMessage,"|")
		greeting := responseMessageParts[0]
		messages := []byte(responseMessageParts[1])
		
		//Decode json response
		if noNewMessages == false {
			var mailbox Mailbox
			if error := json.Unmarshal(messages, &mailbox); error != nil {
				fmt.Printf("Woah! Corrupt message alert!\n")
			}
			
			//Printing messages
			fmt.Printf("%s\n",greeting)
			
			for index,element := range mailbox.Messages {
				fmt.Printf("\nMessage %d:\nTime: %s\nFrom: %s\nMessage: %s\n\n",(index+1),element.SentTime,element.SenderName,element.MessageContent)
			}
		} else {
			fmt.Println("No new messages! Relax!")
		}
	}
}
