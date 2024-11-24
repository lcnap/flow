package event

import (
	"log"
	"testing"
)

type LogEvent struct {
	Id string
}

func (l LogEvent) Fire() {
	log.Printf("fired log event: %s\n", l.Id)
}

func TestFireAt(t *testing.T) {
	RegisterEvent(LogEvent{"1"}, PhaseAuth)
	RegisterEvent(LogEvent{"2"}, PhaseAuth)

	FireAt(PhaseAuth)
}
