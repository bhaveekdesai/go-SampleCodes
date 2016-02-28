package main

import ( "fmt"
	"bufio"
	"net"
	"log"
	"io/ioutil"
	"os"
	"time"
	"encoding/json"
	"strings"
	"strconv"
)

//Count for active clients & other variables
var connections int = 0
var currentClock int = 0
var initi bool = true

type Mailbox struct {
	Messages []MessagesStruct
}

type MessagesStruct struct {
    SentTime string
    SenderName string
    MessageContent string
}

func main() {
	//Open port for listening & start server
	listen,error := net.Listen("tcp",":8083")
	if error != nil {
		log.Fatal(error)
	}
	
	//Print initial messages
	fmt.Println("Welcome to rabbl server")
	
	//Accept incoming connections
	for {
		connection,error := listen.Accept()
		if error != nil {
			log.Println(error)
			continue
		}
	connections++
	fmt.Printf("Someone just joined | Active connections: %d\n",connections)
	
	//Begin local clock on first client connection
	if initi == true{
		initi = false
		go beginClock()
	}
	
	//Serve incoming connections
	go begin(connection)
	
	//Send your own new message
	go sendOwnMessage(connection)
	}
}

func beginClock() {
	//Increment every 2 seconds
	for {
		time.Sleep(2 * time.Second)
		currentClock += 2
	}
}

func begin(connection net.Conn) {
	for {
		//Read incoming message
		incomingMessage,error := bufio.NewReader(connection).ReadString('\n')
		if error != nil {
			connections--
			fmt.Printf("Someone dropped | Active connections: %d\n",connections)
			connection.Close()
			connection = nil
			return
		}
		
		//Process & build response
		responseMessage := process(string(incomingMessage))
		
		//Send Response
		connection.Write([]byte(responseMessage + "`"))
	}	
}

func process(command string) string {
	messages := ""
	if command == "rabbl\n" {
		messages = checkMailbox()
	}
	return (strconv.Itoa(currentClock)+"|Here are your messages: |"+messages)
}

func checkMailbox() string {
	file, error := ioutil.ReadFile("./Mailbox.json")
    if error != nil {
        file = []byte("|EMPTY|")
    } else {
		os.Remove("./Mailbox.json")
	}
	
	return string(file)
}

func sendOwnMessage(connection net.Conn) {
	for {
		readerHandle := bufio.NewReader(os.Stdin)
		fmt.Println("Send a new message:")
			
		fmt.Print("From: ")
		reader,_ := readerHandle.ReadString('\n')
		SenderName := strings.Split(string(reader),"\n")[0]
		
		fmt.Print("Message: ")
		reader,_ = readerHandle.ReadString('\n')
		MessageContent := strings.Split(string(reader),"\n")[0]
		
		SentTime := time.Now().Format(time.RFC1123)
		
		message := MessagesStruct{SentTime,SenderName,MessageContent}
		messages := []MessagesStruct{message}
		mailbox := Mailbox{messages}

		//Decide whether to send push notification or add to message queue (JSON file)
		if connections > 0 {
			pushNotificationMessage,_ := json.Marshal(mailbox)
			sendMessage := string(pushNotificationMessage)
			
			//Send message
			connection.Write([]byte(strconv.Itoa(currentClock)+"|New Message! |"+sendMessage + "`"))
		} else {
			//Append to message queue
			fileExists := checkMailbox()
			
			if strings.Contains(fileExists, "|EMPTY|") {
				
			} else {
				var existingMailbox Mailbox
				json.Unmarshal([]byte(fileExists), &existingMailbox)
				
				messages = append(existingMailbox.Messages, messages...)
				mailbox = Mailbox{messages}
			}
			messageQueue,_ := json.Marshal(mailbox)

			newFile,_ := os.Create("./Mailbox.json")
			newFile.Write(messageQueue)
		}
	}
}
