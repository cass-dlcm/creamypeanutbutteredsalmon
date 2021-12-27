package splatnet

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"os"
)

type shiftList struct {
	Code    *string `json:"code"`
	Summary struct {
		Card struct {
			GoldenIkuraTotal int `json:"golden_ikura_total"`
			HelpTotal        int `json:"help_total"`
			KumaPointTotal   int `json:"kuma_point_total"`
			IkuraTotal       int `json:"ikura_total"`
			KumaPoint        int `json:"kuma_point"`
			JobNum           int `json:"job_num"`
		} `json:"card"`
		Stats []struct {
			DeadTotal            int   `json:"dead_total"`
			MyGoldenIkuraTotal   int   `json:"my_golden_ikura_total"`
			GradePoint           int   `json:"grade_point"`
			TeamGoldenIkuraTotal int   `json:"team_golden_ikura_total"`
			HelpTotal            int   `json:"help_total"`
			TeamIkuraTotal       int   `json:"team_ikura_total"`
			StartTime            int   `json:"start_time"`
			MyIkuraTotal         int   `json:"my_ikura_total"`
			FailureCounts        []int `json:"failure_counts"`
			Schedule             struct {
				Stage struct {
					Image string `json:"image"`
					Name  string `json:"name"`
				} `json:"stage"`
				EndTime   int `json:"end_time"`
				StartTime int `json:"start_time"`
				Weapons   []struct {
					Weapon struct {
						ID        string `json:"id"`
						Image     string `json:"image"`
						Name      string `json:"name"`
						Thumbnail string `json:"thumbnail"`
					} `json:"weapon"`
					ID string `json:"id"`
				} `json:"weapons"`
			} `json:"schedule"`
			JobNum         int `json:"job_num"`
			KumaPointTotal int `json:"kuma_point_total"`
			EndTime        int `json:"end_time"`
			ClearNum       int `json:"clear_num"`
			Grade          struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"grade"`
		} `json:"stats"`
	} `json:"summary"`
	RewardGear struct {
		Thumbnail string `json:"thumbnail"`
		Kind      string `json:"kind"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		Brand     struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Image string `json:"image"`
		} `json:"brand"`
		Rarity int    `json:"rarity"`
		Image  string `json:"image"`
	} `json:"reward_gear"`
	Results shiftSplatnetResults `json:"results"`
}

type shiftSplatnetResults []shiftSplatnet

func (ssr *shiftSplatnetResults) inList(s shiftSplatnet) bool {
	for i := range *ssr {
		if (*ssr)[i].Equals(s) {
			return true
		}
	}
	return false
}

func (s *shiftSplatnet) Equals(s2 shiftSplatnet) bool {
	return s.PlayTime == s2.PlayTime
}

func (s *shiftSplatnet) GetTotalEggs() int {
	sum := 0
	for i := range s.WaveDetails {
		sum += s.WaveDetails[i].GoldenEggs
	}
	return sum
}

type shiftSplatnet struct {
	JobID           int64                   `json:"job_id"`
	DangerRate      float64                 `json:"danger_rate"`
	JobResult       shiftSplatnetJobResult  `json:"job_result"`
	JobScore        int                     `json:"job_score"`
	JobRate         int                     `json:"job_rate"`
	GradePoint      int                     `json:"grade_point"`
	GradePointDelta int                     `json:"grade_point_delta"`
	OtherResults    []shiftSplatnetPlayer   `json:"other_results"`
	KumaPoint       int                     `json:"kuma_point"`
	StartTime       int64                   `json:"start_time"`
	PlayerType      splatnetPlayerType      `json:"player_type"`
	PlayTime        int64                   `json:"play_time"`
	BossCounts      shiftSplatnetBossCounts `json:"boss_counts"`
	EndTime         int64                   `json:"end_time"`
	MyResult        shiftSplatnetPlayer     `json:"my_result"`
	WaveDetails     []shiftSplatnetWave     `json:"wave_details"`
	Grade           shiftSplatnetGrade      `json:"grade"`
	Schedule        shiftSplatnetSchedule   `json:"schedule"`
}

type shiftSplatnetJobResult struct {
	IsClear       bool    `json:"is_clear,omitempty"`
	FailureReason *string `json:"failure_reason,omitempty"`
	FailureWave   *int    `json:"failure_wave,omitempty"`
}

type shiftSplatnetPlayer struct {
	SpecialCounts  []int                           `json:"special_counts"`
	Special        splatnetQuad                    `json:"special"`
	Pid            string                          `json:"pid"`
	PlayerType     splatnetPlayerType              `json:"player_type"`
	WeaponList     []shiftSplatnetPlayerWeaponList `json:"weapon_list"`
	Name           string                          `json:"name"`
	DeadCount      int                             `json:"dead_count"`
	GoldenEggs     int                             `json:"golden_ikura_num"`
	BossKillCounts shiftSplatnetBossCounts         `json:"boss_kill_counts"`
	PowerEggs      int                             `json:"ikura_num"`
	HelpCount      int                             `json:"help_count"`
}

type shiftSplatnetPlayerWeaponList struct {
	ID     string                              `json:"id"`
	Weapon shiftSplatnetPlayerWeaponListWeapon `json:"weapon"`
}

type shiftSplatnetPlayerWeaponListWeapon struct {
	ID        string `json:"id"`
	Image     string `json:"image"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type shiftSplatnetBossCounts struct {
	Goldie    shiftSplatnetBossCountsBoss `json:"3"`
	Steelhead shiftSplatnetBossCountsBoss `json:"6"`
	Flyfish   shiftSplatnetBossCountsBoss `json:"9"`
	Scrapper  shiftSplatnetBossCountsBoss `json:"12"`
	SteelEel  shiftSplatnetBossCountsBoss `json:"13"`
	Stinger   shiftSplatnetBossCountsBoss `json:"14"`
	Maws      shiftSplatnetBossCountsBoss `json:"15"`
	Griller   shiftSplatnetBossCountsBoss `json:"16"`
	Drizzler  shiftSplatnetBossCountsBoss `json:"21"`
}

type shiftSplatnetBossCountsBoss struct {
	Boss  splatnetDouble `json:"boss"`
	Count int            `json:"count"`
}

type shiftSplatnetWave struct {
	WaterLevel   waterLevels `json:"water_level"`
	EventType    eventStruct `json:"event_type"`
	GoldenEggs   int         `json:"golden_ikura_num"`
	GoldenAppear int         `json:"golden_ikura_pop_num"`
	PowerEggs    int         `json:"ikura_num"`
	QuotaNum     int         `json:"quota_num"`
}

type shiftSplatnetGrade struct {
	ID        string `json:"id,omitempty"`
	ShortName string `json:"short_name,omitempty"`
	LongName  string `json:"long_name,omitempty"`
	Name      string `json:"name,omitempty"`
}

type shiftSplatnetSchedule struct {
	StartTime int64                         `json:"start_time"`
	Weapons   []shiftSplatnetScheduleWeapon `json:"weapons"`
	EndTime   int64                         `json:"end_time"`
	Stage     shiftSplatnetScheduleStage    `json:"stage"`
}

type shiftSplatnetScheduleWeapon struct {
	ID                string                                    `json:"id"`
	Weapon            *shiftSplatnetScheduleWeaponWeapon        `json:"weapon"`
	CoopSpecialWeapon *shiftSplatnetScheduleWeaponSpecialWeapon `json:"coop_special_weapon"`
}

type shiftSplatnetScheduleWeaponWeapon struct {
	ID        string `json:"id"`
	Image     string `json:"image"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type shiftSplatnetScheduleWeaponSpecialWeapon struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

type shiftSplatnetScheduleStage struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

type splatnetDouble struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type waterLevels struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type eventStruct struct {
	Key  event  `json:"key"`
	Name string `json:"name"`
}

type splatnetPlayerType struct {
	Gender  string `json:"style,omitempty"`
	Species string `json:"species,omitempty"`
}

type splatnetQuad struct {
	ID     string `json:"id"`
	ImageA string `json:"image_a"`
	ImageB string `json:"image_b"`
	Name   string `json:"name"`
}

const (
	smokeyard = "Salmonid Smokeyard"
	polaris   = "Ruins of Ark Polaris"
	grounds   = "Spawning Grounds"
	bay       = "Marooner's Bay"
	outpost   = "Lost Outpost"
)

type event string

const (
	griller          event = "griller"
	fog              event = "fog"
	cohockCharge     event = "cohock-charge"
	goldieSeeking    event = "goldie-seeking"
	mothership       event = "the-mothership"
	waterLevelsEvent event = "water-levels"
	rush             event = "rush"
)

func (e *event) ToEvent() types.Event {
	switch *e {
	case griller:
		return types.Griller
	case fog:
		return types.Fog
	case cohockCharge:
		return types.CohockCharge
	case goldieSeeking:
		return types.GoldieSeeking
	case mothership:
		return types.Mothership
	case waterLevelsEvent:
		return types.WaterLevels
	case rush:
		return types.Rush
	}
	return -1
}

func (e *event) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type E event
	r := (*E)(e)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	switch *e {
	case griller, fog, cohockCharge, goldieSeeking, mothership, waterLevelsEvent, rush:
		return nil
	}
	return errors.New("Invalid event. Got: " + fmt.Sprint(e))
}

const (
	ht = "high"
	lt = "low"
	nt = "normal"
)

func (s *shiftSplatnet) GetIdentifier(_ string) string {
	return fmt.Sprintf("https://app.splatoon2.nintendo.net/api/coop_results/%d", s.JobID)
}

type shiftSplatnetIterator struct {
	f          *os.File
	buffRead   *bufio.Scanner
	gzipReader *gzip.Reader
}

func (s *shiftSplatnetIterator) Next() (shift core.Shift, errs []error) {
	data := &shiftSplatnet{}
	if s.buffRead.Scan() {
		if err := json.Unmarshal([]byte(s.buffRead.Text()), &data); err != nil {
			errs = append(errs, err)
			if err := s.f.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := s.gzipReader.Close(); err != nil {
				errs = append(errs, err)
			}
			return nil, errs
		}
		if data == nil {
			errs = append(errs, &core.NoMoreShiftsError{})
			return nil, errs
		}
		return data, nil
	}
	if err := s.f.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.gzipReader.Close(); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, &core.NoMoreShiftsError{})
	return nil, errs
}

func (s *shiftSplatnetIterator) GetAddress() string {
	return ""
}

func (s *shiftSplatnet) GetClearWave() int {
	if s.JobResult.IsClear {
		return 3
	}
	return *s.JobResult.FailureWave - 1
}
