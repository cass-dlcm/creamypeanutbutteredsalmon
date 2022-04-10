package salmonstats

import (
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"time"
)

type shiftPage struct {
	CurrentPage  int    `json:"current_page"`
	FirstPageURL string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageURL  string `json:"last_page_url"`
	Links        []struct {
		URL    interface{} `json:"url"`
		Label  interface{} `json:"label"`
		Active bool        `json:"active"`
	} `json:"links"`
	NextPageURL string             `json:"next_page_url"`
	Path        string             `json:"path"`
	PerPage     int                `json:"per_page"`
	PrevPageURL interface{}        `json:"prev_page_url"`
	To          int                `json:"to"`
	Total       int                `json:"total"`
	Results     []shiftSalmonStats `json:"results"`
}

type bosses struct {
	Num3  int `json:"3"`
	Num6  int `json:"6"`
	Num9  int `json:"9"`
	Num12 int `json:"12"`
	Num13 int `json:"13"`
	Num14 int `json:"14"`
	Num15 int `json:"15"`
	Num16 int `json:"16"`
	Num21 int `json:"21"`
}

type shiftSalmonStats struct {
	ID                         int    `json:"id"`
	ScheduleID                 string `json:"schedule_id"`
	PlayerID                   string
	Page                       int
	StartAt                    string    `json:"start_at"`
	Members                    []string  `json:"members"`
	BossAppearances            bosses    `json:"boss_appearances"`
	UploaderUserID             int       `json:"uploader_user_id"`
	ClearWaves                 int       `json:"clear_waves"`
	FailReasonID               int       `json:"fail_reason_id"`
	DangerRate                 string    `json:"danger_rate"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	GoldenEggDelivered         int       `json:"golden_egg_delivered"`
	PowerEggCollected          int       `json:"power_egg_collected"`
	BossAppearanceCount        int       `json:"boss_appearance_count"`
	BossEliminationCount       int       `json:"boss_elimination_count"`
	IsEligibleForNoNightRecord bool      `json:"is_eligible_for_no_night_record"`
	PlayerResults              []struct {
		PlayerID             string      `json:"player_id"`
		GoldenEggs           int         `json:"golden_eggs"`
		PowerEggs            int         `json:"power_eggs"`
		Rescue               int         `json:"rescue"`
		Death                int         `json:"death"`
		SpecialID            int         `json:"special_id"`
		BossEliminationCount int         `json:"boss_elimination_count"`
		GradePoint           interface{} `json:"grade_point"`
		BossEliminations     struct {
			Counts bosses `json:"counts"`
		} `json:"boss_eliminations"`
		SpecialUses []struct {
			Count int `json:"count"`
		} `json:"special_uses"`
		Weapons []struct {
			WeaponID int `json:"weapon_id"`
		} `json:"weapons"`
	} `json:"player_results"`
	Waves []struct {
		Wave                 int `json:"wave"`
		EventID              int `json:"event_id"`
		WaterID              int `json:"water_id"`
		GoldenEggQuota       int `json:"golden_egg_quota"`
		GoldenEggAppearances int `json:"golden_egg_appearances"`
		GoldenEggDelivered   int `json:"golden_egg_delivered"`
		PowerEggCollected    int `json:"power_egg_collected"`
	} `json:"waves"`
}

func (s shiftSalmonStats) GetTotalEggs() int {
	return s.GoldenEggDelivered
}

func (s shiftSalmonStats) GetStage(schedule *types.Schedule) (*types.Stage, []error) {
	var stageRes types.Stage
	var errs []error
	scheduleTime, err := time.Parse("2006-01-02 15:04:05", s.ScheduleID)
	if err != nil {
		errs = []error{err}
		errs = append(errs, types.NewStackTrace())
		return nil, errs
	}
	for i := range schedule.Result {
		if schedule.Result[i].StartUtc.Equal(scheduleTime) {
			switch schedule.Result[i].Stage.Name {
			case "シェケナダム":
				stageRes = types.SpawningGrounds
			case "難破船ドン・ブラコ":
				stageRes = types.MaroonersBay
			case "トキシラズいぶし工房":
				stageRes = types.SalmonidSmokeyard
			case "朽ちた箱舟 ポラリス":
				stageRes = types.RuinsOfArkPolaris
			case "海上集落シャケト場":
				stageRes = types.LostOutpost
			}
			return &stageRes, nil
		}
	}
	errs = []error{}
	errs = append(errs, types.NewStackTrace())
	return nil, errs
}

func (s shiftSalmonStats) GetWeaponSet(weaponSets *types.Schedule) (*types.WeaponSchedule, []error) {
	var weaponRes types.WeaponSchedule
	scheduleTime, err := time.Parse("2006-01-02 15:04:05", s.ScheduleID)
	if err != nil {
		errs := []error{err}
		errs = append(errs, types.NewStackTrace())
		return nil, errs
	}
	for i := range weaponSets.Result {
		if weaponSets.Result[i].StartUtc.Equal(scheduleTime) {
			if weaponSets.Result[i].Weapons[0].ID == -2 && weaponSets.Result[i].Weapons[1].ID == -2 && weaponSets.Result[i].Weapons[2].ID == -2 && weaponSets.Result[i].Weapons[3].ID == -2 {
				weaponRes = types.RandommGrizzco
			}
			if weaponSets.Result[i].Weapons[0].ID == -1 && weaponSets.Result[i].Weapons[1].ID == -1 && weaponSets.Result[i].Weapons[2].ID == -1 && weaponSets.Result[i].Weapons[3].ID == -1 {
				weaponRes = types.FourRandom
			}
			if weaponSets.Result[i].Weapons[0].ID >= 0 && weaponSets.Result[i].Weapons[1].ID >= 0 && weaponSets.Result[i].Weapons[2].ID >= 0 && weaponSets.Result[i].Weapons[3].ID == -1 {
				weaponRes = types.SingleRandom
			}
			if weaponSets.Result[i].Weapons[0].ID >= 0 && weaponSets.Result[i].Weapons[1].ID >= 0 && weaponSets.Result[i].Weapons[2].ID >= 0 && weaponSets.Result[i].Weapons[3].ID >= 0 {
				weaponRes = types.Set
			}
			return &weaponRes, nil
		}
	}
	errs := []error{&types.ErrWeaponsNotFound{}}
	errs = append(errs, types.NewStackTrace())
	return nil, errs
}

func (s shiftSalmonStats) GetEvents() (*types.EventArr, []error) {
	events := types.EventArr{}
	for i := range s.Waves {
		switch s.Waves[i].EventID {
		case 0:
			events = append(events, types.WaterLevels)
		case 1:
			events = append(events, types.CohockCharge)
		case 2:
			events = append(events, types.Fog)
		case 3:
			events = append(events, types.GoldieSeeking)
		case 4:
			events = append(events, types.Griller)
		case 5:
			events = append(events, types.Mothership)
		case 6:
			events = append(events, types.Rush)
		default:
			errs := []error{&types.ErrIntEventNotFound{Event: s.Waves[i].EventID}}
			errs = append(errs, types.NewStackTrace())
			return nil, errs
		}
	}
	return &events, nil
}

func (s shiftSalmonStats) GetTides() (*types.TideArr, []error) {
	tides := types.TideArr{}
	for i := range s.Waves {
		switch s.Waves[i].WaterID {
		case 1:
			tides = append(tides, types.Lt)
		case 2:
			tides = append(tides, types.Nt)
		case 3:
			tides = append(tides, types.Ht)
		default:
			errs := []error{&types.ErrIntTideNotFound{Tide: s.Waves[i].WaterID}}
			errs = append(errs, types.NewStackTrace())
			return nil, errs
		}
	}
	return &tides, nil
}

func (s shiftSalmonStats) GetEggsWaves() []int {
	eggs := []int{}
	for i := range s.Waves {
		eggs = append(eggs, s.Waves[i].GoldenEggDelivered)
	}
	return eggs
}

func (s shiftSalmonStats) GetWaveCount() int {
	return len(s.Waves)
}

func (s shiftSalmonStats) GetTime() (time.Time, []error) {
	startTime, err := time.Parse("2006-01-02 15:04:05", s.StartAt)
	if err != nil {
		errs := []error{err}
		errs = append(errs, types.NewStackTrace())
		return time.Time{}, errs
	}
	return startTime.Local(), nil
}

func (s shiftSalmonStats) GetIdentifier() string {
	return fmt.Sprintf("results/%d/", s.ID)
}
