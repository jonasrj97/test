// input: Struct ElevatorState
// output: Struct ElevatorState
// Funksjonalitet: leser av en struct, og backuper den til en fil.
// funksjonen restore leser av filen og returnerer den sist lagrede struct'en.

package backup

import (
	"encoding/gob"
	"fmt"
	"os"
)

//Import config elevatorstate strukten jee

//ElevatorState boi
type ElevatorState struct {
	IsMaster bool
	ID       int
	//CurrentState fsmState
	CurrentFloor int
	NextFloor    int
	//Dir          io.MotorDirection
}

//Backup - skriver structen til en fil.
func Backup(state ElevatorState) {
	dataFile, err := os.Create("backup.gob")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(state)
	dataFile.Close()
	println("Wrote to file.")
}

// Recover returnerer ElevatorState-Structen som er lagret i filen.
func Recover() ElevatorState {
	var state ElevatorState
	dataFile, err := os.Open("backup.gob")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&state)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return state
}
