package main

import ( "fmt"
	"bufio"
	"net"
	"time"
	"os"
	"strings"
	"log"
	"strconv"
	"math"
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
	connection,error := net.Dial("tcp",serverIP+":8084")
	if error != nil {
		log.Println(error)
	}
	
	//Print initial messages
	fmt.Println("Hi Bob!")

	negotiated := false
	for negotiated == false {
		//Read user entered messages
		readerHandle := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your private key (integer): ")
		privateKey,_ := readerHandle.ReadString('\n')
		privateKey = strings.Split(string(privateKey),"\n")[0]
		
		//Calculate public key
		privateKeyFloat,_ := strconv.ParseFloat(privateKey, 64)
		publicKey := calculate_public_key(privateKeyFloat, 0.0)
		
		//Send own public key
		fmt.Printf("Public key sent to Alice: %s\n",publicKey)
		fmt.Fprint(connection, publicKey + "\n")
		
		//Receive Alice's public key
		incomingResponse,_ := bufio.NewReader(connection).ReadString('\n')
		receivedKey := strings.Split(incomingResponse,"\n")[0]
		receivedKeyFloat,_ := strconv.ParseFloat(receivedKey, 64)
	
		fmt.Printf("Received Alice's Public key: %s\n",receivedKey)
		
		//Calculate shared secret key
		sharedSecretKey := calculate_public_key(privateKeyFloat, receivedKeyFloat)
		
		//Print shared secret key
		fmt.Printf("Shared secret key: %s\n",sharedSecretKey)
		
		negotiated = true
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
