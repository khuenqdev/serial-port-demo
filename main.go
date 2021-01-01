package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/tarm/serial"
)

func main() {
	// Automatically detect open serial port
	s := checkPorts()

	if nil == s {
		panic("No serial port detected!")
	}

	// Send test data to the connected device
	data := "Test data"

	fmt.Println("Writing data")
	_, err := s.Write([]byte(data))

	if nil != err {
		panic("Failed!")
	}

	// Read device response
	fmt.Println("Reading response")

	buf := make([]byte, 2048)
	readLen, err := s.Read(buf)

	if nil != err {
		panic("Cannot read data from terminal")
	}

	fmt.Println("Read length", readLen)
	buf = bytes.Trim(buf, "\x00")

	// Show received response
	fmt.Printf("Received data: %s", string(buf))

	// Close the connection when done
	err = s.Close()

	if nil != err {
		panic("Cannot close terminal connection")
	}

	fmt.Println("Terminal connection closed")
}

func checkPorts() *serial.Port {
	portList := getPortList()

	for _, p := range portList {
		connection := &serial.Config{
			Name:        p,
			Baud:        9600,
		}
		s, err := serial.OpenPort(connection)

		if nil != err {
			continue
		}

		fmt.Println("Opened port", p)
		return s
	}

	return nil
}

func getPortList() []string {
	var portNames []string
	windowsPorts := []string{
		"COM1",
		"COM2",
		"COM3",
		"COM4",
		"COM5",
		"COM6",
		"COM7",
		"COM8",
		"COM9",
		"COM10",
	}

	portNames = append(portNames, windowsPorts...)

	files, err := ioutil.ReadDir("/dev")

	if err != nil {
		return portNames
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "ttyACM") || strings.Contains(file.Name(), "tty.usbmodem") {
			fmt.Println(file.Name())
			portNames = append(portNames, "/dev/" + file.Name())
		}
	}

	return portNames
}
