package main

import (
	"./backup"
)

//ElevatorState aah

func main() {
	recovered := backup.Recover()
	print(recovered.ID)
	/* Network-testing
	recieve, _ := network.Init(10001, 10001, []byte{255, 255, 255, 255})
	go printInformation()
	// go count(send) // push packets
	for {
		select {
		case p := <-recieve:
			if p.Reply != ("hb") {
				fmt.Printf("%+v\n", p)
			}
			//s := string(p.Reply)
			//println(s)
		}
	}
	*/

	/*
	   func printInformation() {
	   	writer := uilive.New()
	   	writer.Start()
	   	println("------Elevator Info-------")
	   	for {
	   		fmt.Fprintf(writer, "ID: %d \nisMaster: %t \nActive elevators (ID:IsMaster): %v\n", network.ID, network.IsMaster, network.ActiveElevators)
	   		time.Sleep(time.Millisecond * 1000)
	   	}
	   }

	   func count(countChannel chan<- network.Packet) {
	   	var melding network.Packet
	   	melding.IDto = 0
	   	melding.IDfrom = network.ID
	   	melding.Reply = "Hallo på do"
	   	for {
	   		//countChannel <- melding
	   		time.Sleep(time.Millisecond * 100)
	   	}
	   }
	*/
}
