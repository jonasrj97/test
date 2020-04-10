package main

import (
	"fmt"

	"./elevio"
)

func main() {
	fmt.Println("Heisann Hoppsann!")
	//Global variables
	const numFloors = 4
	var currentFloor, nextFloor int
	var dir elevio.MotorDirection = elevio.MD_Stop
	//var hallMatrix [2][numFloors]int

	//channels
	buttonChan := make(chan elevio.ButtonEvent)
	floorChan := make(chan int)
	//drv_obstr := make(chan bool) trenger ikke denne siden vi ikke simulerer med obstruksjoner
	//heisStopp := make(chan bool) denne valgte vi også ikke å inkludere

	//goroutines
	go elevio.PollButtons(buttonChan)
	go elevio.PollFloorSensor(floorChan)
	//go elevio.PollObstructionSwitch(drv_obstr)
	//go elevio.PollStopButton(heisStopp)

	//fsm
	currentState := "Initializing"

	//for i := 0; i < 3; i++ {
	//	for j := 0; j < 4; j++ {
	//		elevio.SetButtonLamp(elevio.ButtonType(i), j, false)
	//	}
	//}

	for {
		switch currentState {
		case "Initializing":
			elevio.Init("localhost:15657", numFloors)
			fmt.Printf("Kjører ned til nærmeste etasje..")
			elevio.SetMotorDirection(elevio.MD_Down)
			<-floorChan
			elevio.SetMotorDirection(elevio.MD_Stop)
			currentState = "Idle"
			nextFloor = -1
			fmt.Printf("%s", currentState)

		case "Idle":
			if nextFloor == -1 {
				order := <-buttonChan
				nextFloor = order.Floor
				fmt.Printf("Bestilling til etasje ")
				fmt.Printf("%+v", order.Floor+1)
				fmt.Println(" mottatt.")

				if nextFloor != currentFloor {
					currentState = "Running"
					elevio.SetButtonLamp(order.Button, order.Floor, true)
				} else {
					currentState = "Doors Open"
				}
			} else {
				elevio.SetDoorOpenLamp(false)
				currentState = "Running"
				fmt.Printf("%s", currentState)
			}

		case "Running":
			if nextFloor < currentFloor {
				elevio.SetDoorOpenLamp(false)
				fmt.Println("Døren lukkes.")
				dir = elevio.MD_Down
			} else if nextFloor > currentFloor {
				elevio.SetDoorOpenLamp(false)
				fmt.Println("Døren lukkes.")
				dir = elevio.MD_Up
			} else {
				dir = elevio.MD_Stop
				currentState = "Open Doors"
			}
			elevio.SetMotorDirection(dir)

			currentFloor = <-floorChan
			elevio.SetFloorIndicator(currentFloor)
			if currentFloor == nextFloor {
				currentState = "Open Doors"
				fmt.Printf("%s", currentState)
			}

		case "Open Doors":
			currentFloor = nextFloor
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetButtonLamp(elevio.ButtonType(2), currentFloor, false)
			if dir == elevio.MD_Down {
				elevio.SetButtonLamp(elevio.ButtonType(1), currentFloor, false)
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(0), currentFloor, false)
			}
			elevio.SetDoorOpenLamp(true)
			fmt.Println("*Pling! Døren spretter opp*")
			nextFloor = -1
			currentState = "Idle"
			fmt.Printf("%s", currentState)
		}
	}
}
