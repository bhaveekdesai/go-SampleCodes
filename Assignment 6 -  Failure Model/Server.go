package main

import ( "fmt"
	"strings"
	"bufio"
	"net"
	"log"
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
	fmt.Println("Server started")
	
	//Accept incoming connections
	for {
		connection,error := listen.Accept()
		if error != nil {
			log.Println(error)
			continue
		}
	connections++
	fmt.Printf("A client connected | Active connections: %d\n",connections)
	go begin(connection)
	}
}

func begin(connection net.Conn) {
	for {
		//Read incoming message
		incomingMessage,error := bufio.NewReader(connection).ReadString('\n')
		if error != nil {
			connections--
			fmt.Printf("A client dropped | Active connections: %d\n",connections)
			connection.Close()
			connection = nil
			return
		}
		fmt.Printf("Message received: %s", string(incomingMessage))
		
		//Process & build response
		responseMessage := process(string(incomingMessage))
		
		//Send Response
		connection.Write([]byte(responseMessage + "\n"))
	}	
}

func process(message string) string {
	components := strings.Split(message,"_")
	name := components[0]
	ufid := strings.Split(components[1],"@")[0]
	
	return ("Hello "+name+". I believe your UFID is: "+ufid)
}
