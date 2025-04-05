package main

import (
	"fmt"
	"os"
	"time"
)

func Truncate(s []byte, to int) []byte {
	return s[:to]
}

func filecreate() *os.File {

	currentTime := time.Now()
	filename := currentTime.Format("2006-01-02_15-04-05")
	Myfile, err := os.OpenFile(filename+".ts", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Unable to open file")
	} else {
		fmt.Println("File allocated: " + filename + ".ts")
		return Myfile //return pointer to os.File
	}
	return nil
}

type bufferStruct struct {
	buffer []byte
	seqno  int
	length int
}

func main() {

	// Load environment variables and check they exist
	ingestPassphrase := os.Getenv("PASSPHRASE_IN")
	sender1Passphrase := os.Getenv("PASSPHRASE_OUT1")
	sender2Passphrase := os.Getenv("PASSPHRASE_OUT2")
	record := os.Getenv("record")

	var sender1enabled bool
	var sender2enabled bool
	var recordenabled bool

	var ingestport16 uint16
	var sender1port16 uint16
	var sender2port16 uint16

	ingestport16 = 9800
	sender1port16 = 9801
	sender2port16 = 9802

	// CHECK IF INGEST CHANNEL HAS A PASSPHRASE DEFINED
	if ingestPassphrase == "" { // if no passphrase defined for ingest...
		fmt.Println("Error - No ingest passphrase defined.  Define environment var PASSPHRASE_IN=")
		os.Exit(1)
	}

	// CHECK IF SENDER 1 IS ENABLED, IF NOT - EXIT
	if sender1Passphrase == "" { // if no passphrase defined for sender 1...
		fmt.Println("Error - No sender 1 passphrase defined.  Define environment var PASSPHRASE_OUT1=")
		fmt.Println("We cannot continue")
		os.Exit(1)
	} else {
		sender1enabled = true
	}

	// CHECK IF SENDER 2 IS ENABLED

	if sender2Passphrase == "" { // if no passphrase defined for sender 2, assume it's not in use
		fmt.Println("Warning - No sender 2 passphrase defined.  Define environment var PASSPHRASE_OUT2=")
		fmt.Println("Assuming sender 2 is not in use.")
		sender2enabled = false
	} else {
		sender2enabled = true
	}

	if record == "true" {
		recordenabled = true
	} else {
		recordenabled = false
	}

	if recordenabled {
		filecreate()
	}

	// Make status bools to track if socket is open & streaming or not
	var IngestOpen bool
	var Channel1Open bool
	var Channel2Open bool

	//Make channels
	IngestChannel := make(chan bufferStruct, 100)  // Ingest ("inbound")
	DataChannel := make(chan bufferStruct, 10000)  // Sender 1 ("outbound")
	DataChannel2 := make(chan bufferStruct, 10000) // Sender 2 ("outbound")

	// Call the ingester and await some data
	go ingest(ingestport16, IngestChannel, &IngestOpen, ingestPassphrase)

	fmt.Println("Awaiting ingest connection...")
	for !IngestOpen {
		time.Sleep(10 * time.Millisecond) // Avoid busy-waiting
	}

	// Call the sender and await a connection
	if sender1enabled == true {
		go sender(sender1port16, DataChannel, &Channel1Open, sender1Passphrase)
	}

	if sender2enabled == true {
		go sender(sender2port16, DataChannel2, &Channel2Open, sender2Passphrase)
	}

	for { // multiplex data as we get it from the ingester...
		thisBufferStruct := <-IngestChannel

		if Channel1Open {
			DataChannel <- thisBufferStruct
		}

		if Channel2Open {
			DataChannel2 <- thisBufferStruct
		}

		// If file write is enabled, write this buffer to file:
		//n, _ := thisfile.Write(thisdata)

	}

}
