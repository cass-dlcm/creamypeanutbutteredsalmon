package core

import (
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"time"
)

type recordName string

const (
	totalGoldenEggs         recordName = "Total Golden Eggs"
	totalGoldenEggsTwoNight recordName = "Total Golden Eggs (~2 Night)"
	totalGoldenEggsOneNight recordName = "Total Golden Eggs (~1 Night)"
	totalGoldenEggsNoNight  recordName = "Total Golden Eggs (No Night)"
	ntNormal                recordName = "NT Normal"
	htNormal                recordName = "HT Normal"
	ltNormal                recordName = "LT Normal"
	ntRush                  recordName = "NT Rush"
	htRush                  recordName = "HT Rush"
	ltRush                  recordName = "LT Rush"
	ntFog                   recordName = "NT Fog"
	htFog                   recordName = "HT Fog"
	ltFog                   recordName = "LT Fog"
	ntGoldieSeeking         recordName = "NT Goldie Seeking"
	htGoldieSeeking         recordName = "HT Goldie Seeking"
	ntGriller               recordName = "NT Griller"
	htGriller               recordName = "HT Griller"
	ntMothership            recordName = "NT Mothership"
	htMothership            recordName = "HT Mothership"
	ltMothershp             recordName = "LT Mothership"
	ltCohockCharge          recordName = "LT Cohock Charge"
)

// Shift is a generic source of match results, with only the necessary details available, accessible as methods.
type Shift interface {
	GetTotalEggs() int
	GetStage(types.Schedule) (*types.Stage, []error)
	GetWeaponSet(types.Schedule) (*types.WeaponSchedule, []error)
	GetEvents() (*types.EventArr, []error)
	GetTides() (*types.TideArr, []error)
	GetEggsWaves() []int
	GetWaveCount() int
	GetClearWave() int
	GetTime() (time.Time, []error)
	GetIdentifier(string) string
}

type record struct {
	Time         time.Time
	RecordAmount int
	Identifier   []string
}

func getRecordNames() []recordName {
	return []recordName{
		totalGoldenEggs,
		totalGoldenEggsTwoNight,
		totalGoldenEggsOneNight,
		totalGoldenEggsNoNight,
		ntNormal,
		htNormal,
		ltNormal,
		ntRush,
		htRush,
		ltRush,
		ntFog,
		htFog,
		ltFog,
		ntGoldieSeeking,
		htGoldieSeeking,
		ntGriller,
		htGriller,
		ntMothership,
		htMothership,
		ltMothershp,
		ltCohockCharge,
	}
}

func getAllRecords() map[recordName]*map[string]*map[types.WeaponSchedule]*record {
	records := map[recordName]*map[string]*map[types.WeaponSchedule]*record{}
	recordNames := getRecordNames()
	for i := range recordNames {
		records[recordNames[i]] = nil
	}
	return records
}

// ShiftIterator fulfils the design pattern of iterating through a set of Shift, only being able to progress one way.
// New in v4:
//  • GetAddress function
type ShiftIterator interface {
	Next() (Shift, []error)
	GetAddress() string
}

// NoMoreShiftsError implements the error interface.
// New in v4
type NoMoreShiftsError struct{}

// Error returns a static message of "no more shifts".
// New in v4
func (_ *NoMoreShiftsError) Error() string {
	return "no more shifts"
}
