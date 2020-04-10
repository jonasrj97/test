package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func count(countChannel chan<- int) {
	for {
		countChannel <- 1
		time.Sleep(time.Millisecond * 500)
	}
}

func netPoll() {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 10001, Zone: ""})
	defer ServerConn.Close()
	timeoutDuration := 2 * time.Second
	ServerConn.SetReadDeadline(time.Now().Add(timeoutDuration))
	buf := make([]byte, 1024)
	n, _, _ := ServerConn.ReadFromUDP(buf)
	if len(string(string(buf[0:n]))) > 0 {
		fmt.Println("Connection established! ")
		isListening = true
		return
	}
	fmt.Print("No connection found. ")
	isListening = false
	return
}

func netPush(toBroadcast <-chan string) {
	for {
		select {
		case x := <-toBroadcast:
			Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 10001, Zone: ""})
			defer Conn.Close()
			Conn.Write([]byte(x))
		}
	}
}

func netListen(recieved chan<- int, timedOut chan<- int) {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 10001, Zone: ""})
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		timeoutDuration := 2 * time.Second
		ServerConn.SetReadDeadline(time.Now().Add(timeoutDuration))
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Print("Connection timed out. ")
			timedOut <- 1
			return
		}
		fmt.Print("Received ", string(buf[0:n]), " from ", addr, " ")
		tempString := string(buf[0:n])
		tempNumber, _ := strconv.Atoi(tempString)
		recieved <- 1
		backup = tempNumber
	}
}

var backup int = 0
var isListening bool

func main() {
	var teller = 0
	countChannel := make(chan int)
	toBroadcast := make(chan string)
	recieved := make(chan int)
	timedOut := make(chan int)

	fmt.Println("Listening..")
	netPoll()

	for {
		if !isListening {
			fmt.Println("Counting initiated.")
			go count(countChannel)
			go netPush(toBroadcast)
			for {
				select {
				case <-countChannel:
					teller++
					tempNumber := strconv.Itoa(teller)
					toBroadcast <- tempNumber
					fmt.Println(teller)
				}
			}
		} else {
			go netListen(recieved, timedOut)
			for {
				select {
				case <-recieved:
					fmt.Print("Locally saved value: ")
					fmt.Println(backup)
				case <-timedOut:
					fmt.Println("Resuming counting..")
					isListening = false
					teller = backup
				}
				if isListening == false {
					break
				}
			}
		}
	}
}
