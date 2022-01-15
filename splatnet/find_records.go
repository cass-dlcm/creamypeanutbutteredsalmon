package splatnet

import (
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"runtime"
	"time"
)

const (
	randomGrizzcoSchedule = "-2"
	randomSchedule        = "-1"
)

func (s *shiftSplatnet) GetWeaponSet(_ types.Schedule) (*types.WeaponSchedule, []error) {
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
		errs := []error{fmt.Errorf("WeaponSchedule not found: %s %s %s %s", s.Schedule.Weapons[0].ID, s.Schedule.Weapons[1].ID, s.Schedule.Weapons[2].ID, s.Schedule.Weapons[3].ID)}
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
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
			errs := []error{fmt.Errorf("no event found: %s", s.WaveDetails[i].EventType.Key)}
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return nil, errs
		}
		events = append(events, e)
	}
	return &events, nil
}

func (s *shiftSplatnet) GetStage(_ types.Schedule) (*types.Stage, []error) {
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
		errs := []error{fmt.Errorf("stage not found: %s", s.Schedule.Stage.Name)}
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
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
			errs := []error{fmt.Errorf("tide not found: %s", s.WaveDetails[i].WaterLevel.Key)}
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
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
