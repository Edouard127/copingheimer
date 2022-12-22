package packet

type State uint

const (
	Offline State = iota
	Wait
)

type States struct {
	states map[State]bool
}

func NewStates() *States {
	return &States{
		states: make(map[State]bool),
	}
}

func (s *States) Set(state State, value bool) {
	s.states[state] = value
}

func (s *States) Get(state State) bool {
	return s.states[state]
}

func (s *States) Has(state State) bool {
	if st, ok := s.states[state]; ok {
		return st == true
	}
	return false
}

func (s *States) SetAll(value bool) {
	for k := range s.states {
		s.states[k] = value
	}
}

func (s *States) SetIf(state State, value bool, condition bool) {
	if condition {
		s.states[state] = value
	}
}

func ToInt(s bool) int32 {
	if s {
		return 1
	}
	return 0
}
