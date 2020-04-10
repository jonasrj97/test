package esm

import (
	io "../elevio"
	c "../config"
	"fmt"
)


func RunElevator(chans c.ElevatorChannels, NewOrderChan <-chan io.ButtonEvent, LightMatrixChan chan []bool) {
	currentOrders := make([]io.ButtonEvent,0,10) //cab halls
	doorLightChan := make(chan int)
	elevator := c.ElevatorState{
		ID:				1,			//må hentes fra nettverksmodulen
		CurrentState:	0,
		CurrentFloor:	-1,
		NextFloor:		-1,
		Dir:			io.MD_Stop,
	}
	
	go toggleFloorButtonLights(LightMatrixChan)

	for {

		//send ut heisens tilstand
		// chans.StateChan <- elevator
		
		switch elevator.CurrentState {		
		case 0:
			//initialiserer
			fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
			io.SetDoorOpenLamp(false)
			io.SetMotorDirection(io.MD_Down)
			select{
			case elevator.CurrentFloor = <- chans.FloorChan:
				io.SetMotorDirection(io.MD_Stop)
				io.SetFloorIndicator(elevator.CurrentFloor)
				elevator.CurrentState = 1
				fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
			}
				
		
		case 1:
			
			select{
			case buttonPress := <- chans.ButtonChan:
				if buttonPress.Button == io.BT_Cab {
					addCabHallToCurrentOrders(buttonPress, &currentOrders)
					io.SetButtonLamp(buttonPress.Button, buttonPress.Floor, true)
					fmt.Println(currentOrders)
				} else {
					chans.ReceivedFloorHallChan <- buttonPress
				}
			
			case elevator.CurrentFloor = <- chans.FloorChan:
				// i tilfelle heisen styres med tastene 7,8 eller 9
				io.SetFloorIndicator(elevator.CurrentFloor)
			case newOrder := <- NewOrderChan:
				temp := make([]io.ButtonEvent,1)
				temp[0] = newOrder
				currentOrders = append(temp, currentOrders...)
				
			default:
				if len(currentOrders) != 0 {
					elevator.NextFloor = currentOrders[0].Floor
					if elevator.NextFloor != elevator.CurrentFloor {
						elevator.CurrentState = 2
						fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
					} else {
						fmt.Println(currentOrders)
						elevator.CurrentState = 3
						fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
					}
				}
			}

		case 2:
			elevator.NextFloor = currentOrders[0].Floor
			if elevator.NextFloor < elevator.CurrentFloor{
				elevator.Dir = io.MD_Down
			}else if elevator.NextFloor > elevator.CurrentFloor{
				elevator.Dir = io.MD_Up
			}
			io.SetMotorDirection(elevator.Dir)

			select {
			case elevator.CurrentFloor = <- chans.FloorChan:
				fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
				io.SetFloorIndicator(elevator.CurrentFloor)
				for i := 0; i < len(currentOrders); i++ {
					if elevator.CurrentFloor == currentOrders[i].Floor {
						fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
						elevator.CurrentState = 3
					}
				}

			case buttonPress := <- chans.ButtonChan:
				if buttonPress.Button == io.BT_Cab {
					addCabHallToCurrentOrders(buttonPress, &currentOrders)
					io.SetButtonLamp(buttonPress.Button, buttonPress.Floor, true)
					fmt.Println(currentOrders)
				} else {
					chans.ReceivedFloorHallChan <- buttonPress
				}
			case newOrder := <- NewOrderChan:
				temp := make([]io.ButtonEvent,1)
				temp[0] = newOrder
				currentOrders = append(temp, currentOrders...)
			}

		case 3:
			fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
			io.SetMotorDirection(io.MD_Stop)

			for i := 0; i < len(currentOrders); i++ {
				if currentOrders[i].Floor == elevator.CurrentFloor {
					go toggleDoorLight(doorLightChan)
					if currentOrders[i].Button == io.BT_Cab {
						io.SetButtonLamp(io.BT_Cab, currentOrders[i].Floor, false)
					}
					currentOrders = append(currentOrders[:i],currentOrders[i+1:]...)
					fmt.Println(currentOrders)
				}
			}
			
			select{
			case <- doorLightChan:
				elevator.CurrentState = 1
			case buttonPress := <- chans.ButtonChan:
				if buttonPress.Button == io.BT_Cab {
					addCabHallToCurrentOrders(buttonPress, &currentOrders)
					io.SetButtonLamp(buttonPress.Button, buttonPress.Floor, true)
					fmt.Println(currentOrders)
				} else {
					chans.ReceivedFloorHallChan <- buttonPress
				}
			case newOrder := <- NewOrderChan:
				temp := make([]io.ButtonEvent,1)
				temp[0] = newOrder
				currentOrders = append(temp, currentOrders...)				
			}
		}

	}
}
