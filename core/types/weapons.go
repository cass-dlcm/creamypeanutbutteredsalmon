package types

// WeaponSchedule consists of an indicator as to which weapons are in the shift.
type WeaponSchedule string

// The four possibilities for weapon set types.
const (
	RandommGrizzco WeaponSchedule = "random_gold"
	SingleRandom   WeaponSchedule = "single_random"
	FourRandom     WeaponSchedule = "four_random"
	Set            WeaponSchedule = "set"
)

// IsElementExists finds whether the given WeaponSchedule is in the WeaponSchedule slice.
func (w *WeaponSchedule) IsElementExists(arr []WeaponSchedule) bool {
	for _, v := range arr {
		if v == *w {
			return true
		}
	}
	return false
}
