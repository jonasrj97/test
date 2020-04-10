package config

import (
	io "../elevio"
	"sync"
)

type Backup struct {
	CurrentFloor 	int 
	CurrentOrders	[]int
}

//struct med egenskapene til den heisen som passer best til å utføre en aktuell bestilling.
type BestFitStruct struct{
	Mux		sync.Mutex
	ID 		int
	Cost 	int
}

//hvert element i ordrematrisen er en struct med et felt som er sant hvis det har kommet en bestilling til den aktuelle etasjen og et felt som er sant
//hvis bestillingen holder på å bli fullført.
type MatrixElement struct{
	PlacedOrder			bool
	ElevID				int 
}

//chans:
type OrderDistributionChannels struct{
	ActiveElevatorsChan		chan map[int]ElevatorState
	NewOrderChan			chan io.ButtonEvent
	IsMasterChan			chan bool
	LightMatrixChan			chan []bool
	FloorHallMatrixChan		chan []MatrixElement
}

type fsmState int

const (
	Initializing fsmState = 0
	Idle                  = 1
	Running               = 2
	OpenDoors             = 3
)

type ElevatorState struct {
	IsMaster	 bool
	ID 			 int
	CurrentState fsmState
	CurrentFloor int
	NextFloor    int
	Dir          io.MotorDirection	
}

type ElevatorChannels struct {
	FloorChan  				chan int
	StateChan  				chan ElevatorState
	ButtonChan				chan io.ButtonEvent
	ReceivedFloorHallChan	chan io.ButtonEvent
	KillChan				chan bool
}

//hjelpefunksjon som konverterer listeindeks til matrisekoordinater
func MatrixIndexToListIndex(row int, col int) int { //zero indexing for both col and row
	if row == 0{
		return row * 3 + col
	} else if row == 1 {
		return row * 3 + col + 1
	} else {
		return row * 3 + col + 2
	}
}

//hjelpefunksjon som konverterer matrisekoodrinater til listeindeks
func ListIndexToMatrixIndex(i int, length int) (row int, col int) {
	nFloors := length/2
	if i < nFloors{
		row = 0
		col = i
	} else {
		row = 1
		col = i-nFloors
	}
	return 
}