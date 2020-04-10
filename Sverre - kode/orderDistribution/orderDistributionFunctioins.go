package orderDistribution

import (
	io "../elevio"
	c "../config"
	// "../esm"
	"strconv"
	"fmt"
	"sync"
)


//funksjonen regner ut hvilken heis som er best skikket til å utføre den aktuelle bestillingen
func calculateCost(state c.ElevatorState, orderFloor int, orderDir io.MotorDirection, bestFit *c.BestFitStruct, wg *sync.WaitGroup) {
	defer wg.Done()

	cost := 0
	notFeasible := false 

	if state.Dir != orderDir && orderDir != io.MD_Stop {
		notFeasible = true
	} else if state.CurrentState == 0 {
		notFeasible = true
	} else if state.CurrentState == 2 && state.CurrentFloor == orderFloor {
		notFeasible = true
	} else if state.Dir == io.MD_Down && orderDir == io.MD_Down && state.NextFloor > orderFloor {
		notFeasible = true
	} else if state.Dir == io.MD_Up && orderDir == io.MD_Up && state.NextFloor < orderFloor {
		notFeasible = true
	}

	if state.Dir == io.MD_Down && orderDir == io.MD_Down && orderFloor >= state.NextFloor {
		cost -= 1
	} else if state.Dir == io.MD_Up && orderDir == io.MD_Up && orderFloor <= state.NextFloor {
		cost -= 1
	}

	if state.Dir == 0 {
		cost -= 2
	}
	
	diff := orderFloor - state.CurrentFloor
	if diff < 0 {
		diff = -diff
	}
	cost += diff

	if !notFeasible && cost <= (*bestFit).Cost && state.ID < (*bestFit).ID {
		(*bestFit).Mux.Lock()
		(*bestFit).ID = state.ID
		(*bestFit).Cost = cost
		(*bestFit).Mux.Unlock()
	}
}

//sjekke tilstandene til heisene for om en ordre har blitt fullført
//IKKE TESTET
func removeCompletedOrders(activeElevators map[int]c.ElevatorState, orderMatrix *[]c.MatrixElement){
	var floor int
	for _, element := range activeElevators {
		for i := 0; i < len(*orderMatrix); i++ {
			_, floor = c.ListIndexToMatrixIndex(i, len(*orderMatrix))
			if (element.ID == (*orderMatrix)[i].ElevID) && (element.CurrentFloor == floor) && (element.CurrentState == 3){
				(*orderMatrix)[i].PlacedOrder 	= false
				(*orderMatrix)[i].ElevID 		= 0
			}
		}
	}
}

//helpefunksjon for å oppdatere bestillingsmatrisen
// func updateLocalOrderMatrix(buttonChan <-chan io.ButtonEvent, orderMatrix *[]OrderMatrixElement, currentOrders *[]int, wg *sync.WaitGroup){
// 	loop:
// 		for{
// 			select{
// 			case buttonPress := <- buttonChan:
// 				var row int
// 				//oppdater den lokale ordrematrisen
// 				if buttonPress.Button == io.BT_HallUp{
// 					row = 0
// 				} else if buttonPress.Button == io.BT_HallDown{
// 					row = 1
// 				}
// 				i := matrixIndexToListIndex(row, buttonPress.Floor)
// 				(*orderMatrix)[i].PlacedOrder = true

// 				//oppdater currentOrders
// 				if buttonPress.Button == io.BT_Cab {
// 					checkLoop:
// 						for _, floor := range *currentOrders{
// 							if buttonPress.Floor == floor{
// 								break checkLoop
// 							}
// 						}
// 						*currentOrders = append(*currentOrders,buttonPress.Floor)	
// 					}
					
// 			default:
// 				wg.Done()
// 				break loop
// 			}	
// 		}
// }

func addLocalHallsToOrderMatrix(buttonPress io.ButtonEvent, orderMatrix *[]c.MatrixElement){
	var row int
	if buttonPress.Button == io.BT_HallUp{
		row = 0
	} else if buttonPress.Button == io.BT_HallDown{
		row = 1
	}
	i := c.MatrixIndexToListIndex(row, buttonPress.Floor)
	(*orderMatrix)[i].PlacedOrder = true
}

//sjekker om ordre som er aktive egt. skal være det, i tilfelle en heis skulle dette ut
func checkActiveOrders(activeElevators map[int]c.ElevatorState, orderMatrix *[]c.MatrixElement){
	indexActiveOrders := make([]int, 0)
	for _, elem := range activeElevators{
		col := elem.NextFloor
		row := 2 
		indexActiveOrders = append(indexActiveOrders, c.MatrixIndexToListIndex(row, col))
		if elem.Dir == io.MD_Up{
			row = 0
			indexActiveOrders = append(indexActiveOrders, c.MatrixIndexToListIndex(row, col))
		} else if elem.Dir == io.MD_Down{
			row = 1
			indexActiveOrders = append(indexActiveOrders, c.MatrixIndexToListIndex(row, col))
		}
	}
	var correct bool
	for i := 0; i < len(*orderMatrix); i++{
		correct = false
		if (*orderMatrix)[i].ElevID == 0 {
			for j := 0; j < len(indexActiveOrders); j++ {
				if i == indexActiveOrders[j] {
					correct = true
					break
				}
			}
			if correct == false {
				(*orderMatrix)[i].ElevID = 0
			}
		}
	}
}

//hjelpefunksjon som printer ordrematrisen
//hallUp 	--> [etg_2 ... ... etg_n  ]
//hallDown 	--> [etg_1 ... ... etg_n-1]
//cab		--> [etg_1 ... ... etg_n  ]
func printOrderMatrix(matrix []c.MatrixElement){
	var output 				string
	var outputPlacedOrder 	string
	var outputElevID		string
	for i := 0; i < 3; i++{
		output = 			""
		outputPlacedOrder = ""
		outputElevID = 		""
		for j := 0; j < 4; j++{
			if (i == 0 && j == 3) || (i == 1 && j == 0){
				outputPlacedOrder 	+= "|                |"
				outputElevID		+= "|                |"
				continue
			} 
			index := c.MatrixIndexToListIndex(i, j)
			
			outputPlacedOrder 	+= "| Order  = " + strconv.FormatBool(matrix[index].PlacedOrder) 
			outputElevID		+= "| ElevID = " + strconv.Itoa(matrix[index].ElevID)
			if matrix[index].PlacedOrder{
				outputPlacedOrder += " "
			} 
			outputPlacedOrder += " |"
			outputElevID += "     |"
		}
		output = outputPlacedOrder + "\n" + outputElevID + "\n"
		fmt.Printf(output)
		fmt.Printf("\n")
	}
}


//hjelpefunksjon som oppdaterer listen over hvilke heiser som er aktive.
func updateActiveElevators(channel <-chan map[int]c.ElevatorState, activeElevators *map[int]c.ElevatorState){	
	loop:
		for{
			select{
			case *activeElevators = <- channel:
			default:	
				break loop
		}
	}
}