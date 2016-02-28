package main

import ( "fmt"
	"bufio"
	"os"
	"crypto/aes"
	"crypto/cipher"
)
	
func main() {

	//Initialize key
	key := []byte("91682456317480253987523640978432")
	initializationVector := []byte(key)[:aes.BlockSize]
	
	//Print initial messages
	fmt.Println("Hi! Enter your message and get AES encrypted cipher!")
	
	//Start user interaction
	for {
		//Read user entered messages
		readerHandle := bufio.NewReader(os.Stdin)
		fmt.Print("Message: ")
		inputText,_ := readerHandle.ReadString('\n')
		message := []byte(inputText)
		
		cipherText := make([]byte, len(message))
		encrypt(message, key, initializationVector, cipherText)
		fmt.Printf("Encrypted Message: %s\n", cipherText)
		
	}
}

func encrypt (message []byte, key []byte, initializationVector []byte, cipherText []byte) {
	
	//Encrypt message
	aesEncrypter, _ := aes.NewCipher(key)
	(cipher.NewCFBEncrypter(aesEncrypter, initializationVector)).XORKeyStream(cipherText, message)
}
