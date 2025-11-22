package multi

type SessionEventChannels struct {
	StartChan    chan StartEvent
	UpdateChan   chan UpdateEvent
	ErrorChan    chan ErrorEvent
	GameOverChan chan GameOverEvent
	JoinedChan   chan JoinedEvent
}

func NewSessionEventChannels() *SessionEventChannels {
	return &SessionEventChannels{
		StartChan:    make(chan StartEvent),
		UpdateChan:   make(chan UpdateEvent),
		ErrorChan:    make(chan ErrorEvent),
		GameOverChan: make(chan GameOverEvent),
		JoinedChan:   make(chan JoinedEvent),
	}
}
