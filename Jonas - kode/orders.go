package main

import (
	"fmt"
)

// Vet ikke helt hva som må importeres...

const max_number_of_elevators = 10
const number_of_floors = 4

type STATE int
	const(
		IDLE STATE = 0
		RUNNING STATE = 1
		DOORS_OPEN STATE = 2
		DISCONNECTED STATE = 3
	)

type MotorDirection int
	const(
		MD_Up   MotorDirection = 1
		MD_Down                = -1
		MD_Stop                = 0
	)

type Elevator struct {
	Active bool
    State STATE
    CurrentDirection MotorDirection
	CurrentFloor int
	NextFloor int
	ID int
}

var elevator_list [max_number_of_elevators]Elevator

var order_matrix [number_of_floors][3]int

var desired_floor int

var closest_difference int
var closest_index int
var difference_temp int

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//func choose_elevator(elevator_list []Elevator, desired_floor int, max_number_of_elevators int) int {

func choose_elevator(desired_floor int) int {

	//returns the ID of the elevator that should take the order

	if (desired_floor == -1) {
		return -1
	}

	closest_difference = 1000
	closest_index = 1000

	for i := 0; i < max_number_of_elevators; i++ {
		if (elevator_list[i].Active == true){
			if ((elevator_list[i].State == RUNNING) && (elevator_list[i].CurrentDirection == MD_Up)) {
				difference_temp = desired_floor - elevator_list[i].CurrentFloor
				if difference_temp < 0 {
					difference_temp = number_of_floors + Abs(difference_temp)
				}
				if difference_temp < closest_difference {
					closest_difference = difference_temp
					closest_index = i;
				}
			}
			if ((elevator_list[i].State == RUNNING) && (elevator_list[i].CurrentDirection == MD_Down)) {
				difference_temp = elevator_list[i].CurrentFloor - desired_floor
				if difference_temp < 0 {
					difference_temp = number_of_floors + Abs(difference_temp)
				}
				if difference_temp < closest_difference {
					closest_difference = difference_temp
					closest_index = i;
				}
			}
			if (elevator_list[i].State == IDLE) {
				difference_temp = Abs(desired_floor - elevator_list[i].CurrentFloor)
				if difference_temp < closest_difference {
					closest_difference = difference_temp
					closest_index = i;
				}
			}
		}
	}
	if (closest_index != 1000){
		return elevator_list[closest_index].ID
	} else {
		return -1
	}
}

//func update_elevator_list() ([]Elevator, int) {
	// fyll inn kode for å oppdatere elevator_list. Via channels? Nettverk?
	// oppdater max_number_of_elevators. Vil være lengden til elevator_list.
	// Returner elevator_list og max_number_of_elevators
//}

func get_order() int {
	for i := 0; i < number_of_floors; i++ {
		for j := 0; j < 3; j++ {
			if order_matrix[i][j] == 1 {
				order_matrix[i][j] = 0;
				return i
			}
		}
	}
	return -1
	// Hent en ny bestilling, og returner hvilken etasje en av heisene må reise til.
}

func orders_init() {
	for i := 0; i < max_number_of_elevators; i++ {
		elevator_list[i].Active = false
		elevator_list[i].State = IDLE
		elevator_list[i].CurrentDirection = MD_Stop
		elevator_list[i].CurrentFloor = -1
		elevator_list[i].NextFloor = -1
		elevator_list[i].ID = i
	}
	for i := 0; i < number_of_floors; i++ {
		for j := 0; j < 3; j++{
			order_matrix[i][j] = 0
		}
	}
}

func orders_set_elevator(x int) {
	elevator_list[x].Active = true
}

func orders_reset_elevator(x int) {
	elevator_list[x].Active = false
}
 
func main() {
	for i := 0; i < number_of_floors; i++ {
		for j := 0; j < 3; j++ {
			order_matrix[i][j] = 0
		}
	}

	//fmt.Printf("%d", elevator_list[3].CurrentFloor)

	orders_init()

	//fmt.Println(order_matrix)
	//fmt.Println(elevator_list)

	
	elevator_list[0].CurrentFloor = 3
	elevator_list[0].Active = true
	elevator_list[0].CurrentDirection = MD_Up
	elevator_list[0].State = RUNNING

	/*
	
	elevator_list[1].CurrentFloor = 1
	elevator_list[1].Active = true
	elevator_list[1].CurrentDirection = MD_Up
	elevator_list[1].State = RUNNING
	
	elevator_list[2].CurrentFloor = 6
	elevator_list[2].Active = true
	elevator_list[2].CurrentDirection = MD_Down
	elevator_list[2].State = RUNNING

	elevator_list[3].CurrentFloor = 5
	elevator_list[3].Active = true
	elevator_list[3].CurrentDirection = MD_Down
	elevator_list[3].State = RUNNING
	
	elevator_list[4].CurrentFloor = 8
	elevator_list[4].Active = true
	elevator_list[4].CurrentDirection = MD_Down
	elevator_list[4].State = IDLE
	*/
	

	fmt.Println(elevator_list)

	order_matrix[2][2] = 1
	order_matrix[3][2] = 1
	order_matrix[1][2] = 1
	fmt.Println(order_matrix)
	fmt.Println(choose_elevator(get_order()))
	fmt.Println(order_matrix)
}

// Jeg vet ikke helt hvilke typer og konstanter som må deklareres i denne koden. 
// Sannsynligvis hentes mye fra andre moduler, men skrev en del inn her enn så lenge.