package core

import (
	"fmt"
	"runtime"
)

type Signal int

const (
	Pause  Signal = iota // Pauses the program
	Resume               // Resumes the program
	Stop                 // Terminates the program
)

func (s Signal) String() string {
	if s < Pause || s > Stop {
		return "Signal(" + string(s) + ")"
	}
	return [...]string{"Pause", "Resume", "Stop"}[s]
}

func (s *Signal) Parse(str string) error {
	switch str {
	case "Pause":
		*s = Pause
	case "Resume":
		*s = Resume
	case "Stop":
		*s = Stop
	default:
		return fmt.Errorf("invalid signal: %s", str)
	}
	return nil
}

func RunSignal(f func(), signaller <-chan Signal) {
	state := Resume

	for {
		select {
		case state = <-signaller:
			switch state {
			case Pause:
				fmt.Println("Pausing...")
			case Resume:
				fmt.Println("Resuming...")
			case Stop:
				fmt.Println("Stopping...")
				return
			}
		default:
			// We use runtime.Gosched() to prevent a deadlock in this case.
			// It will not be needed of work is performed here which yields
			// to the scheduler.
			runtime.Gosched()

			if state == Pause {
				break
			}

			f()
		}
	}
}
