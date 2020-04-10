package orderDistribution

import (
	// "fmt"
	"sync"
	// "../esm"
	io "../elevio"
	c "../config"
)

//hovedfunksjon
func DistributeOrders(chans c.OrderDistributionChannels, ReceivedFloorHallChan <-chan io.ButtonEvent, numFloors int){
	floorHallMatrix := make([]c.MatrixElement, numFloors*2)
	activeElevators := make(map[int]c.ElevatorState)
	isMaster := true
	ID := 1
	var wg sync.WaitGroup
	for{

		//må kontinuerlig sjekke om denne statusen har forandret seg
		select{
		case isMaster = <- chans.IsMasterChan:
		default:
		}

		//oppdater ordrematrisen med den siste oppdateringen av den som er synkronisert med de andre
		loop:
			for{
				select{
				case syncFloorHallMatrix := <- chans.FloorHallMatrixChan:
					for i := 0; i < len(floorHallMatrix); i++ {
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
				row, col := c.ListIndexToMatrixIndex(i, len(floorHallMatrix))
				var newOrder io.ButtonEvent
				if row == 0 {
					newOrder.Button = io.BT_HallUp
				} else {
					newOrder.Button = io.BT_HallDown
				}
				newOrder.Floor = col
				chans.NewOrderChan <- newOrder
				break
			}
		}


		//hvis isMaster == true skal det også deles ut ordrer til de andre heisene
		if isMaster {

			updateActiveElevators(chans.ActiveElevatorsChan, &activeElevators)

			//hvilke heiser har allerede mottatt -- tror ikke denne er nødvendig 
			// preoccupiedElevs := make([]int,0)
			// for i := 0; i < len(floorHallMatrix); i++ {
			// 	if floorHallMatrix[i].PlacedOrder && floorHallMatrix[i].ElevID != 0 {
			// 		preoccupiedElevs = append(preoccupiedElevs,floorHallMatrix[i].ElevID)
			// 	}
			// }
			
			//sjekk om ordrene som har blitt utdelt fortsatt er aktive
			checkActiveOrders(activeElevators, &floorHallMatrix)
			//sjekk om en ordre har blitt fullført og om en ordre kan fjernes fra syncFloorHallMatrix
			removeCompletedOrders(activeElevators, &floorHallMatrix)
			
			//sjekk om det er anledning til å sende noen av heisene ordre fra syncFloorHallMatrix (calculateCost())
			//oppdater matrisen i så fall med riktig elevID
			for i, order := range floorHallMatrix {
				if order.PlacedOrder && order.ElevID == 0 {
					bestFit := c.BestFitStruct{
						ID: 	0,
						Cost:	100,
					}
					
					row, col := c.ListIndexToMatrixIndex(i,len(floorHallMatrix))
					dir := io.MD_Up
					if row == 1 {
						dir = io.MD_Down
					}

					for _, elev := range activeElevators {
						wg.Add(1)
						calculateCost(elev,col,dir,&bestFit,&wg)
					}
					wg.Wait()

					order.ElevID = bestFit.ID
				}
			}
		}

		//send lokal ordrematrise til network for synkronisering (og eventuelt til master)
		chans.FloorHallMatrixChan <- floorHallMatrix

	}
}