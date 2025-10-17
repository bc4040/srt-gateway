package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"

	"github.com/haivision/srtgo"
)

func sender(portno uint16, DataChannel chan bufferStruct, status *bool, passphrase string) {
	thisPortString := strconv.FormatUint(uint64(portno), 10)
	fmt.Println("Sender @ " + thisPortString + ": Starting up")

	if utf8.RuneCountInString(passphrase) < 10 {
		fmt.Println("[SENDER] Passphrase too short, must be at least 10 characters")
		os.Exit(1) // exit non-zero
	}

	//  socket setup
	options := make(map[string]string)
	options["transtype"] = "live"
	options["passphrase"] = passphrase

	for { // outer for
		sck := srtgo.NewSrtSocket("0.0.0.0", portno, options)
		//Mode: Listener
		defer sck.Close()
		sck.Listen(1)
		fmt.Println("[SENDER] " + thisPortString + " Waiting for connection in SRT LISTENER mode")
		s, peeraddr, _ := sck.Accept() //socket, peeraddr, err := sck.Accept()
		fmt.Println("[SENDER] " + thisPortString + " Got connection from client: " + peeraddr + ", starting to send data...")
		*status = true

		lastseqno := 0
		outOfOrderCount := 0
		for { // inner for
			/// DATA TRANSMISSION OUTBOUND
			thisBufferStruct := <-DataChannel

			if thisBufferStruct.seqno >= 1 { // catch odd condition where seqno could be out of order? (seqno is an internal value)
				if lastseqno != thisBufferStruct.seqno-1 {
					fmt.Println("Sender @ " + thisPortString + ": Seqno out of order!")
					fmt.Println("[SENDER] " + "Seq number out-of-order was detected on port " + thisPortString)
					fmt.Println(thisBufferStruct.seqno)
					outOfOrderCount += 1
				}
			}
			thisdata := Truncate(thisBufferStruct.buffer, thisBufferStruct.length)

			// If file write is enabled, write this buffer to file:
			//n, _ := thisfile.Write(thisdata)

			n, _ := s.Write(thisdata) // WRITE INTO THE SOCKET
			if n == 0 {
				// Something went wrong or the channel closed gracefully?
				fmt.Println("[SENDER] " + "Sender socket closed (no data) on port " + thisPortString)
				break // break inner and await next connection
			}
			lastseqno = thisBufferStruct.seqno

		} // end inner for

		fmt.Println("[SENDER] " + "Sender socket closed on port " + thisPortString)
		*status = false
		sck.Close()

	} // end outer for

}
