package main

import "fmt"

// Vet ikke helt hva som må importeres...

var number_of_elevators := 3

type STATE int
	const(
		IDLE STATE = 0
		RUNNING STATE = 1
		DOORS_OPEN STATE = 2
		DISCONNECTED STATE = 3
	)

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type Elevator struct {
    State STATE
    CurrentDirection MotorDirection
	CurrentFloor int
	NextFloor int
	ID int
}

var elevator_list := [number_of_elevators]Elevator

var desired_floor int

closest_difference = inf

func choose_elevator(elevator_list []Elevator, desired_floor int, number_of_elevators int) int {

	var closest_difference int
	var closest_index int
	var difference_temp int

	closest_difference = inf

	for i := 0; i < number_of_elevators; i++ {
		if ((elevator_list[i].STATE == IDLE) || ((elevator_list[i].STATE == MOVING) && (elevator_list[i].CurrentDirection == MD_Up)) && (elevator_list[i].CurrentFloor < desired_floor)) {
			difference_temp = desired_floor - elevator_list[i].CurrentFloor
			if difference_temp < closest_difference {
				closest_difference = difference_temp
				closest_index = i;
			}
		}
		if ((elevator_list[i].STATE == IDLE) || ((elevator_list[i].STATE == MOVING) && (elevator_list[i].CurrentDirection == MD_Down)) && (elevator_list[i].CurrentFloor > desired_floor)) {
			difference_temp = elevator_list[i].CurrentFloor - desired_floor
			if difference_temp < closest_difference {
				closest_difference = difference_temp
				closest_index = i;
			}
		}

		return closest_index
	}
}

func update_elevator_list() ([]Elevator, int) {
	// fyll inn kode for å oppdatere elevator_list. Via channels? Nettverk?
	// oppdater number_of_elevators. Vil være lengden til elevator_list.
	// Returner listen over Elevators, og antall Elevators.
}

func get_order() int {
	// Hent en ny bestilling, og returner hvilken etasje en av heisene må reise til.
}

func main() {
	var ID_to_take_order int

	elevator_list, number_of_elevators = update_elevator_list()
	desired_floor = get_order()
	elevator_index = choose_elevator(elevator_list, desired_floor, number_of_elevators)
	ID_to_take_order = elevator_list[elevator_index].ID

	// send order to elevator with ID ID_to_take_order. Via UDP?
}

// Jeg vet ikke helt hvilke typer og konstanter som må deklareres i denne koden. 
// Sannsynligvis hentes mye fra andre moduler, men skrev en del inn her enn så lenge.
