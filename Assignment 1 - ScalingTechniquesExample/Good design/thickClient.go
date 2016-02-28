package main

import ( "fmt"
	"bufio"
	"net"
	"time"
	"os"
	"strings"
	"regexp"
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
	connection,error := net.Dial("tcp",serverIP+":8083")
	if error != nil {
		log.Println(error)
	}
	
	//Print initial messages
	fmt.Println("Thick Client started")
	fmt.Println("Send your email-ID in following format: <name>_<ufid>@ufl.edu")
	fmt.Println("Example: bhaveekdesai_45616629@ufl.edu")

	for {
		//Read user entered messages
		readerHandle := bufio.NewReader(os.Stdin)
		fmt.Print("Enter email ID: ")
		outgoingMessage,_ := readerHandle.ReadString('\n')
		outgoingMessage = strings.Split(string(outgoingMessage),"\n")[0]
		
		//Validate
		valid := validate(string(outgoingMessage))
		responseMessage := "Email-ID not in desired format (<name>_<ufid>@ufl.edu). Please check."
		latency := time.Since(time.Now())
		
		//Send message to server & fetch response if validation successful
		if valid == true {
			clockStart := time.Now()
			
			fmt.Fprint(connection, outgoingMessage + "\n")
			incomingResponse,_ := bufio.NewReader(connection).ReadString('\n')
			
			latency = time.Since(clockStart)
			responseMessage = strings.Split(incomingResponse,"\n")[0]
		}
		
		//Print response
		fmt.Printf("Response: \n-----%s\n-----Latency: %s\n",responseMessage, latency)
		
	}
}

func validate(message string) bool {
	//Check for valid email in the format (eg): bhaveekdesai_45616629@ufl.edu 
	validID := regexp.MustCompile(`^[a-z]+\_[0-9]+@ufl.edu$`)
	return validID.MatchString(message)
}
