package types

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

/*
Schedule is a struct of all the Salmon Run rotations since the game's launch.
*/
type Schedule struct {
	Result []ScheduleItem `json:"result"`
}

type ScheduleItem struct {
	Start    string    `json:"start"`
	StartUtc time.Time `json:"start_utc"`
	StartT   int       `json:"start_t"`
	End      string    `json:"end"`
	EndUtc   time.Time `json:"end_utc"`
	EndT     int       `json:"end_t"`
	Stage    struct {
		Image string `json:"image"`
		Name  string `json:"name"`
	} `json:"stage"`
	Weapons []struct {
		ID    int    `json:"id"`
		Image string `json:"image"`
		Name  string `json:"name"`
	} `json:"weapons"`
}

/*
GetSchedules downloads and returns a filled Schedule, or all the errors it encounters along with a stack trace.
*/
func GetSchedules(client *http.Client) (Schedule, []error) {
	var errs []error
	url := "https://spla2.yuu26.com/coop"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errs = append(errs, err, NewStackTrace())
		return Schedule{}, errs
	}
	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, err, NewStackTrace())
		return Schedule{}, errs
	}
	data := Schedule{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		errs = append(errs, err, NewStackTrace())
		return Schedule{}, errs
	}
	return data, nil
}
