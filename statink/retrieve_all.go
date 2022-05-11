package statink

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
GetAllShifts downloads every shiftStatInk from the provided stat.ink server and saves it to the provided database.
*/
func GetAllShifts(db *sql.DB, dbType string, statInkServer *types.Server, client *http.Client, quiet bool) (errs []error) {
	schedule, errs2 := types.GetSchedules(client)
	if errs2 != nil {
		errs = append(errs, append(errs2, types.NewStackTrace())...)
		return errs
	}
	if statInkServer.APIKey == "" {
		for len(statInkServer.APIKey) != 43 {
			log.Println("Please get your stat.ink API key, paste it here, and press enter: ")
			if _, err := fmt.Scanln(&statInkServer.APIKey); err != nil {
				errs = append(errs, err, types.NewStackTrace())
				return errs
			}
		}
	}
	getShift := func(id int) (data []shiftStatInk, errs []error) {
		url := fmt.Sprintf("%suser-salmon", statInkServer.Address)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errs = append(errs, err, types.NewStackTrace())
			return nil, errs
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", statInkServer.APIKey))
		query := req.URL.Query()
		query.Set("newer_than", fmt.Sprint(id))
		query.Set("order", "asc")
		req.URL.RawQuery = query.Encode()
		if !quiet {
			log.Println(req.URL)
		}
		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, err, types.NewStackTrace())
			return nil, errs
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				errs = append(errs, err)
			}
		}()
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Println(resp.Status)
			errs = append(errs, err, types.NewStackTrace())
			return nil, errs
		}
		var highestID int
		if err := db.QueryRow("SELECT id FROM Shifts ORDER BY id DESC LIMIT 1;").Scan(&highestID); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				errs = append(errs, err, types.NewStackTrace())
				return nil, errs
			}
		}
		for i := range data {
			events, errs2 := data[i].GetEvents()
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return nil, errs
			}
			tides, errs2 := data[i].GetTides()
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return nil, errs
			}
			golden := data[i].GetEggsWaves()
			totalGolden := data[i].GetTotalEggs()
			clearWave := data[i].GetClearWave()
			princess := 0
			if data[i].MyData.GoldenEggDelivered == totalGolden || (len(data[i].Teammates) > 0 && (data[i].Teammates[0].GoldenEggDelivered == totalGolden || (len(data[i].Teammates) > 1 && (data[i].Teammates[1].GoldenEggDelivered == totalGolden || (len(data[i].Teammates) > 2 && data[i].Teammates[2].GoldenEggDelivered == totalGolden))))) {
				princess = 1
			}
			stage, errs2 := data[i].GetStage(nil)
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return nil, errs
			}
			weaponSet, errs2 := data[i].GetWeaponSet(&schedule)
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return nil, errs
			}
			t, errs2 := data[i].GetTime()
			if errs2 != nil {
				errs = append(errs, append(errs2, types.NewStackTrace())...)
				return nil, errs
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
				if err := db.QueryRow("SELECT id FROM Shifts WHERE identifier = ?", data[i].GetIdentifier(*statInkServer)).Scan(&id); err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						errs = append(errs, err, types.NewStackTrace())
						return nil, errs
					}
					if _, err := db.Exec("INSERT INTO Shifts VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);", data[i].GetIdentifier(*statInkServer), w1event, w1tide, w1g, w2event, w2tide, w2g, w3event, w3tide, w3g, clearWave, princess, stage, *weaponSet, t.Unix(), highestID+i+1); err != nil {
						errs = append(errs, err, types.NewStackTrace())
						return nil, errs
					}
				}
			} else if dbType == "postgresql" {
				if err := db.QueryRow("SELECT id FROM Shifts WHERE identifier = $1", data[i].GetIdentifier(*statInkServer)).Scan(&id); err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						errs = append(errs, err, types.NewStackTrace())
						return nil, errs
					}
					if _, err := db.Exec("INSERT INTO Shifts VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);", data[i].GetIdentifier(*statInkServer), w1event, w1tide, w1g, w2event, w2tide, w2g, w3event, w3tide, w3g, clearWave, princess, stage, *weaponSet, t.Unix(), highestID+i+1); err != nil {
						errs = append(errs, err, types.NewStackTrace())
						return nil, errs
					}
				}
			}
		}
		log.Println(id)
		return data, nil
	}
	id := 1
	var ident string
	if err := db.QueryRow("SELECT identifier FROM Shifts WHERE identifier LIKE ? ORDER BY id DESC;", fmt.Sprintf("%s%%", statInkServer.Address)).Scan(&ident); err == nil {
		identStrs := strings.Split(ident, "/")
		idBig, err := strconv.ParseInt(identStrs[len(identStrs)-1], 10, 32)
		if err != nil {
			errs = append(errs, err, types.NewStackTrace())
			return errs
		}
		id = int(idBig)
	}
	for {
		tempData, errs2 := getShift(id)
		if len(errs2) > 0 {
			errs = append(errs, errs2...)
			return errs
		}
		if len(tempData) == 0 {
			return nil
		}
		id = tempData[len(tempData)-1].ID
	}
}

func (s *shiftStatInk) GetClearWave() int {
	return *s.ClearWaves
}
