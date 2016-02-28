package main

import ( "fmt"
	"strings"
	"bufio"
	"net"
	"regexp"
	"log"
)

//Count for active clients
var connections int= 0

func main() {
	//Open port for listening & start server
	listen,error := net.Listen("tcp",":8082")
	if error != nil {
		log.Fatal(error)
	}
	
	//Print initial messages
	fmt.Println("Thick Server started")
	fmt.Println("Send me your email-ID in following format: <name>_<ufid>@ufl.edu")
	fmt.Println("Example: bhaveekdesai_45616629@ufl.edu")

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
		incomingMessage = strings.Split(string(incomingMessage),"\n")[0]
		
		//Validate
		valid := validate(string(incomingMessage))
		
		//Process & build response
		responseMessage := "Email-ID not in desired format (<name>_<ufid>@ufl.edu). Please check."
		if valid {
			responseMessage = process(string(incomingMessage))
		}
		
		//Send Response
		connection.Write([]byte(responseMessage + "\n"))
	}
}

func validate(message string) bool {
	//Check for valid email in the format (eg): bhaveekdesai_45616629@ufl.edu 
	validID := regexp.MustCompile(`^[a-z]+\_[0-9]+@ufl.edu$`)
	return validID.MatchString(message)
}

func process(message string) string {
	components := strings.Split(message,"_")
	name := components[0]
	ufid := strings.Split(components[1],"@")[0]
	
	return ("Hello "+name+". I believe your UFID is: "+ufid)
}
