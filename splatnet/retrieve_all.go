package splatnet

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"github.com/cass-dlcm/splatnetiksm"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
GetAllShifts downloads every shiftSplatnet from the SplatNet server and saves it to the provided database.
*/
func GetAllShifts(db *sql.DB, dbType, sessionToken, cookie, locale, userID string, client *http.Client, quiet bool) (*string, *string, *string, []error) {
	var errs []error
	_, timezone := time.Now().Zone()
	timezone = -timezone / 60
	appHead := http.Header{
		"Host":              []string{"app.splatoon2.nintendo.net"},
		"x-unique-id":       []string{"32449507786579989235"},
		"x-requested-with":  []string{"XMLHttpRequest"},
		"x-timezone-offset": []string{fmt.Sprint(timezone)},
		"User-Agent":        []string{"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"Accept":            []string{"*/*"},
		"Referer":           []string{"https://app.splatoon2.nintendo.net/home"},
		"Accept-Encoding":   []string{"gzip deflate"},
		"Accept-Language":   []string{locale},
	}

	if !quiet {
		log.Println("Pulling Salmon Run data from online...")
	}

	url := "https://app.splatoon2.nintendo.net/api/coop_results"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errs = append(errs, err, types.NewStackTrace())
		return &sessionToken, &cookie, &userID, errs
	}

	req.Header = appHead

	if cookie == "" {
		newSessionToken, newCookie, errs2 := splatnetiksm.GenNewCookie(locale, sessionToken, "blank", client)
		if len(errs2) > 0 {
			errs = append(errs, errs2...)
			return &sessionToken, &cookie, &userID, errs
		}
		sessionToken = *newSessionToken
		cookie = *newCookie
	}

	req.AddCookie(&http.Cookie{Name: "iksm_session", Value: cookie})

	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, err, types.NewStackTrace())
		return &sessionToken, &cookie, &userID, errs
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			errs = append(errs, err)
		}
	}()

	var data shiftList

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		newSessionToken, newCookie, errs2 := splatnetiksm.GenNewCookie(locale, sessionToken, "blank", client)
		if len(errs2) > 0 {
			errs = append(errs, append(errs2, types.NewStackTrace())...)
			return &sessionToken, &cookie, &userID, errs
		}
		sessionToken = *newSessionToken
		cookie = *newCookie
		newSessionToken, newCookie, newID, errs2 := GetAllShifts(db, dbType, sessionToken, cookie, locale, userID, client, quiet)
		if len(errs2) > 0 {
			errs = append(errs, append(errs2, types.NewStackTrace())...)
			return &sessionToken, &cookie, &userID, errs
		}
		return newSessionToken, newCookie, newID, nil
	}

	if data.Code != nil {
		newSessionToken, newCookie, errs2 := splatnetiksm.GenNewCookie(locale, sessionToken, "auth", client)
		if len(errs2) > 0 {
			errs = append(errs, append(errs2, types.NewStackTrace())...)
			return &sessionToken, &cookie, &userID, errs
		}
		sessionToken = *newSessionToken
		cookie = *newCookie
		newSessionToken, newCookie, newID, errsRec := GetAllShifts(db, dbType, sessionToken, cookie, locale, userID, client, quiet)
		if len(errsRec) > 0 {
			errs = append(errs, errsRec...)
			return &sessionToken, &cookie, &userID, errs
		}
		return newSessionToken, newCookie, newID, nil
	}
	var highestID int
	if err := db.QueryRow("SELECT id FROM Shifts ORDER BY id DESC LIMIT 1;").Scan(&highestID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			errs = append(errs, err, types.NewStackTrace())
			return &sessionToken, &cookie, &userID, errs
		}
	}
	for i := range data.Results {
		userID = data.Results[i].MyResult.Pid
		result := &data.Results[i]
		var w1event, w1tide, w2event, w2tide, w3event, w3tide string
		var w1g, w2g, w3g int
		if result.GetClearWave() == 3 {
			w3event = strings.ReplaceAll(string(result.WaveDetails[2].EventType.Key), "-", "_")
			w3tide = result.WaveDetails[2].WaterLevel.Key
			w3g = result.WaveDetails[2].GoldenEggs
			if w3event == "the_mothership" {
				w3event = "mothership"
			}
		}
		if result.GetClearWave() >= 2 {
			w2event = strings.ReplaceAll(string(result.WaveDetails[1].EventType.Key), "-", "_")
			w2tide = result.WaveDetails[1].WaterLevel.Key
			w2g = result.WaveDetails[1].GoldenEggs
			if w2event == "the_mothership" {
				w2event = "mothership"
			}
		}
		if result.GetClearWave() >= 1 {
			w1event = strings.ReplaceAll(string(result.WaveDetails[0].EventType.Key), "-", "_")
			w1tide = result.WaveDetails[0].WaterLevel.Key
			w1g = result.WaveDetails[0].GoldenEggs
			if w1event == "the_mothership" {
				w1event = "mothership"
			}
		}
		totalGolden := 0
		for i := range result.WaveDetails {
			totalGolden += result.WaveDetails[i].GoldenEggs
		}
		princess := 0
		if totalGolden == result.MyResult.GoldenEggs || (len(result.OtherResults) > 0 && (totalGolden == result.OtherResults[0].GoldenEggs || (len(result.OtherResults) > 1 && (totalGolden == result.OtherResults[1].GoldenEggs || (len(result.OtherResults) > 2 && totalGolden == result.OtherResults[2].GoldenEggs))))) {
			princess = 1
		}
		stage, errs2 := result.GetStage(nil)
		if errs2 != nil {
			errs = append(errs, append(errs2, types.NewStackTrace())...)
			return &sessionToken, &cookie, &userID, errs
		}
		weaponSet, errs2 := result.GetWeaponSet(nil)
		if errs2 != nil {
			errs = append(errs, append(errs2, types.NewStackTrace())...)
			return &sessionToken, &cookie, &userID, errs
		}
		var id int
		if dbType == "sqlite" {
			if err := db.QueryRow("SELECT id FROM Shifts WHERE time = ?", result.PlayTime).Scan(&id); err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					errs = append(errs, err, types.NewStackTrace())
					return &sessionToken, &cookie, &userID, errs
				}
				if _, err := db.Exec("INSERT INTO Shifts VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);", fmt.Sprintf("https://app.splatoon2.nintendo.net/api/coop_results/%d", result.JobID), w1event, w1tide, w1g, w2event, w2tide, w2g, w3event, w3tide, w3g, result.GetClearWave(), princess, stage, *weaponSet, result.PlayTime, highestID+i+1); err != nil {
					errs = append(errs, err, types.NewStackTrace())
					return &sessionToken, &cookie, &userID, errs
				}
			}
		} else if dbType == "postgresql" {
			if err := db.QueryRow("SELECT id FROM Shifts WHERE identifier = $1", fmt.Sprintf("https://app.splatoon2.nintendo.net/api/coop_results/%d", result.JobID)).Scan(&id); err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					errs = append(errs, err, types.NewStackTrace())
					return &sessionToken, &cookie, &userID, errs
				}
				if _, err := db.Exec("INSERT INTO Shifts VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);", fmt.Sprintf("https://app.splatoon2.nintendo.net/api/coop_results/%d", result.JobID), w1event, w1tide, w1g, w2event, w2tide, w2g, w3event, w3tide, w3g, result.GetClearWave(), princess, stage, *weaponSet, result.PlayTime, highestID+i+1); err != nil {
					errs = append(errs, err, types.NewStackTrace())
					return &sessionToken, &cookie, &userID, errs
				}
			}
		}
	}
	return &sessionToken, &cookie, &userID, nil
}
