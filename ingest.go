package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"

	"github.com/haivision/srtgo"
)

func ingest(portno uint16, ReturnChannel chan bufferStruct, status *bool, passphrase string) {

	thisPortString := strconv.FormatUint(uint64(portno), 10)

	if utf8.RuneCountInString(passphrase) < 10 {
		fmt.Println("[INGEST] Passphrase too short, must be at least 10 characters")
		os.Exit(1) // exit non-zero
	}

	seqno := 0
	options := make(map[string]string)
	options["transtype"] = "live"
	options["passphrase"] = passphrase // passphrase must be atleast 10 char

	sck := srtgo.NewSrtSocket("0.0.0.0", portno, options)
	//   defer sck.Close()
	sck.Listen(1)
	fmt.Println("[Ingest] Listening for a new connection on port " + thisPortString + "...")

	//double-cast uint16 for clean console printout

	s, _, _ := sck.Accept() // Await a connection...
	fmt.Println("[Ingest] Connection accepted!")

	*status = true

	buf := make([]byte, 1316) // the inbound buffer
	for {
		n, _ := s.Read(buf)
		if n == 0 {
			fmt.Println("break!")
			break
		}

		// We must copy the bytes buffer into a new bytearray
		// before sending it to the consumer routine on the channel
		newbuf := make([]byte, 1316)
		copy(newbuf, buf) //newbuf out of scope next iteration
		thisBufferToSend := bufferStruct{newbuf, seqno, n}
		ReturnChannel <- thisBufferToSend
		seqno += 1
		//fmt.Printf("\rSequence no: " + strconv.Itoa(seqno) + " | Packet size: " + strconv.Itoa(n))

	} // end inner for

	fmt.Println("[INGEST] Connection was torn down, exiting gracefully...")

	sck.Close()
	s.Close()

	os.Exit(0) // Close the whole program on loss of ingest

}
