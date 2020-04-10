package esm

import (
	"time"
	io "../elevio"
	c "../config"
)

func addCabHallToCurrentOrders(buttonPress io.ButtonEvent, currentOrders *[]int){
	//oppdater currentOrders
	shouldAppend := true
	for _, floor := range *currentOrders{
		if buttonPress.Floor == floor{
			shouldAppend = false
		}
	}
	if shouldAppend {
		*currentOrders = append(*currentOrders,buttonPress.Floor)
	}	
}

//funksjonen nedenfor skrur på og av dørlyset
func toggleDoorLight(doorLightChan chan int){
	io.SetDoorOpenLamp(true)
	time.Sleep(1*time.Second)
	io.SetDoorOpenLamp(false)
	doorLightChan <- 1
}

func toggleButtonLights(LightMatrixChan <-chan []bool){
	for{
		select{
		case lightMatrix := <- LightMatrixChan:
			for i := 0; i < len(lightMatrix); i++ {
				row, col := c.ListIndexToMatrixIndex(i, len(lightMatrix))
				button := io.BT_HallUp
				if row == 1 {
					button = io.BT_HallDown
				} else if row == 2 {
					button = io.BT_Cab
				}
				io.SetButtonLamp(button, col, lightMatrix[i])
		
			}
		default: 
			time.Sleep(20*time.Millisecond)
		}
	}
}



// func checkIfValidOrder()
//når heisen mottar en ordre må den sjekke om den henger i tråd med heisens tilstand
