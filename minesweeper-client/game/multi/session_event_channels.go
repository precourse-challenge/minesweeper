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
		StartChan:    make(chan StartEvent, 5),
		UpdateChan:   make(chan UpdateEvent, 5),
		ErrorChan:    make(chan ErrorEvent, 5),
		GameOverChan: make(chan GameOverEvent, 5),
		JoinedChan:   make(chan JoinedEvent, 5),
	}
}
