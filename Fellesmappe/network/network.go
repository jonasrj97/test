// Takk til Joachim Holwech for god dokumentasjon på broadcast- & listen-funksjonene.

package network

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"

	c "../config"
)

// TODO: Lagre alt av ID, IsMaster etc inn i elevatorState
// endre Heartbeat til å kun sende ut elevatorstate
// endre activeElevators til å være en struct
// fikse kommunikasjonen med nettverk på esm-delen
// faktisk teste

// Globale variabler

// ID til noden. -1 om noden ikke er konfigurert.
var ID int = -1

// IsMaster - om noden har rollen som master eller ikke.
var IsMaster bool

// ActiveElevators - et map som inneholder alle aktive noder (at disse har sendt ut hearbeat innen det siste sekundet)
var ActiveElevators map[int]bool

//setup -- om noden er konfigurert ved oppstart
var setup bool = false

// broadcastaddress er adressen som pakkene pushes til. []byte{127,0,0,1} - local , {255,255,255,255} - broadcast
var broadcastaddress []byte

var elevatorState c.ElevatorState

// Init initialiserer portene og starter goroutinene for sending, mottaking og scanning.
func Init(readPort int, writePort int, broadcastAddress []byte, StateChan chan c.ElevatorState) (<-chan c.Packet, chan<- c.Packet) {
	broadcastaddress = broadcastAddress
	recieve := make(chan c.Packet, 10)
	send := make(chan c.Packet, 10)
	internal := make(chan c.Packet, 10)
	start := make(chan int)
	toESM := make(chan c.ElevatorState)
	go listen(recieve, internal, readPort)
	go broadcast(send, writePort)
	go heartbeat(send, start)
	go scan(internal, start, toESM)
	return recieve, send
}

func getSetState(comms chan c.ElevatorState, internal chan c.ElevatorState) {
	r := <-internal // motta ID og IsMaster fra scan
	comms <- r      // send ID og IsMaster til esm
	for {
		r := <-comms // motta state fra esm
		elevatorState = r
	}

}

func timeOut(signal chan<- int) {
	time.Sleep(time.Second)
	signal <- 1
	return
}

// heartbeat viser nodens tilstedeværelse ved å regelmessig sende ut heartbeats
func heartbeat(outgoing chan<- c.Packet, signal <-chan int) {
	var heartbeat c.Packet
	heartbeat.IDto = 0     // 0 - til alle
	heartbeat.Reply = "hb" // heartbeat
	<-signal               // vent på startsignal
	for {
		heartbeat.IDfrom = ID
		heartbeat.IsMaster = IsMaster
		heartbeat.State = elevatorState
		outgoing <- heartbeat
		time.Sleep(time.Millisecond * 100) // 10Hz
	}
}

// Scan skanner nettverket for å oppdage andre noders tilstedevære, og om disse er master eller ikke. Den tildeler også noden en ID ved oppstart.
// Scan holder også et map oppdatert med alle aktive noder på nettverket.
func scan(incoming <-chan c.Packet, startSignal chan<- int, toESM chan<- c.ElevatorState) {
	for {
		timeoutSignal := make(chan int)
		go timeOut(timeoutSignal)
		activeElevators := make(map[int]bool)
	loop:
		for {
			select {
			case r := <-incoming:
				if r.IDto == 0 && r.Reply == "hb" {
					activeElevators[r.IDfrom] = r.IsMaster
				}
			case <-timeoutSignal:
				break loop
			}
		}
		if !setup && len(activeElevators) == 0 {
			var r c.ElevatorState
			r.ID = 1
			r.IsMaster = true
			toESM <- r
			ID = 1          // overflødig
			IsMaster = true // overflødig
			startSignal <- 1
			setup = true
		}
		highestID := 0
		for k := range activeElevators {
			if k > highestID {
				highestID = k
			}
		}
		if !setup && len(activeElevators) > 0 {
			var r c.ElevatorState
			r.ID = highestID + 1

			ID = highestID + 1 // overflødig
			startSignal <- 1
			setup = true
		}
		ActiveElevators = activeElevators
	}
}

// listen - lytter til nettverket og mottar pakker.
func listen(recieve chan<- c.Packet, internal chan<- c.Packet, port int) {
	connection, err := net.ListenUDP("udp", &net.UDPAddr{IP: broadcastaddress, Port: port, Zone: ""})
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	for {
		var message c.Packet
		inputBytes := make([]byte, 4096)
		length, _, err := connection.ReadFromUDP(inputBytes)
		if err != nil {
			log.Fatal(err)
		}
		buffer := bytes.NewBuffer(inputBytes[:length])
		decoder := gob.NewDecoder(buffer)
		decoder.Decode(&message)
		if message.IDto == 0 || message.IDto == ID {
			recieve <- message
			internal <- message
		}
	}
}

// broadcast - sender pakker ut på nettet.
func broadcast(send <-chan c.Packet, port int) {
	Conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: broadcastaddress, Port: port, Zone: ""})
	if err != nil {
		log.Fatal(err)
	}
	defer Conn.Close()
	for {
		message := <-send
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		encoder.Encode(message)
		Conn.Write(buffer.Bytes())
		buffer.Reset()
	}
}
