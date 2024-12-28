package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func Truncate(s []byte, to int) []byte {
	return s[:to]
}

func filecreate(filename string) *os.File {
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

	_ingestport := os.Getenv("PORT_IN")
	_sender1port := os.Getenv("PORT_OUT1")
	_sender2port := os.Getenv("PORT_OUT2")

	var ingestport16 uint16
	var sender1port16 uint16
	var sender2port16 uint16

	if _ingestport == "" {
		fmt.Println("Error - No ingest port defined.  Define environment var PORT_IN=")
		os.Exit(1)
	} else {
		ingestport, _ := strconv.ParseUint(_ingestport, 10, 16)
		ingestport16 = uint16(ingestport)
	}

	if _sender1port == "" {
		fmt.Println("Error - No sender 1 port defined.  Define environment var PORT_OUT1=")
		os.Exit(1)
	} else {
		sender1port, _ := strconv.ParseUint(_sender1port, 10, 16)
		sender1port16 = uint16(sender1port)

	}

	if _sender2port == "" {
		fmt.Println("WARNING - No sender 2 port defined.  Define environment var PORT_OUT2=")
	} else {
		sender2port, _ := strconv.ParseUint(_sender2port, 10, 16)
		sender2port16 = uint16(sender2port)

	}

	if ingestPassphrase == "" {
		fmt.Println("Error - No ingest passphrase defined.  Define environment var PASSPHRASE_IN=")
		os.Exit(1)
	}

	if sender1Passphrase == "" {
		fmt.Println("Error - No sender 1 passphrase defined.  Define environment var PASSPHRASE_OUT1=")
		os.Exit(1)
	}

	if _sender2port != "" && sender2Passphrase == "" {
		fmt.Println("Error - No sender 2 passphrase defined.  Define environment var PASSPHRASE_OUT2=")
		os.Exit(1)
	}

	// Make status bools to track if socket is open & streaming or not
	var IngestOpen bool
	var Channel1Open bool
	var Channel2Open bool

	//Make channels
	IngestChannel := make(chan bufferStruct, 100)  // Inbound
	DataChannel := make(chan bufferStruct, 10000)  // Outbound 1
	DataChannel2 := make(chan bufferStruct, 10000) // Outbound 1

	// Call the ingester and await some data
	go ingest(ingestport16, IngestChannel, &IngestOpen, ingestPassphrase)

	fmt.Println("Awaiting ingest connection...")
	for !IngestOpen {
		time.Sleep(10 * time.Millisecond) // Avoid busy-waiting
	}

	// Call the sender and await a connection
	go sender(sender1port16, DataChannel, &Channel1Open, sender1Passphrase)

	if _sender2port != "" {
		go sender(sender2port16, DataChannel2, &Channel2Open, sender2Passphrase)
	}

	for { // multiplex data as we get it from the ingester...
		//fmt.Printf("-") // tick

		thisBufferStruct := <-IngestChannel
		//fmt.Println(thisBufferGlob.seqno)

		if Channel1Open {
			DataChannel <- thisBufferStruct
		}

		if Channel2Open {
			DataChannel2 <- thisBufferStruct
		}
	}

	os.Exit(0) // Main loop died
}
