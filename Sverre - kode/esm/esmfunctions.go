package esm

import (
	"time"
	io "../elevio"
	c "../config"
)

func addCabHallToCurrentOrders(buttonPress io.ButtonEvent, currentOrders *[]io.ButtonEvent){
	//oppdater currentOrders
	shouldAppend := true
	for i := 0; i < len(*currentOrders); i++ {
		if buttonPress == (*currentOrders)[i] {
			shouldAppend = false
		}
	}
	if shouldAppend {
		*currentOrders = append(*currentOrders,buttonPress)
	}	
}

//funksjonen nedenfor skrur på og av dørlyset
func toggleDoorLight(doorLightChan chan int){
	io.SetDoorOpenLamp(true)
	time.Sleep(2*time.Second)
	io.SetDoorOpenLamp(false)
	doorLightChan <- 1
}

func toggleFloorButtonLights(LightMatrixChan <-chan []bool){
	for{
		select{
		case lightMatrix := <- LightMatrixChan:
			for i := 0; i < len(lightMatrix); i++ {
				row, col := c.ListIndexToMatrixIndex(i, len(lightMatrix))
				button := io.BT_HallUp
				if row == 1 {
					button = io.BT_HallDown
				} 
				io.SetButtonLamp(button, col, lightMatrix[i])
		
			}
		default: 
			time.Sleep(20*time.Millisecond)
		}
	}
}
