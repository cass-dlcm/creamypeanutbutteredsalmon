package salmonstats

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"log"
	"net/http"
	"os"
	"time"
)

/*
GetAllShifts downloads every shiftSalmonStats from the provided salmon-stats/api server and saves it to a gzipped jsonlines file.
*/
func GetAllShifts(db *sql.DB, dbType, userID string, server types.Server, client *http.Client, quiet bool) (errs []error) {
	schedule, errs2 := types.GetSchedules(client)
	if errs2 != nil {
		errs = append(errs, append(errs2, types.NewStackTrace())...)
		return errs
	}
	if !quiet {
		log.Println("Pulling Salmon Run data from online...")
	}
	getShifts := func(page int) (found bool, errs []error) {
		url := fmt.Sprintf("%splayers/%s/results", server.Address, userID)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errs = append(errs, err, types.NewStackTrace())
			return false, errs
		}
		query := req.URL.Query()
		query.Set("raw", "1")
		query.Set("count", "200")
		query.Set("page", fmt.Sprint(page))
		req.URL.RawQuery = query.Encode()

		if !quiet {
			log.Println(req.URL)
		}

		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, err, types.NewStackTrace())
			return false, errs
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				errs = append(errs, err)
			}
		}()
		var data shiftPage
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			errs = append(errs, err, types.NewStackTrace())
			return false, errs
		}
		var highestID int
		if err := db.QueryRow("SELECT id FROM Shifts ORDER BY id DESC LIMIT 1;").Scan(&highestID); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				errs = append(errs, err, types.NewStackTrace())
				return false, errs
			}
		}
		for i := range data.Results {
			events, errs2 := data.Results[i].GetEvents()
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return false, errs
			}
			tides, errs2 := data.Results[i].GetTides()
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return false, errs
			}
			golden := data.Results[i].GetEggsWaves()
			totalGolden := data.Results[i].GetTotalEggs()
			clearWave := data.Results[i].GetClearWave()
			princess := 0
			if data.Results[i].PlayerResults[0].GoldenEggs == totalGolden || (len(data.Results[i].PlayerResults) > 1 && (data.Results[i].PlayerResults[1].GoldenEggs == totalGolden || (len(data.Results[i].PlayerResults) > 2 && (data.Results[i].PlayerResults[2].GoldenEggs == totalGolden || (len(data.Results[i].PlayerResults) > 3 && data.Results[i].PlayerResults[3].GoldenEggs == totalGolden))))) {
				princess = 1
			}
			stage, errs2 := data.Results[i].GetStage(&schedule)
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return false, errs
			}
			weaponSet, errs2 := data.Results[i].GetWeaponSet(&schedule)
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return false, errs
			}
			t, errs2 := data.Results[i].GetTime()
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return false, errs
			}
			var w1event, w1tide, w2event, w2tide, w3event, w3tide string
			var w1g, w2g, w3g int
			if clearWave == 3 {
				w3event = (*events)[2].EventToKeyString()
				w3tide = string((*tides)[2])
				w3g = golden[2]
			}
			if clearWave >= 2 {
				w2event = (*events)[1].EventToKeyString()
				w2tide = string((*tides)[1])
				w2g = golden[2]
			}
			if clearWave >= 1 {
				w1event = (*events)[0].EventToKeyString()
				w1tide = string((*tides)[0])
				w1g = golden[0]
			}
			if dbType == "sqlite" {
				var id int
				if err := db.QueryRow("SELECT id FROM Shifts WHERE identifier = ?", data.Results[i].GetIdentifier(server)).Scan(&id); err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						errs = append(errs, err, types.NewStackTrace())
						return false, errs
					}
					if _, err := db.Exec("INSERT INTO Shifts VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);", data.Results[i].GetIdentifier(server), w1event, w1tide, w1g, w2event, w2tide, w2g, w3event, w3tide, w3g, clearWave, princess, stage, *weaponSet, t.Unix(), highestID+i+1); err != nil {
						errs = append(errs, err, types.NewStackTrace())
						return false, errs
					}
				}
			} else if dbType == "postgresql" {
				var id int
				if err := db.QueryRow("SELECT id FROM Shifts WHERE identifier = $1", data.Results[i].GetIdentifier(server)).Scan(&id); err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						errs = append(errs, err, types.NewStackTrace())
						return false, errs
					}
					if _, err := db.Exec("INSERT INTO Shifts VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);", data.Results[i].GetIdentifier(server), w1event, w1tide, w1g, w2event, w2tide, w2g, w3event, w3tide, w3g, clearWave, princess, stage, *weaponSet, t.Unix(), highestID+i+1); err != nil {
						errs = append(errs, err, types.NewStackTrace())
						return false, errs
					}
				}
			}
		}
		return len(data.Results) > 0, nil
	}
	if _, err := os.Stat(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName)); err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName))
			if err != nil {
				errs = append(errs, err, types.NewStackTrace())
				return errs
			}
			if err := f.Close(); err != nil {
				errs = append(errs, err, types.NewStackTrace())
				return errs
			}
		}
	}
	var count int
	if dbType == "sqlite" {
		_ = db.QueryRow("SELECT count() FROM Shifts WHERE identifier LIKE ?", fmt.Sprintf("%s%%", server.Address)).Scan(&count)
	} else if dbType == "postgresql" {
		_ = db.QueryRow("SELECT count() FROM Shifts WHERE identifier LIKE $1", fmt.Sprintf("%s%%", server.Address)).Scan(&count)
	}
	page := count/200 + 1
	hasPages := true
	for hasPages {
		hasPages, errs = getShifts(page)
		if len(errs) > 0 {
			return errs
		}
		page++
	}
	return nil
}

func (s *shiftSalmonStats) GetClearWave() int {
	return s.ClearWaves
}
