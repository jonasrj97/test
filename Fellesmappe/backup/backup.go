// input: Struct ElevatorState
// output: Struct ElevatorState
// Funksjonalitet: leser av en struct, og backuper den til en fil.
// funksjonen restore leser av filen og returnerer den sist lagrede struct'en.

package backup

import (
	"encoding/gob"
	"fmt"
	"os"

	c "../config"
)

//Backup - skriver structen til en fil.
// func Backup(state config.ElevatorState) {
// 	dataFile, err := os.Create("backup.gob")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	dataEncoder := gob.NewEncoder(dataFile)
// 	dataEncoder.Encode(state)
// 	dataFile.Close()
// }

func Backup(currentFloor int, currentOrders []int) {
	state := c.Backup{currentFloor, currentOrders}
	dataFile, err := os.Create("backup.gob")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(state)
	dataFile.Close()
}

func Recover() (state c.Backup, err error) {
	state = c.Backup{0,nil}
	dataFile, err := os.Open("backup.gob")
	if err != nil {
		fmt.Println(err)
		return
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&state)

	if err != nil {
		fmt.Println(err)
	}

	deleteError := os.Remove()
	return
}

// Recover returnerer ElevatorState-Structen som er lagret i filen.
// func Recover() config.ElevatorState {
// 	var state config.ElevatorState
// 	dataFile, err := os.Open("backup.gob")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	dataDecoder := gob.NewDecoder(dataFile)
// 	err = dataDecoder.Decode(&state)

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	return state
// }
