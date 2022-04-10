package statink

import (
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"time"
)

type shiftStatInk struct {
	ID             int           `json:"id"`
	UUID           string        `json:"uuid"`
	SplatnetNumber int           `json:"splatnet_number"`
	URL            string        `json:"url"`
	APIEndpoint    string        `json:"api_endpoint"`
	Stage          statinkTriple `json:"stage"`
	IsCleared      *bool         `json:"is_cleared"`
	FailReason     *struct {
		Key  string `json:"key"`
		Name name   `json:"name"`
	} `json:"fail_reason"`
	ClearWaves      *int        `json:"clear_waves"`
	DangerRate      interface{} `json:"danger_rate"`
	Quota           []int       `json:"quota"`
	Title           *title      `json:"title"`
	TitleExp        *int        `json:"title_exp"`
	TitleAfter      *title      `json:"title_after"`
	TitleExpAfter   *int        `json:"title_exp_after"`
	BossAppearances []bossCount `json:"boss_appearances"`
	Waves           []struct {
		KnownOccurrence      *statinkTriple `json:"known_occurrence"`
		WaterLevel           *statinkTriple `json:"water_level"`
		GoldenEggQuota       *int           `json:"golden_egg_quota"`
		GoldenEggAppearances *int           `json:"golden_egg_appearances"`
		GoldenEggDelivered   *int           `json:"golden_egg_delivered"`
		PowerEggCollected    *int           `json:"power_egg_collected"`
	} `json:"waves"`
	MyData    *player  `json:"my_data"`
	Teammates []player `json:"teammates"`
	Agent     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"agent"`
	Automated    bool         `json:"automated"`
	Note         *string      `json:"note"`
	LinkURL      *string      `json:"link_url"`
	ShiftStartAt *statInkTime `json:"shift_start_at"`
	StartAt      *statInkTime `json:"start_at"`
	EndAt        *statInkTime `json:"end_at"`
	RegisterAt   statInkTime  `json:"register_at"`
}

type statinkTriple struct {
	Key      string  `json:"key"`
	Splatnet *string `json:"splatnet"`
	Name     name    `json:"name"`
}

type title struct {
	Key         string `json:"key"`
	Splatnet    *int   `json:"splatnet"`
	Name        name   `json:"name"`
	GenericName name   `json:"generic_name"`
}

type statInkTime struct {
	Time    int64     `json:"time"`
	Iso8601 time.Time `json:"iso8601"`
}

type boss struct {
	Key         string `json:"key"`
	Splatnet    *int   `json:"splatnet"`
	SplatnetStr string `json:"splatnet_str"`
	Name        name   `json:"name"`
}

type bossCount struct {
	Boss  boss `json:"boss"`
	Count int  `json:"count"`
}

type player struct {
	SplatnetID string `json:"splatnet_id"`
	Name       string `json:"name"`
	Special    struct {
		Key  string `json:"key"`
		Name name   `json:"name"`
	} `json:"special"`
	Rescue             int `json:"rescue"`
	Death              int `json:"death"`
	GoldenEggDelivered int `json:"golden_egg_delivered"`
	PowerEggCollected  int `json:"power_egg_collected"`
	Species            struct {
		Key  string `json:"key"`
		Name name   `json:"name"`
	} `json:"species"`
	Gender struct {
		Key     string `json:"key"`
		Iso5218 int    `json:"iso5218"`
		Name    name   `json:"name"`
	} `json:"gender"`
	SpecialUses []int `json:"special_uses"`
	Weapons     []struct {
		Key      string `json:"key"`
		Splatnet int    `json:"splatnet"`
		Name     name   `json:"name"`
	} `json:"weapons"`
	BossKills []bossCount `json:"boss_kills"`
}

type name struct {
	DeDE string `json:"de_DE"`
	EnGB string `json:"en_GB"`
	EnUS string `json:"en_US"`
	EsES string `json:"es_ES"`
	EsMX string `json:"es_MX"`
	FrCA string `json:"fr_CA"`
	FrFR string `json:"fr_FR"`
	ItIT string `json:"it_IT"`
	JaJP string `json:"ja_JP"`
	NlNL string `json:"nl_NL"`
	RuRU string `json:"ru_RU"`
	ZhCN string `json:"zh_CN"`
	ZhTW string `json:"zh_TW"`
}

func (s *shiftStatInk) GetTotalEggs() int {
	sum := 0
	for i := range s.Waves {
		sum += *s.Waves[i].GoldenEggDelivered
	}
	return sum
}

func (s *shiftStatInk) GetStage(_ *types.Schedule) (*types.Stage, []error) {
	var stageRes types.Stage
	switch s.Stage.Key {
	case "dam":
		stageRes = types.SpawningGrounds
	case "donburako":
		stageRes = types.MaroonersBay
	case "polaris":
		stageRes = types.RuinsOfArkPolaris
	case "shaketoba":
		stageRes = types.LostOutpost
	case "tokishirazu":
		stageRes = types.SalmonidSmokeyard
	default:
		errs := []error{&types.ErrStrStageNotFound{s.Stage.Key}}
		errs = append(errs, types.NewStackTrace())
		return nil, errs
	}
	return &stageRes, nil
}

func (s *shiftStatInk) GetWeaponSet(weaponSets *types.Schedule) (*types.WeaponSchedule, []error) {
	var weaponRes types.WeaponSchedule
	for i := range weaponSets.Result {
		if weaponSets.Result[i].StartUtc.Equal(s.ShiftStartAt.Iso8601) {
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

func (s *shiftStatInk) GetEvents() (*types.EventArr, []error) {
	events := types.EventArr{}
	for i := range s.Waves {
		if s.Waves[i].KnownOccurrence == nil {
			events = append(events, types.WaterLevels)
		} else {
			event, errs := types.StringToEvent(s.Waves[i].KnownOccurrence.Key)
			if errs != nil {
				return nil, errs
			}
			events = append(events, *event)
		}
	}
	return &events, nil
}

func (s *shiftStatInk) GetTides() (*types.TideArr, []error) {
	tides := types.TideArr{}
	for i := range s.Waves {
		switch s.Waves[i].WaterLevel.Key {
		case "low":
			tides = append(tides, types.Lt)
		case "normal":
			tides = append(tides, types.Nt)
		case "high":
			tides = append(tides, types.Ht)
		default:
			errs := []error{&types.ErrStrTideNotFound{Tide: s.Waves[i].WaterLevel.Key}}
			errs = append(errs, types.NewStackTrace())
			return nil, errs
		}
	}
	return &tides, nil
}

func (s *shiftStatInk) GetEggsWaves() []int {
	eggs := []int{}
	for i := range s.Waves {
		eggs = append(eggs, *s.Waves[i].GoldenEggDelivered)
	}
	return eggs
}

func (s *shiftStatInk) GetWaveCount() int {
	return len(s.Waves)
}

func (s *shiftStatInk) GetTime() (time.Time, []error) {
	return s.StartAt.Iso8601.Local(), nil
}

func (s *shiftStatInk) GetIdentifier() string {
	return fmt.Sprintf("%d", s.ID)
}
