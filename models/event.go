package models

type EventCategory string

const (
	Movie    EventCategory = "movie"
	Sports   EventCategory = "sports"
	Concert  EventCategory = "concert"
	Workshop EventCategory = "workshop"
	Party    EventCategory = "party"
)

type Event struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	HypeMeter   int           `json:"hype_meter"`
	Artists     []string      `json:"artists"`
	Duration    string        `json:"duration"`
	Category    EventCategory `json:"category"`
}

func (e *Event) AddHype() {
	//think of a way to accumulate logic--may be search how bms does it
}
