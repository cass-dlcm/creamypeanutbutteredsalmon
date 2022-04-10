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
String returns the name of the Event, currently hardcoded as the en-US locale, or an error if the Event isn't a valid value.
*/
func (e Event) String() (string, []error) {
	var errs []error
	switch e {
	case WaterLevels:
		return "Water Levels", nil
	case Rush:
		return "Rush", nil
	case Fog:
		return "Fog", nil
	case GoldieSeeking:
		return "Goldie Seeking", nil
	case Griller:
		return "Griller", nil
	case CohockCharge:
		return "Cohock Charge", nil
	case Mothership:
		return "Mothership", nil
	}
	errs = append(errs, &ErrIntEventNotFound{Event: int(e)})
	errs = append(errs, NewStackTrace())
	return "", errs
}

/*
StringToEvent returns a pointer to an Event if the Event matches the inputted string, otherwise it returns an error.
*/
func StringToEvent(inStr string) (*Event, []error) {
	var eventRes Event
	switch inStr {
	case "water_levels":
		eventRes = WaterLevels
	case "rush":
		eventRes = Rush
	case "fog":
		eventRes = Fog
	case "goldie_seeking":
		eventRes = GoldieSeeking
	case "griller":
		eventRes = Griller
	case "cohock_charge":
		eventRes = CohockCharge
	case "mothership":
		eventRes = Mothership
	default:
		errs := []error{&ErrStrEventNotFound{Event: inStr}}
		errs = append(errs, NewStackTrace())
		return nil, errs
	}
	return &eventRes, nil
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
