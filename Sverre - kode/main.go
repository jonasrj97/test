package main

import (
	"fmt"
	io "./elevio"
	od "./orderDistribution"
	c "./config"
	"./esm"
	"os/signal"
	"os"
	"syscall"
)

const numFloors = 4


func main() {
	fmt.Println("cheeky fokkah")

	//channels
	killSig := make(chan os.Signal,1)
	signal.Notify(killSig, syscall.SIGINT)

	elevatorChans := c.ElevatorChannels{
		FloorChan:  			make(chan int),
		StateChan:  			make(chan c.ElevatorState), //for bruk mellom network og esm (orders?)
		ButtonChan:				make(chan io.ButtonEvent),
		ReceivedFloorHallChan:	make(chan io.ButtonEvent),
	}

	orderDistributorChans := c.OrderDistributionChannels{
		ActiveElevatorsChan:	make(chan map[int]c.ElevatorState),
		NewOrderChan:			make(chan io.ButtonEvent),
		IsMasterChan:			make(chan bool),
		LightMatrixChan:		make(chan []bool),
		FloorHallMatrixChan:	make(chan []c.MatrixElement),
	}
	
	//Initializing
	io.Init("localhost:15650", numFloors) 

	//goroutines
	go esm.RunElevator(elevatorChans, orderDistributorChans.NewOrderChan, orderDistributorChans.LightMatrixChan)
	go od.DistributeOrders(orderDistributorChans, elevatorChans.ReceivedFloorHallChan, numFloors)
	go io.PollButtons(elevatorChans.ButtonChan)
	go io.PollFloorSensor(elevatorChans.FloorChan)


	//Avslutter main hvis en melding sendes over KillChan
	select{
	case <- killSig:
		fmt.Printf("\n\nelevator is kill. nooooooo\n")
		os.Exit(1)
	}
}
