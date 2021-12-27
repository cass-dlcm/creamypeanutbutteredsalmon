package types

/*
Tide is a string enum for denoting the water level of the wave.
*/
type Tide string

/*
The three tides.
*/
const (
	Ht Tide = "HT"
	Nt Tide = "NT"
	Lt Tide = "LT"
)

/*
TideArr is a wrapper around a Tide slice for the purpose of using the IsAllElementExist and HasElement functions.
*/
type TideArr []Tide

/*
IsAllElementExist finds whether the given TideArr contains every element in the TideArr being called upon.
*/
func (t *TideArr) IsAllElementExist(arr TideArr) bool {
	for _, i := range *t {
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
HasElement finds whether the given Tide is in the TideArr.
*/
func (t *TideArr) HasElement(tide Tide) bool {
	for _, i := range *t {
		if i == tide {
			return true
		}
	}
	return false
}
