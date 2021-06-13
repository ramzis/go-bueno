package bueno

type StateValue string

const (
	STATE_Null                             StateValue = "STATE_Null"
	STATE_Stop                             StateValue = "STATE_Stop"
	STATE_NewGame                          StateValue = "STATE_NewGame"
	STATE_InitVars                         StateValue = "STATE_InitVars"
	STATE_InitDeck                         StateValue = "STATE_InitDeck"
	STATE_DealStartingHands                StateValue = "STATE_DealStartingHands"
	STATE_DrawFirstCard                    StateValue = "STATE_DrawFirstCard"
	STATE_StartGame                        StateValue = "STATE_StartGame"
	STATE_GameRunning                      StateValue = "STATE_GameRunning"
	STATE_GameRunning_WaitingForMove       StateValue = "STATE_GameRunning_WaitingForMove"
	STATE_GameRunning_MakingMove           StateValue = "STATE_GameRunning_MakingMove"
	STATE_GameRunning_ResolvingMoveEffects StateValue = "STATE_GameRunning_ResolvingMoveEffects"
)

type state struct {
	value    StateValue
	onChange chan StateValue
}

type State interface {
	Set(StateValue)
	Get() StateValue
	OnChangeChan() <-chan StateValue
}

func (s *state) Set(state StateValue) {
	if s.value == state {
		return
	}
	s.value = state
	if s.value == STATE_Stop {
		close(s.onChange)
		return
	}
	s.onChange <- s.value
}

func (s *state) Get() StateValue {
	return s.value
}

func (s *state) OnChangeChan() <-chan StateValue {
	return s.onChange
}

func NewState() State {
	return &state{
		value:    STATE_Null,
		onChange: make(chan StateValue),
	}
}
