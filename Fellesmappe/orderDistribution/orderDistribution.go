package orderDistribution

import (
	// "fmt"
	// "sync"
	// "../esm"
	// "fmt"
	io "../elevio"
	c "../config"
)

var isMaster bool


//hovedfunksjon
func DistributeOrders(chans c.OrderDistributionChannels, ReceivedFloorHallChan <-chan io.ButtonEvent, /*ActiveElevatorsChan <-chan map[int]c.ElevatorState,*/ numFloors int){
	floorHallMatrix := make([]c.MatrixElement, numFloors*2)
	// activeElevators := make(map[int]c.ElevatorState)
	isMaster := true
	ID := 1
	// var wg sync.WaitGroup
	for{
		//oppdater ordrematrisen med den siste oppdateringen av den som er synkronisert med de andre
		loop:
			for{
				select{
				case syncFloorHallMatrix := <- chans.FloorHallMatrixChan:
					for i := 0; i < len(floorHallMatrix); i++ {
						//sørg for at cabHalls ikke synkroniseres med de andre heisene
						floorHallMatrix[i].PlacedOrder 	= syncFloorHallMatrix[i].PlacedOrder
						floorHallMatrix[i].ElevID 		= syncFloorHallMatrix[i].ElevID
					}
					lightMatrix := make([]bool, len(floorHallMatrix))
					for i := 0; i < len(lightMatrix); i++ {
						lightMatrix[i] = floorHallMatrix[i].PlacedOrder
					}
					chans.LightMatrixChan <- lightMatrix
				
				case newFloorHall := <- ReceivedFloorHallChan:
					col := newFloorHall.Floor
					row := 0
					if newFloorHall.Button == io.BT_HallDown {
						row = 1
					}
					i := c.MatrixIndexToListIndex(row,col)
					floorHallMatrix[i].PlacedOrder = true

				default:
					break loop
				}
			}	
		
		//hvis det har kommet en ordre fra master, skal denne forrerst i køen.
		for i := 0; i < len(floorHallMatrix); i++ {
			if floorHallMatrix[i].ElevID == ID {
				_, floor := c.ListIndexToMatrixIndex(i, len(floorHallMatrix))
				chans.NewOrderChan <- floor
				break
			}
		}


		//send lokal ordrematrise til master
		//FloorHallMatrixChan <- floorHallMatrix

		//hvis isMaster == true skal det også deles ut ordrer til de andre heisene
		if isMaster {

			//oppdater listen med aktive heiser og hvilke tilstander de har
			//update activeElevators() <- activeElevatorsChan

			//motta synkronisert floorHallMatrix fra Network
			
			//sjekk om en ordre har blitt fullført og om en ordre kan fjernes fra syncFloorHallMatrix
			
			//sjekk om det er anledning til å sende noen av heisene ordre fra syncFloorHallMatrix (calculateCost())
			//oppdater matrisen i så fall med riktig elevID

			//floorHallMatrixChan <- syncFloorHallMatrix (sørg for at channelen er tom først)


		}


	}
}