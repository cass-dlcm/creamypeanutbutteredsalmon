package types

/*
Event is an integer enum for denoting the event of a wave.
*/
type Event int

/*
The seven Salmon Run events.
*/
const (
	WaterLevels Event = iota
	Rush
	Fog
	GoldieSeeking
	Griller
	CohockCharge
	Mothership
)

/*
String returns the name of the Event, currently hardcoded as the en-US locale.
*/
func (e Event) String() string {
	switch e {
	case WaterLevels:
		return "Water Levels"
	case Rush:
		return "Rush"
	case Fog:
		return "Fog"
	case GoldieSeeking:
		return "Goldie Seeking"
	case Griller:
		return "Griller"
	case CohockCharge:
		return "Cohock Charge"
	case Mothership:
		return "Mothership"
	}
	return ""
}

func DisplayStringToEvent(inStr string) *Event {
	var event Event
	switch inStr {
	case WaterLevels.String():
		event = WaterLevels
	case Rush.String():
		event = Rush
	case Fog.String():
		event = Fog
	case GoldieSeeking.String():
		event = GoldieSeeking
	case Griller.String():
		event = Griller
	case CohockCharge.String():
		event = CohockCharge
	case Mothership.String():
		event = Mothership
	}
	return &event
}

/*
StringToEvent returns a pointer to an Event if the Event matches the inputted string, otherwise it returns an error.
*/
func StringToEvent(inStr string) Event {
	switch inStr {
	case "water_levels":
		return WaterLevels
	case "rush":
		return Rush
	case "fog":
		return Fog
	case "goldie_seeking":
		return GoldieSeeking
	case "griller":
		return Griller
	case "cohock_charge":
		return CohockCharge
	case "mothership":
		return Mothership
	}
	return -1
}

/*
EventToKeyString returns the key of an Event
*/
func (e Event) EventToKeyString() string {
	switch e {
	case WaterLevels:
		return "water_levels"
	case Rush:
		return "rush"
	case Fog:
		return "fog"
	case GoldieSeeking:
		return "goldie_seeking"
	case Griller:
		return "griller"
	case CohockCharge:
		return "cohock_charge"
	case Mothership:
		return "mothership"
	}
	return ""
}

/*
EventArr is a wrapper around an Event slice for the purpose of using the IsAllElementExist and HasElement functions.
*/
type EventArr []Event

/*
IsAllElementExist finds whether the given EventArr contains every element in the EventArr being called on.
*/
func (e *EventArr) IsAllElementExist(arr EventArr) bool {
	for _, i := range *e {
		found := false
		for _, j := range arr {
			if i == j {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

/*
HasElement finds whether the given Event is in the EventArr.
*/
func (e *EventArr) HasElement(event Event) bool {
	for _, i := range *e {
		if i == event {
			return true
		}
	}
	return false
}
