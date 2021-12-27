package types

// Tide is a string enum for denoting the water level of the wave.
type Tide string

// The three tides.
const (
	Ht Tide = "HT"
	Nt Tide = "NT"
	Lt Tide = "LT"
)

// TideArr is a wrapper around a Tide slice for the purpose of using the IsAllElementExist function.
type TideArr []Tide

// IsAllElementExist finds whether the given Tide slice contains every element in the TideArr.
func (t *TideArr) IsAllElementExist(arr []Tide) bool {
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

func (t *TideArr) HasElement(tide Tide) bool {
	for _, i := range *t {
		if i == tide {
			return true
		}
	}
	return false
}
