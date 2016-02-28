package main

import ( "fmt"
	"strings"
	"bufio"
	"net"
	"log"
	"strconv"
	"math"
	"os"
)

//Count for active clients
var connections int= 0
	
func main() {
	//Open port for listening & start server
	listen,error := net.Listen("tcp",":8084")
	if error != nil {
		log.Fatal(error)
	}
	
	//Print initial messages
	fmt.Println("Hi Alice!")
	
	//Accept incoming connections
	for {
		connection,error := listen.Accept()
		if error != nil {
			log.Println(error)
			continue
		}
	connections++
	//fmt.Printf("A client connected | Active connections: %d\n",connections)
	go begin(connection)
	}
}

func begin(connection net.Conn) {
	
	for {
		//Receive Bob's public key
		incomingResponse,error := bufio.NewReader(connection).ReadString('\n')
		receivedKey := strings.Split(incomingResponse,"\n")[0]
		receivedKeyFloat,_ := strconv.ParseFloat(receivedKey, 64)
		
		if error != nil {
			connections--
			//fmt.Printf("A client dropped | Active connections: %d\n",connections)
			connection.Close()
			connection = nil
			return
		}
		fmt.Printf("Received Bob's Public key: %s\n",receivedKey)
		
		//Read user entered messages
		readerHandle := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your private key (integer): ")
		privateKey,_ := readerHandle.ReadString('\n')
		privateKey = strings.Split(string(privateKey),"\n")[0]
		
		//Calculate public key
		privateKeyFloat,_ := strconv.ParseFloat(privateKey, 64)
		publicKey := calculate_public_key(privateKeyFloat, 0.0)
		
		//Calculate shared secret key
		sharedSecretKey := calculate_public_key(privateKeyFloat, receivedKeyFloat)
		
		//Send own public key
		fmt.Printf("Public key sent to Bob: %s\n",publicKey)
		connection.Write([]byte(publicKey + "\n"))
		
		//Print shared secret key
		fmt.Printf("Shared secret key: %s\n",sharedSecretKey)
	}
}

func calculate_public_key(privateKey float64, keyReceived float64) string {
	//Calculate public key
	public_base_g := 5.0
	public_modulus_p := 23.0
	
	if (keyReceived != 0.0) {
		public_base_g = keyReceived
	}
		
	public_key := math.Mod(math.Pow(public_base_g, privateKey), public_modulus_p)
	public_key_string := strconv.FormatFloat(public_key,'g',-1,64)
	return public_key_string
}
