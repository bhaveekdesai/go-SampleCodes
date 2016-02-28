package main

import ( "fmt"
	"bufio"
	"net"
	"time"
	"os"
	"strings"
	"regexp"
	"log"
	"encoding/json"
)
	
type Mailbox struct {
	Messages []MessagesStruct
}

type MessagesStruct struct {
    SenderName    string
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
		log.Println(error)
	}
	
	//Print initial messages
	fmt.Println("Welcome to Rabbl")
	fmt.Println("say 'rabbl' to get your messages")


	for {
		//Read user entered messages
		readerHandle := bufio.NewReader(os.Stdin)
		outgoingMessage,_ := readerHandle.ReadString('\n')
		outgoingMessage = strings.Split(string(outgoingMessage),"\n")[0]
		
		//Validate
		valid := validate(string(outgoingMessage))
		responseMessage := "Couldn't hear you. Say rabbl to get your messages"
		latency := time.Since(time.Now())
		
		//Send message to server & fetch response if validation successful
		if valid == true {
			clockStart := time.Now()
			
			fmt.Fprint(connection, outgoingMessage + "\n")
			incomingResponse,_ := bufio.NewReader(connection).ReadString('`')
			
			latency = time.Since(clockStart)
			responseMessage = strings.Split(incomingResponse,"`")[0]
		
			//Decode response
			responseMessageParts := strings.Split(responseMessage,"|")
			greeting := responseMessageParts[0]
			messages := []byte(responseMessageParts[1])
			
			//Decode json response
			var mailbox Mailbox
			if error := json.Unmarshal(messages, &mailbox); error != nil {
				fmt.Printf("Woah! Corrupt message alert!\n")
			}
			
			//Printing messages
			fmt.Printf("%s\n",greeting)
			
			for index,element := range mailbox.Messages {
				fmt.Printf("\nMessage %d:\nFrom: %s\nMessage: %s\n\n",(index+1),element.SenderName,element.MessageContent)
			}
			
			fmt.Printf("Latency: %s\n",latency)
		}else {
			fmt.Printf("%s\n",responseMessage)
		}
	}
}

func validate(message string) bool {
	//Check for correct command: rabbl
	validCommand := regexp.MustCompile(`^rabbl$`)
	return validCommand.MatchString(message)
}
