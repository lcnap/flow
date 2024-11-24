package event

type Phase int

const (
	PhaseAuth Phase = 1
)

type Event interface {
	Fire()
}

var (
	eventMap map[Phase][]Event = map[Phase][]Event{}
)

func RegisterEvent(e Event, phase Phase) {
	if events, ok := eventMap[phase]; ok {
		events = append(events, e)
		eventMap[phase] = events
	} else {
		events = append([]Event{}, e)
		eventMap[phase] = events
	}
}

func FireAt(phase Phase) {
	events := eventMap[phase]
	for _, e := range events {
		e.Fire()
	}
}
