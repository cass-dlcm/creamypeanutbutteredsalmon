package core

import (
	"errors"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"net/http"
	"time"
)

func filterStages(stages []types.Stage, data Shift, schedules *types.Schedule) (Shift, []error) {
	stage, errs := data.GetStage(schedules)
	if errs != nil {
		return nil, errs
	}
	if stage.IsElementExists(stages) {
		return data, nil
	}
	return nil, nil
}

func filterEvents(events []types.Event, data Shift) (Shift, []error) {
	shiftEvents, errs := data.GetEvents()
	if len(errs) > 0 {
		return nil, errs
	}
	if shiftEvents.IsAllElementExist(events) {
		return data, nil
	}
	return nil, nil
}

func filterTides(tides []types.Tide, data Shift) (Shift, []error) {
	shiftTides, errs := data.GetTides()
	if errs != nil {
		return nil, errs
	}
	if shiftTides.IsAllElementExist(tides) {
		return data, nil
	}
	return nil, nil
}

func filterWeapons(weapons []types.WeaponSchedule, data Shift, schedules *types.Schedule) (Shift, []error) {
	weaponSet, errs := data.GetWeaponSet(schedules)
	if errs != nil {
		return nil, errs
	}
	if weaponSet.IsElementExists(weapons) {
		return data, nil
	}
	return nil, nil
}

type CompleteRecordsMap map[recordName]*map[string]*map[types.WeaponSchedule]*record

/*
FindRecords uses the given iterators to pull the shift data from the various sources based on the parameters, and finds records based on the set filters.
*/
func FindRecords(iterators []ShiftIterator, stages []types.Stage, hasEvents types.EventArr, tides types.TideArr, weapons []types.WeaponSchedule, client *http.Client) (CompleteRecordsMap, []error) {
	var errs []error
	scheduleList, errs2 := types.GetSchedules(client)
	if errs2 != nil {
		errs = append(errs, errs2...)
		return nil, errs
	}
	records := getAllRecords()
	for i := range iterators {
		shift, errs2 := iterators[i].Next()
		for errs2 == nil {
			if shift == nil {
				break
			}
			shift, errs2 = filterStages(stages, shift, &scheduleList)
			if errs2 != nil {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			shift, errs2 = filterEvents(hasEvents, shift)
			if len(errs2) > 0 {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			shift, errs2 = filterTides(tides, shift)
			if errs2 != nil {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			shift, errs2 = filterWeapons(weapons, shift, &scheduleList)
			if len(errs2) > 0 {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			totalEggs := shift.GetTotalEggs()
			weaponsType, _ := shift.GetWeaponSet(&scheduleList)
			stage, _ := shift.GetStage(&scheduleList)
			nightCount := 0
			waveEggs := shift.GetEggsWaves()
			waveEvents, _ := shift.GetEvents()
			waveWaterLevel, _ := shift.GetTides()
			clearWaves := shift.GetClearWave()
			var shiftTime time.Time
			shiftTime, errs2 = shift.GetTime()
			if errs2 != nil {
				errs = append(errs, errs2...)
				return nil, errs
			}
			for l := 0; l < clearWaves; l++ {
				if (*waveEvents)[l] == types.WaterLevels && hasEvents.HasElement(types.WaterLevels) {
					if records[recordName(string((*waveWaterLevel)[l])+" Normal")] == nil {
						records[recordName(string((*waveWaterLevel)[l])+" Normal")] = &map[string]*map[types.WeaponSchedule]*record{}
					}
					if (*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()] == nil {
						(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()] = &map[types.WeaponSchedule]*record{}
					}
					if (*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType] == nil || waveEggs[l] > (*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType].RecordAmount {
						(*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType] = &record{
							Time:         shiftTime,
							RecordAmount: waveEggs[l],
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if (*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType].Identifier = append((*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier())
					}
					continue
				}
				eventStr := (*waveEvents)[l].String()
				if eventStr == "" {
					errs = append(errs, &types.ErrIntEventNotFound{Event: int((*waveEvents)[l])}, types.NewStackTrace())
					return nil, errs
				}
				if hasEvents.HasElement((*waveEvents)[l]) &&
					tides.HasElement((*waveWaterLevel)[l]) {
					if records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)] == nil {
						records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)] = &map[string]*map[types.WeaponSchedule]*record{}
					}
					if (*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()] == nil {
						(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()] = &map[types.WeaponSchedule]*record{}
					}
					if (*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType] == nil || waveEggs[l] > (*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType].RecordAmount {
						(*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType] = &record{
							Time:         shiftTime,
							RecordAmount: waveEggs[l],
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if (*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType].Identifier = append((*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier())
					}
				}
				nightCount++
			}
			if clearWaves == 3 {
				if records[totalGoldenEggs] == nil {
					records[totalGoldenEggs] = &map[string]*map[types.WeaponSchedule]*record{}
				}
				if (*records[totalGoldenEggs])[stage.String()] == nil {
					(*records[totalGoldenEggs])[stage.String()] = &map[types.WeaponSchedule]*record{}
				}
				if (*(*records[totalGoldenEggs])[stage.String()])[*weaponsType] == nil || (*(*records[totalGoldenEggs])[stage.String()])[*weaponsType].RecordAmount < totalEggs {
					(*(*records[totalGoldenEggs])[stage.String()])[*weaponsType] = &record{
						Time:         shiftTime,
						RecordAmount: totalEggs,
						Identifier:   []string{shift.GetIdentifier()},
					}
				} else if (*(*records[totalGoldenEggs])[stage.String()])[*weaponsType].Time == shiftTime {
					(*(*records[totalGoldenEggs])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggs])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier())
				}
				if nightCount == 2 {
					if records[totalGoldenEggsTwoNight] == nil {
						records[totalGoldenEggsTwoNight] = &map[string]*map[types.WeaponSchedule]*record{}
					}
					if (*records[totalGoldenEggsTwoNight])[stage.String()] == nil {
						(*records[totalGoldenEggsTwoNight])[stage.String()] = &map[types.WeaponSchedule]*record{}
					}
					if (*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType] == nil || (*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType].RecordAmount < totalEggs {
						(*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType] = &record{
							Time:         shiftTime,
							RecordAmount: totalEggs,
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if (*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier())
					}
				}
				if nightCount == 1 {
					if records[totalGoldenEggsOneNight] == nil {
						records[totalGoldenEggsOneNight] = &map[string]*map[types.WeaponSchedule]*record{}
					}
					if (*records[totalGoldenEggsOneNight])[stage.String()] == nil {
						(*records[totalGoldenEggsOneNight])[stage.String()] = &map[types.WeaponSchedule]*record{}
					}
					if (*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType] == nil || (*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType].RecordAmount < totalEggs {
						(*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType] = &record{
							Time:         shiftTime,
							RecordAmount: totalEggs,
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if (*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier())
					}
				}
				if nightCount == 0 {
					if records[totalGoldenEggsNoNight] == nil {
						records[totalGoldenEggsNoNight] = &map[string]*map[types.WeaponSchedule]*record{}
					}
					if (*records[totalGoldenEggsNoNight])[stage.String()] == nil {
						(*records[totalGoldenEggsNoNight])[stage.String()] = &map[types.WeaponSchedule]*record{}
					}
					if (*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType] == nil || (*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType].RecordAmount < totalEggs {
						(*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType] = &record{
							Time:         shiftTime,
							RecordAmount: totalEggs,
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if (*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier())
					}
				}
			}
			shift, errs2 = iterators[i].Next()
		}
		if len(errs2) > 0 {
			if errors.Is(errs2[0], &NoMoreShiftsError{}) {
				continue
			}
			errs = append(errs, errs2...)
			return nil, errs
		}
	}
	return records, nil
}

func filterSchedule(shift Shift, latest *types.ScheduleItem, schedule *types.Schedule) (Shift, []error) {
	stage, errs := shift.GetStage(schedule)
	if errs != nil {
		return nil, errs
	}
	if stage.StringJP() != latest.Stage.Name {
		return nil, nil
	}
	weapons, errs := shift.GetWeaponSet(schedule)
	if errs != nil {
		return nil, errs
	}
	var weaponRes types.WeaponSchedule
	weaponSets := latest.Weapons
	if weaponSets[0].ID == -2 && weaponSets[1].ID == -2 && weaponSets[2].ID == -2 && weaponSets[3].ID == -2 {
		weaponRes = types.RandommGrizzco
	} else if weaponSets[0].ID == -1 && weaponSets[1].ID == -1 && weaponSets[2].ID == -1 && weaponSets[3].ID == -1 {
		weaponRes = types.FourRandom
	} else if weaponSets[0].ID >= 0 && weaponSets[1].ID >= 0 && weaponSets[2].ID >= 0 && weaponSets[3].ID == -1 {
		weaponRes = types.SingleRandom
	} else if weaponSets[0].ID >= 0 && weaponSets[1].ID >= 0 && weaponSets[2].ID >= 0 && weaponSets[3].ID >= 0 {
		weaponRes = types.Set
	} else {
		errs = []error{&types.ErrWeaponsNotFound{}}
		errs = append(errs, types.NewStackTrace())
		return nil, errs
	}
	if *weapons != weaponRes {
		return nil, nil
	}
	return shift, nil
}

type PartialRecordsMap map[recordName]*record

/*
FindLatest uses the given iterators to pull the shift data from the various sources based on the parameters, and finds records based on the set filters.
*/
func FindLatest(iterators []ShiftIterator, hasEvents types.EventArr, tides types.TideArr, client *http.Client) (PartialRecordsMap, []error) {
	var errs []error
	scheduleList, errs2 := types.GetSchedules(client)
	latest := &types.ScheduleItem{StartUtc: time.Unix(0, 0)}
	for i := range scheduleList.Result {
		if scheduleList.Result[i].StartUtc.Before(time.Now()) && scheduleList.Result[i].StartUtc.After(latest.StartUtc) {
			latest = &scheduleList.Result[i]
		}
	}
	if errs2 != nil {
		errs = append(errs, errs2...)
		return nil, errs
	}
	records := getLatestRecords()
	for i := range iterators {
		shift, errs2 := iterators[i].Next()
		for errs2 == nil {
			if shift == nil {
				break
			}
			shift, errs2 = filterSchedule(shift, latest, &scheduleList)
			if errs2 != nil {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			shift, errs2 = filterEvents(hasEvents, shift)
			if len(errs2) > 0 {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			shift, errs2 = filterTides(tides, shift)
			if errs2 != nil {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			totalEggs := shift.GetTotalEggs()
			nightCount := 0
			waveEggs := shift.GetEggsWaves()
			waveEvents, _ := shift.GetEvents()
			waveWaterLevel, _ := shift.GetTides()
			clearWaves := shift.GetClearWave()
			var shiftTime time.Time
			shiftTime, errs2 = shift.GetTime()
			if errs2 != nil {
				errs = append(errs, errs2...)
				return nil, errs
			}
			for l := 0; l < clearWaves; l++ {
				if (*waveEvents)[l] == types.WaterLevels && hasEvents.HasElement(types.WaterLevels) {
					if records[recordName(string((*waveWaterLevel)[l])+" Normal")] == nil || waveEggs[l] > records[recordName(string((*waveWaterLevel)[l])+" Normal")].RecordAmount {
						records[recordName(string((*waveWaterLevel)[l])+" Normal")] = &record{
							Time:         shiftTime,
							RecordAmount: waveEggs[l],
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if records[recordName(string((*waveWaterLevel)[l])+" Normal")].Time == shiftTime {
						records[recordName(string((*waveWaterLevel)[l])+" Normal")].Identifier = append(records[recordName(string((*waveWaterLevel)[l])+" Normal")].Identifier, shift.GetIdentifier())
					}
					continue
				}
				eventStr := (*waveEvents)[l].String()
				if eventStr == "" {
					errs = append(errs, &types.ErrIntEventNotFound{Event: int((*waveEvents)[l])}, types.NewStackTrace())
					return nil, errs
				}
				if hasEvents.HasElement((*waveEvents)[l]) &&
					tides.HasElement((*waveWaterLevel)[l]) {
					if records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)] == nil || waveEggs[l] > records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)].RecordAmount {
						records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)] = &record{
							Time:         shiftTime,
							RecordAmount: waveEggs[l],
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)].Time == shiftTime {
						records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)].Identifier = append(records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)].Identifier, shift.GetIdentifier())
					}
				}
				nightCount++
			}
			if clearWaves == 3 {
				if records[totalGoldenEggs] == nil || records[totalGoldenEggs].RecordAmount < totalEggs {
					records[totalGoldenEggs] = &record{
						Time:         shiftTime,
						RecordAmount: totalEggs,
						Identifier:   []string{shift.GetIdentifier()},
					}
				} else if records[totalGoldenEggs].Time == shiftTime {
					records[totalGoldenEggs].Identifier = append(records[totalGoldenEggs].Identifier, shift.GetIdentifier())
				}
				if nightCount == 2 {
					if records[totalGoldenEggsTwoNight] == nil || records[totalGoldenEggsTwoNight].RecordAmount < totalEggs {
						records[totalGoldenEggsTwoNight] = &record{
							Time:         shiftTime,
							RecordAmount: totalEggs,
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if records[totalGoldenEggsTwoNight].Time == shiftTime {
						records[totalGoldenEggsTwoNight].Identifier = append(records[totalGoldenEggsTwoNight].Identifier, shift.GetIdentifier())
					}
				}
				if nightCount == 1 {
					if records[totalGoldenEggsOneNight] == nil || records[totalGoldenEggsOneNight].RecordAmount < totalEggs {
						records[totalGoldenEggsOneNight] = &record{
							Time:         shiftTime,
							RecordAmount: totalEggs,
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if records[totalGoldenEggsOneNight].Time == shiftTime {
						records[totalGoldenEggsOneNight].Identifier = append(records[totalGoldenEggsOneNight].Identifier, shift.GetIdentifier())
					}
				}
				if nightCount == 0 {
					if records[totalGoldenEggsNoNight] == nil || records[totalGoldenEggsNoNight].RecordAmount < totalEggs {
						records[totalGoldenEggsNoNight] = &record{
							Time:         shiftTime,
							RecordAmount: totalEggs,
							Identifier:   []string{shift.GetIdentifier()},
						}
					} else if records[totalGoldenEggsNoNight].Time == shiftTime {
						records[totalGoldenEggsNoNight].Identifier = append(records[totalGoldenEggsNoNight].Identifier, shift.GetIdentifier())
					}
				}
			}
			shift, errs2 = iterators[i].Next()
		}
		if len(errs2) > 0 {
			if errors.Is(errs2[0], &NoMoreShiftsError{}) {
				continue
			}
			errs = append(errs, errs2...)
			return nil, errs
		}
	}
	return records, nil
}
