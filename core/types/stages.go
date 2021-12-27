package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// Stage is an int enum of the stage for the rotation.
type Stage int

// The five Salmon Run stages.
const (
	SpawningGrounds Stage = iota
	MaroonersBay
	LostOutpost
	SalmonidSmokeyard
	RuinsOfArkPolaris
)

// String returns the name of the Stage, currently hardcoded as the en-US locale.
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

// IsElementExists finds whether the given Stage is in the Stage slice.
func (s *Stage) IsElementExists(arr []Stage) bool {
	for _, v := range arr {
		if v == *s {
			return true
		}
	}
	return false
}

func (s Stage) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	jsonValue, err := json.Marshal(s.String())
	if err != nil {
		return nil, err
	}
	buffer.WriteString(string(jsonValue))
	return buffer.Bytes(), nil
}

func (s *Stage) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type S Stage
	r := (*S)(s)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	switch s.String() {
	case SpawningGrounds.String(), MaroonersBay.String(), LostOutpost.String(), SalmonidSmokeyard.String(), RuinsOfArkPolaris.String():
		return nil
	}
	return errors.New("Invalid StageENum. Got: " + fmt.Sprint(*s))
}
