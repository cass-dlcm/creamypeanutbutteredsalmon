package types

/*
Stage is an int enum of the stage for the rotation.
*/
type Stage int

/*
The five Salmon Run stages.
*/
const (
	SpawningGrounds Stage = iota
	MaroonersBay
	LostOutpost
	SalmonidSmokeyard
	RuinsOfArkPolaris
)

/*
String returns the name of the Stage, currently hardcoded as the en-US locale.
*/
func (s Stage) String() string {
	switch s {
	case SpawningGrounds:
		return "Spawning Grounds"
	case MaroonersBay:
		return "Marooner's Bay"
	case LostOutpost:
		return "Lost Outpost"
	case SalmonidSmokeyard:
		return "Salmonid Smokeyard"
	case RuinsOfArkPolaris:
		return "Ruins of Ark Polaris"
	}
	return ""
}

/*
StringJP returns the name of the Stage as the ja-JP locale.
*/
func (s Stage) StringJP() string {
	switch s {
	case SpawningGrounds:
		return "シェケナダム"
	case MaroonersBay:
		return "難破船ドン・ブラコ"
	case LostOutpost:
		return "海上集落シャケト場"
	case SalmonidSmokeyard:
		return "トキシラズいぶし工房"
	case RuinsOfArkPolaris:
		return "朽ちた箱舟ポラリス"
	}
	return ""
}

/*
IsElementExists finds whether the given Stage is in the Stage slice.
*/
func (s *Stage) IsElementExists(arr []Stage) bool {
	for _, v := range arr {
		if v == *s {
			return true
		}
	}
	return false
}
