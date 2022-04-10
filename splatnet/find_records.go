package splatnet

import (
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"time"
)

const (
	randomGrizzcoSchedule = "-2"
	randomSchedule        = "-1"
)

func (s *shiftSplatnet) GetWeaponSet(_ *types.Schedule) (*types.WeaponSchedule, []error) {
	var wS types.WeaponSchedule
	if s.Schedule.Weapons[0].ID == randomGrizzcoSchedule &&
		s.Schedule.Weapons[1].ID == randomGrizzcoSchedule &&
		s.Schedule.Weapons[2].ID == randomGrizzcoSchedule &&
		s.Schedule.Weapons[3].ID == randomGrizzcoSchedule {
		wS = types.RandommGrizzco
	}
	if s.Schedule.Weapons[0].ID == randomSchedule &&
		s.Schedule.Weapons[1].ID == randomSchedule &&
		s.Schedule.Weapons[2].ID == randomSchedule &&
		s.Schedule.Weapons[3].ID == randomSchedule {
		wS = types.FourRandom
	}
	if s.Schedule.Weapons[0].Weapon != nil &&
		s.Schedule.Weapons[1].Weapon != nil &&
		s.Schedule.Weapons[2].Weapon != nil &&
		s.Schedule.Weapons[3].ID == randomSchedule {
		wS = types.SingleRandom
	}
	if s.Schedule.Weapons[0].Weapon != nil &&
		s.Schedule.Weapons[1].Weapon != nil &&
		s.Schedule.Weapons[2].Weapon != nil &&
		s.Schedule.Weapons[3].Weapon != nil {
		wS = types.Set
	}
	if wS == "" {
		errs := []error{&ErrWeaponsNotFound{s.Schedule.Weapons[0].ID, s.Schedule.Weapons[1].ID, s.Schedule.Weapons[2].ID, s.Schedule.Weapons[3].ID}}
		errs = append(errs, types.NewStackTrace())
		return nil, errs
	}
	return &wS, nil
}

func (s *shiftSplatnet) GetEvents() (*types.EventArr, []error) {
	events := types.EventArr{}
	for i := range s.WaveDetails {
		var e types.Event
		switch s.WaveDetails[i].EventType.Key {
		case griller:
			e = types.Griller
		case fog:
			e = types.Fog
		case cohockCharge:
			e = types.CohockCharge
		case goldieSeeking:
			e = types.GoldieSeeking
		case mothership:
			e = types.Mothership
		case waterLevelsEvent:
			e = types.WaterLevels
		case rush:
			e = types.Rush
		default:
			errs := []error{&types.ErrStrEventNotFound{Event: string(s.WaveDetails[i].EventType.Key)}}
			errs = append(errs, types.NewStackTrace())
			return nil, errs
		}
		events = append(events, e)
	}
	return &events, nil
}

func (s *shiftSplatnet) GetStage(_ *types.Schedule) (*types.Stage, []error) {
	var stageResult types.Stage
	switch s.Schedule.Stage.Name {
	case polaris:
		stageResult = types.RuinsOfArkPolaris
	case outpost:
		stageResult = types.LostOutpost
	case bay:
		stageResult = types.MaroonersBay
	case smokeyard:
		stageResult = types.SalmonidSmokeyard
	case grounds:
		stageResult = types.SpawningGrounds
	default:
		errs := []error{&types.ErrStrStageNotFound{Stage: s.Schedule.Stage.Name}}
		errs = append(errs, types.NewStackTrace())
		return nil, errs
	}
	return &stageResult, nil
}

func (s *shiftSplatnet) GetTides() (*types.TideArr, []error) {
	tides := types.TideArr{}
	for i := range s.WaveDetails {
		var t types.Tide
		switch s.WaveDetails[i].WaterLevel.Key {
		case ht:
			t = types.Ht
		case lt:
			t = types.Lt
		case nt:
			t = types.Nt
		default:
			errs := []error{&types.ErrStrTideNotFound{Tide: s.WaveDetails[i].WaterLevel.Key}}
			errs = append(errs, types.NewStackTrace())
			return nil, errs
		}
		tides = append(tides, t)
	}
	return &tides, nil
}

func (s *shiftSplatnet) GetEggsWaves() []int {
	eggs := []int{}
	for i := range s.WaveDetails {
		eggs = append(eggs, s.WaveDetails[i].GoldenEggs)
	}
	return eggs
}

func (s *shiftSplatnet) GetWaveCount() int {
	return len(s.WaveDetails)
}

func (s *shiftSplatnet) GetTime() (time.Time, []error) {
	return time.Unix(s.PlayTime, 0).Local(), nil
}
