package core

import (
	"database/sql"
	"errors"
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

/*
Shift is a generic source of match results, with only the necessary details available, accessible as methods.
*/
type Shift interface {
	GetTotalEggs() int
	GetStage(*types.Schedule) (*types.Stage, []error)
	GetWeaponSet(*types.Schedule) (*types.WeaponSchedule, []error)
	GetEvents() (*types.EventArr, []error)
	GetTides() (*types.TideArr, []error)
	GetEggsWaves() []int
	GetClearWave() int
	GetTime() (time.Time, []error)
	GetIdentifier() string
}

type DBShift struct {
	totalEggs  int
	stage      types.Stage
	weaponSet  types.WeaponSchedule
	events     types.EventArr
	tides      types.TideArr
	eggsWaves  []int
	princess   bool
	clearWave  int
	time       int
	identifier string
	id         int
}

func (D *DBShift) GetTotalEggs() int {
	sum := 0
	for i := 0; i < D.clearWave; i++ {
		sum += D.eggsWaves[i]
	}
	return sum
}

func (D *DBShift) GetStage(_ *types.Schedule) (*types.Stage, []error) {
	return &D.stage, nil
}

func (D *DBShift) GetWeaponSet(_ *types.Schedule) (*types.WeaponSchedule, []error) {
	return &D.weaponSet, nil
}

func (D *DBShift) GetEvents() (*types.EventArr, []error) {
	return &D.events, nil
}

func (D *DBShift) GetTides() (*types.TideArr, []error) {
	return &D.tides, nil
}

func (D *DBShift) GetEggsWaves() []int {
	return D.eggsWaves
}

func (D *DBShift) GetClearWave() int {
	return D.clearWave
}

func (D *DBShift) GetTime() (time.Time, []error) {
	return time.Unix(int64(D.time), 0), nil
}

func (D *DBShift) GetIdentifier() string {
	return D.identifier
}

type DBShiftIterator struct {
	db     *sql.DB
	dbType string
	id     int
}

func NewDBShiftIterator(db *sql.DB, dbType string) *DBShiftIterator {
	return &DBShiftIterator{db: db, dbType: dbType}
}

func (d *DBShiftIterator) Next() (Shift, []error) {
	eventStrs := []string{"", "", ""}
	shift := &DBShift{
		totalEggs:  0,
		stage:      0,
		weaponSet:  "",
		events:     types.EventArr{types.Event(-1), types.Event(-1), types.Event(-1)},
		tides:      types.TideArr{types.Tide(""), types.Tide(""), types.Tide("")},
		eggsWaves:  []int{-1, -1, -1},
		princess:   false,
		clearWave:  0,
		time:       0,
		identifier: "",
		id:         -1,
	}
	if d.dbType == "sqlite" {
		if err := d.db.QueryRow("SELECT * FROM Shifts WHERE id > ? ORDER BY id LIMIT 1;", d.id).Scan(&shift.identifier, &eventStrs[0], &shift.tides[0], &shift.eggsWaves[0], &eventStrs[1], &shift.tides[1], &shift.eggsWaves[1], &eventStrs[2], &shift.tides[2], &shift.eggsWaves[2], &shift.clearWave, &shift.princess, &shift.stage, &shift.weaponSet, &shift.time, &shift.id); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, []error{err, types.NewStackTrace()}
			}
			return nil, nil
		}
	} else if d.dbType == "postgresql" {
		if err := d.db.QueryRow("SELECT * FROM Shifts WHERE id > $1 ORDER BY id LIMIT 1;", d.id).Scan(&shift.identifier, &eventStrs[0], &shift.tides[0], &shift.eggsWaves[0], &eventStrs[1], &shift.tides[1], &shift.eggsWaves[1], &eventStrs[2], &shift.tides[2], &shift.eggsWaves[2], &shift.clearWave, &shift.princess, &shift.stage, &shift.weaponSet, &shift.time, &shift.id); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, []error{err, types.NewStackTrace()}
			}
			return nil, nil
		}
	}
	d.id++
	for i := range eventStrs {
		event := types.DisplayStringToEvent(eventStrs[i])
		if event == nil {
			return nil, []error{&types.ErrStrEventNotFound{Event: eventStrs[i]}, types.NewStackTrace()}
		}
		shift.events[i] = *event
	}
	if shift.clearWave < 3 {
		shift.eggsWaves = shift.eggsWaves[0:shift.clearWave]
		shift.events = shift.events[0:shift.clearWave]
		shift.tides = shift.tides[0:shift.clearWave]
	}
	return shift, nil
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

func getLatestRecords() map[recordName]*record {
	records := map[recordName]*record{}
	recordNames := getRecordNames()
	for i := range recordNames {
		records[recordNames[i]] = nil
	}
	return records
}

/*
ShiftIterator fulfils the design pattern of iterating through a set of Shift, only being able to progress one way.
*/
type ShiftIterator interface {
	Next() (Shift, []error)
}

/*
NoMoreShiftsError implements the error interface.
*/
type NoMoreShiftsError struct{}

/*
Error returns a static message of "no more shifts".
*/
func (*NoMoreShiftsError) Error() string {
	return "no more shifts"
}
