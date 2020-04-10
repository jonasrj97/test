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
		
	elevatorChans := esm.ESMChannels{
		ButtonChan: make(chan io.ButtonEvent),
		FloorChan:  make(chan int),
		StateChan:  make(chan esm.ElevatorState), //for bruk mellom network og esm (orders?)
		KillChan:   make(chan bool),
	}

	elevatorChans := c.ElevatorChannels{
		FloorChan:  			make(chan int),
		StateChan:  			make(chan c.ElevatorState), //for bruk mellom network og esm (orders?)
		ButtonChan:				make(chan io.ButtonEvent),
		ReceivedFloorHallChan:	make(chan io.ButtonEvent),
	}

	orderDistributorChans := c.OrderDistributionChannels{
		// ActiveElevatorsChan:	make(chan map[int]esm.ElevatorState),
		NewOrderChan:			make(chan int),
		LightMatrixChan:		make(chan []bool),
		FloorHallMatrixChan:	make(chan []c.MatrixElement),
	}
	
	//Initializing
	io.Init("localhost:15650", numFloors)
	// recieve, send := network.Init(10001, 10001, []byte{255, 255, 255, 255}) // recieve og send er channels for hhv. mottaking og sending

	//goroutines
	go esm.RunElevator(elevatorChans, orderDistributorChans.NewOrderChan, orderDistributorChans.LightMatrixChan, killSig)
	go od.DistributeOrders(orderDistributorChans, elevatorChans.ReceivedFloorHallChan, numFloors)
	go io.PollButtons(elevatorChans.ButtonChan)
	go io.PollFloorSensor(elevatorChans.FloorChan)

	select {}
}
