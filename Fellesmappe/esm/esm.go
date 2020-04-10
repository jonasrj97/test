package esm

import (
	io "../elevio"
	c "../config"
	b "../backup"
	"fmt"
	"os"
)

func RunElevator(chans c.ElevatorChannels, NewOrderChan <-chan int, LightMatrixChan <-chan []bool, killSig <-chan os.Signal) {
	currentOrders := make([]int,0,10) //cab halls
	doorLightChan := make(chan int)

	// motta ID og isMaster fra network
	elevator := c.ElevatorState{
		ID:           1, //må hentes fra nettverksmodulen
		CurrentState: 0,
		CurrentFloor: -1,
		NextFloor:    -1,
		Dir:          io.MD_Stop,
	}

	//overlagrer heisens tilstand og currentOrders hvis det finnes en backup
	backup, err := b.Recover()
	if err == nil {
		fmt.Printf("\nElevator resumes last known state.\n")
		fmt.Println(currentOrders)
		elevator.CurrentFloor = backup.CurrentFloor
		elevator.CurrentState = 1
		currentOrders = backup.CurrentOrders
	}
	
	go toggleButtonLights(LightMatrixChan)
	//kill-funksjonen ligger her slik at heisens tilstand kan lagres før prosessen termineres.
	go func(){
		select{
		case <- killSig:
			b.Backup(elevator.CurrentFloor, currentOrders)
			fmt.Printf("\n\n\nelevator is kill. noooooooo\n\n\n")
			os.Exit(1)
		}
	}()

	for {

		//send ut heisens tilstand
		// chans.StateChan <- elevator

		switch elevator.CurrentState {
		case 0:
			//initialiserer
			fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)

			io.SetMotorDirection(io.MD_Down)
			select {
			case elevator.CurrentFloor = <-chans.FloorChan:
				io.SetMotorDirection(io.MD_Stop)
				io.SetFloorIndicator(elevator.CurrentFloor)
				elevator.CurrentState = 1
				fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
			}

		case 1:

			select {
			case buttonPress := <-chans.ButtonChan:
				if buttonPress.Button == io.BT_Cab {
					addCabHallToCurrentOrders(buttonPress, &currentOrders)
				} else {
					chans.ReceivedFloorHallChan <- buttonPress
				}

			case elevator.CurrentFloor = <-chans.FloorChan:
				// i tilfelle heisen styres med tastene 7,8 eller 9
				io.SetFloorIndicator(elevator.CurrentFloor)
			case newOrder := <-NewOrderChan:
				currentOrders = append(currentOrders, newOrder)

			default:
				if len(currentOrders) != 0 {
					elevator.NextFloor = currentOrders[0]
					if elevator.NextFloor != elevator.CurrentFloor {
						elevator.CurrentState = 2
						fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
					} else {
						elevator.CurrentState = 3
						fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
					}
				}
			}

		case 2:
			if elevator.NextFloor < elevator.CurrentFloor {
				elevator.Dir = io.MD_Down
			} else if elevator.NextFloor > elevator.CurrentFloor {
				elevator.Dir = io.MD_Up
			}
			io.SetMotorDirection(elevator.Dir)

			select {
			case elevator.CurrentFloor = <-chans.FloorChan:
				fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
				io.SetFloorIndicator(elevator.CurrentFloor)
				for i, floor := range currentOrders {
					if elevator.CurrentFloor == floor {
						currentOrders = append(currentOrders[:i], currentOrders[i+1:]...)

						fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
						elevator.CurrentState = 3
					}
				}
			case buttonPress := <-chans.ButtonChan:
				if buttonPress.Button == io.BT_Cab {
					addCabHallToCurrentOrders(buttonPress, &currentOrders)
				} else {
					chans.ReceivedFloorHallChan <- buttonPress
				}
			case newOrder := <-NewOrderChan:
				currentOrders = append(currentOrders, newOrder)
			}

		case 3:
			io.SetMotorDirection(io.MD_Stop)
			go toggleDoorLight(doorLightChan)

			select {
			case <-doorLightChan:
				if elevator.CurrentFloor == elevator.NextFloor {
					elevator.CurrentState = 1
					fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
				} else {
					elevator.CurrentState = 2
					fmt.Printf("%v %d %v %d %v %d %v %d", "\n\nState:", elevator.CurrentState, "\nDirection: ", elevator.Dir, "\nCurrent floor: ", elevator.CurrentFloor, "\nNext floor: ", elevator.NextFloor)
				}

			}
		}

	}
}
