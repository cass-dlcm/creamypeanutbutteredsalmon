package core

import (
	"errors"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"net/http"
	"time"
)

func filterStages(stages []types.Stage, data Shift, schedules types.Schedule) (Shift, []error) {
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

func filterWeapons(weapons []types.WeaponSchedule, data Shift, schedules types.Schedule) (Shift, []error) {
	weaponSet, errs := data.GetWeaponSet(schedules)
	if errs != nil {
		return nil, errs
	}
	if weaponSet.IsElementExists(weapons) {
		return data, nil
	}
	return nil, nil
}

/*
FindRecords uses the given iterators to pull the shift data from the various sources based on the parameters, and finds records based on the set filters.
*/
func FindRecords(iterators []ShiftIterator, stages []types.Stage, hasEvents types.EventArr, tides types.TideArr, weapons []types.WeaponSchedule, client *http.Client) (map[recordName]*map[string]*map[types.WeaponSchedule]*record, []error) {
	var errs []error
	scheduleList, errs2 := types.GetSchedules(client)
	if errs2 != nil {
		errs = append(errs, errs2...)
		return nil, errs
	}
	records := getAllRecords()
	for i := range iterators {
		addr := iterators[i].GetAddress()
		shift, errs2 := iterators[i].Next()
		for errs2 == nil {
			if shift == nil {
				break
			}
			shift, errs2 = filterStages(stages, shift, scheduleList)
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
			shift, errs2 = filterWeapons(weapons, shift, scheduleList)
			if len(errs2) > 0 {
				errs = append(errs, errs2...)
				return nil, errs
			}
			if shift == nil {
				shift, errs2 = iterators[i].Next()
				continue
			}
			totalEggs := shift.GetTotalEggs()
			weaponsType, _ := shift.GetWeaponSet(scheduleList)
			stage, _ := shift.GetStage(scheduleList)
			nightCount := 0
			waveCount := shift.GetWaveCount()
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
			for l := 0; l < waveCount && i < clearWaves; l++ {
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
							Identifier:   []string{shift.GetIdentifier(addr)},
						}
					} else if (*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType].Identifier = append((*(*records[recordName(string((*waveWaterLevel)[l])+" Normal")])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier(addr))
					}
					continue
				}
				eventStr, errs2 := (*waveEvents)[l].String()
				if len(errs2) > 0 {
					errs = append(errs, errs2...)
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
							Identifier:   []string{shift.GetIdentifier(addr)},
						}
					} else if (*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType].Identifier = append((*(*records[recordName(string((*waveWaterLevel)[l])+" "+eventStr)])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier(addr))
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
						Identifier:   []string{shift.GetIdentifier(addr)},
					}
				} else if (*(*records[totalGoldenEggs])[stage.String()])[*weaponsType].Time == shiftTime {
					(*(*records[totalGoldenEggs])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggs])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier(addr))
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
							Identifier:   []string{shift.GetIdentifier(addr)},
						}
					} else if (*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggsTwoNight])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier(addr))
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
							Identifier:   []string{shift.GetIdentifier(addr)},
						}
					} else if (*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggsOneNight])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier(addr))
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
							Identifier:   []string{shift.GetIdentifier(addr)},
						}
					} else if (*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType].Time == shiftTime {
						(*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType].Identifier = append((*(*records[totalGoldenEggsNoNight])[stage.String()])[*weaponsType].Identifier, shift.GetIdentifier(addr))
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
