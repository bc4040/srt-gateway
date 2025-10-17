package main

import (
	"fmt"
	"os"
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

type DataChannel struct {
	bufferStruct chan bufferStruct
	channelOpen  bool
	peerAddr     string
}

func main() {

	// Load environment variables and check they exist
	ingestPassphrase := os.Getenv("PASSPHRASE_IN")
	sender1Passphrase := os.Getenv("PASSPHRASE_OUT1")
	sender2Passphrase := os.Getenv("PASSPHRASE_OUT2")

	var sender1enabled bool
	var sender2enabled bool

	var ingestport16 uint16
	var sender1port16 uint16
	var sender2port16 uint16

	ingestport16 = 9800
	sender1port16 = 9801
	sender2port16 = 9802

	if ingestPassphrase == "" { // if no passphrase defined for ingest, exit
		fmt.Println("Error - No ingest passphrase defined.  Define environment var PASSPHRASE_IN=")
		os.Exit(1)
	}

	if sender1Passphrase == "" { // if no passphrase defined for sender 1, exit
		fmt.Println("Error - No sender 1 passphrase defined.  Define environment var PASSPHRASE_OUT1=")
		os.Exit(1)
	} else {
		sender1enabled = true
	}

	if sender2Passphrase == "" { // if no passphrase defined for sender 2, assume it's not in use
		fmt.Println("Warning - No sender 2 passphrase defined.  Define environment var PASSPHRASE_OUT2=")
		fmt.Println("Assuming sender 2 is not in use.")
		sender2enabled = false
	} else {
		sender2enabled = true
	}

	// Make status bools to track if socket is open & streaming or not
	var IngestOpen bool

	//Make list of active data channels
	var ActiveDataChannels []DataChannel

	//Make internal channel
	IngestChannel := make(chan bufferStruct, 100) // Inbound

	//Make flexible channels
	DataChannel1 := DataChannel{
		bufferStruct: make(chan bufferStruct, 10000),
		channelOpen:  false,
		peerAddr:     "0.0.0.0",
	}

	DataChannel2 := DataChannel{
		bufferStruct: make(chan bufferStruct, 10000),
		channelOpen:  false,
		peerAddr:     "0.0.0.0",
	}

	ActiveDataChannels = append(ActiveDataChannels, DataChannel1)
	ActiveDataChannels = append(ActiveDataChannels, DataChannel2)

	// Call the ingester and await some data
	go ingest(ingestport16, IngestChannel, &IngestOpen, ingestPassphrase)

	fmt.Println("Awaiting ingest connection...")
	for !IngestOpen {
		time.Sleep(10 * time.Millisecond) // Avoid busy-waiting
	}

	// Call the sender and await a connection
	if sender1enabled == true {
		go sender(sender1port16, &ActiveDataChannels[0], sender1Passphrase)
	}

	if sender2enabled == true {
		go sender(sender2port16, &ActiveDataChannels[1], sender2Passphrase)
	}

	for { // multiplex data as we get it from the ingester...
		//fmt.Printf("-") // tick

		thisBufferStruct := <-IngestChannel // as we get data from ingest
		//fmt.Println(thisBufferGlob.seqno)

		// feed data into the active data channels...
		for _, DataChannel := range ActiveDataChannels {
			if DataChannel.channelOpen {
				DataChannel.bufferStruct <- thisBufferStruct
			}
		}

	}
}
