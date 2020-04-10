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
//MÅ TESTES
func calculateCost(elev c.ElevatorState, orderFloor int, orderDir int, bestFit *c.BestFitStruct, wg *sync.WaitGroup){
	//kaller wg.Done() når funksjonen er ferdig
	defer wg.Done()
	var cost, elevDir, difference int
	//i enkelte tilfeller er ikke heisen skikket i det hele tatt og bør overses
	notFitted := false
	//konverterer heisens retning til en int
	if elev.Dir == io.MD_Up{
		elevDir = 1
	} else if elev.Dir == io.MD_Stop{
		elevDir = 0
	} else if elev.Dir == io.MD_Down{
		elevDir = -1
	}

	//utregning av den aktuelle heisens kostnad for å utføre den aktuelle bestillingen
	difference = elev.CurrentFloor-orderFloor
	if difference < 0{
		difference = -difference
	}
	cost += difference

	if elev.CurrentState != 2 {
		cost -= 1
	} else if elevDir == orderDir{
		cost -= 1
	}

	//structen må låses slik at ikke flere go-rutiner forandrer på den samtidig. 
	if (cost < bestFit.Cost) && !notFitted {
		bestFit.Mux.Lock()
		bestFit.ID = elev.ID
		bestFit.Cost = cost
		bestFit.Mux.Unlock()
	}
}

//sjekke tilstandene til heisene for om en ordre har blitt fullført
//IKKE TESTET
func checkForCompletedOrders(activeElevators map[int]c.ElevatorState, orderMatrix *[]c.MatrixElement){
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
//fjerner ordre som er fullført
//UNØDVENDIG SIDEN VI KUTTER UT COMPLETEDORDERSCHAN
// func removeCompletedOrders(completedOrdersChan <-chan esm.ElevatorState, orderMatrix *[]OrderMatrixElement, wg *sync.WaitGroup){
// 	loop:
// 		for{
// 			select{
// 			case completedOrder := <- completedOrdersChan:
// 				//kan immideately remove the cab order, hence row = 2
// 				row := 2
// 				i := matrixIndexToListIndex(row, completedOrder.CurrentFloor)
// 				(*orderMatrix)[i].PlacedOrder 	= false
// 				(*orderMatrix)[i].Active	 	= false
// 				(*orderMatrix)[i].ElevID 		= 0
// 				//koden nedenfor antar at heisens CurrentDirection ikke forandres etter at den stopper. Den må først forandres når ny ordre er blitt gitt.
// 				if completedOrder.Dir == io.MD_Down {
// 					row = 1
// 				} else if completedOrder.Dir == io.MD_Up{
// 					row = 0
// 				}
// 				j := matrixIndexToListIndex(row, completedOrder.CurrentFloor)
// 				(*orderMatrix)[j].PlacedOrder 	= false
// 				(*orderMatrix)[j].Active	 	= false
// 				(*orderMatrix)[j].ElevID 		= 0
// 			default:
// 				wg.Done()
// 				break loop
// 			}
// 		}
// }

//sjekker om ordre som er aktive egt. skal være det, i tilfelle en heis skulle dette ut
func checkActiveOrders(activeElevators map[int]c.ElevatorState, orderMatrix *[]c.MatrixElement){
	indexActiveOrders := make([]int, 0, 100)
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
func updateActiveElvatorsList(channel <-chan c.ElevatorState/*map[int]esm.ElevatorState*/, activeElevators *map[int]c.ElevatorState, wg *sync.WaitGroup){	
	loop:
		for{
			select{
			// case *activeElevators = <- channel:
			case elev := <- channel:
				(*activeElevators)[elev.ID] = elev
			default:	
				wg.Done()
				break loop
		}
	}
}