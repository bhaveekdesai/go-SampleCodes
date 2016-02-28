package main

import ( "fmt"
	"bufio"
	"net"
	"log"
	"io/ioutil"
)

//Count for active clients
var connections int= 0
	
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
	go begin(connection)
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
		fmt.Printf("Command received: %s", string(incomingMessage))
		
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
	return ("Here are your messages: |"+messages)
}

func checkMailbox() string {
	file, error := ioutil.ReadFile("./Mailbox.json")
    if error != nil {
        fmt.Printf("File error: %v\n", error)
    }

	return string(file)
}
