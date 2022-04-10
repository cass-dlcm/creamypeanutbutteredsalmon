package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/salmonstats"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/splatnet"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/statink"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

func setLanguage() (string, []error) {
	var errs []error
	log.Println("Please enter your locale (see readme for list).")

	var locale string
	// Taking input from user
	if _, err := fmt.Scanln(&locale); err != nil {
		errs = append(errs, err, types.NewStackTrace())
		return "", errs
	}
	languageList := map[string]string{
		"en-US": "en-US",
		"es-MX": "es-MX",
		"fr-CA": "fr-CA",
		"ja-JP": "ja-JP",
		"en-GB": "en-GB",
		"es-ES": "es-ES",
		"fr-FR": "fr-FR",
		"de-DE": "de-DE",
		"it-IT": "it-IT",
		"nl-NL": "nl-NL",
		"ru-RU": "ru-RU",
	}
	_, exists := languageList[locale]
	for !exists {
		log.Println("Invalid language code. Please try entering it again.")

		if _, err := fmt.Scanln(&locale); err != nil {
			errs = append(errs, err, types.NewStackTrace())
			return "", errs
		}

		_, exists = languageList[locale]
	}
	return locale, nil
}

func getFlags(statInkURLConf []*types.Server, salmonStatsURLConf []*types.Server) (bool, []types.Stage, []types.Event, []types.Tide, []types.WeaponSchedule, []*types.Server, bool, []*types.Server, string, []error) {
	var errs []error
	mostRecentBool := flag.Bool("recent", false, "To calculate personal bests for the most recent rotation only.")
	stagesStr := flag.String("stage", "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris", "To set a specific set of stages.")
	hasEventsStr := flag.String("event", "water_levels rush fog goldie_seeking griller cohock_charge mothership", "To set a specific set of events.")
	hasTides := flag.String("tide", "LT NT HT", "To set a specific set of tides.")
	hasWeapons := flag.String("weapon", "set single_random four_random random_gold", "To restrict to a specific set of weapon types.")
	statInk := flag.String("statink", "", "To read data from stat.ink. Use \"official\" for the server at stat.ink.")
	useSplatnet := flag.Bool("splatnet", false, "To read data from splatnet.")
	salmonStats := flag.String("salmonstats", "", "To read data from salmon-stats. Use \"official\" for the server at salmon-stats-api.yuki.games")
	outFile := flag.String("outfile", "", "To output data to a JSON file.")
	flag.Parse()

	if *mostRecentBool && *stagesStr != "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris" {
		errs = append(errs, errors.New("incorrect flags; recent cannot be used with the stages flag"), types.NewStackTrace())
		return false, nil, nil, nil, nil, nil, false, nil, "", errs
	}
	if *mostRecentBool && *hasWeapons != "set single_random four_random random_gold" {
		errs = append(errs, errors.New("Incorrect flags; recent cannot be used with the weapons flag"), types.NewStackTrace())
		return false, nil, nil, nil, nil, nil, false, nil, "", errs
	}

	var stages []types.Stage
	stagesStrArr := strings.Split(*stagesStr, " ")
	for i := range stagesStrArr {
		var stageRes types.Stage
		switch stagesStrArr[i] {
		case "spawning_grounds":
			stageRes = types.SpawningGrounds
		case "marooners_bay":
			stageRes = types.MaroonersBay
		case "lost_outpost":
			stageRes = types.LostOutpost
		case "salmonid_smokeyard":
			stageRes = types.SalmonidSmokeyard
		case "ruins_of_ark_polaris":
			stageRes = types.RuinsOfArkPolaris
		default:
			errs = append(errs, &types.ErrStrStageNotFound{Stage: stagesStrArr[i]}, types.NewStackTrace())
			return false, nil, nil, nil, nil, nil, false, nil, "", errs
		}
		stages = append(stages, stageRes)
	}
	var hasEvents []types.Event
	eventsStrArr := strings.Split(*hasEventsStr, " ")
	for i := range eventsStrArr {
		eventRes, errs2 := types.StringToEvent(eventsStrArr[i])
		if errs2 != nil {
			errs = append(errs, append(errs2, types.NewStackTrace())...)
			return false, nil, nil, nil, nil, nil, false, nil, "", errs
		}
		hasEvents = append(hasEvents, *eventRes)
	}
	var weapons []types.WeaponSchedule
	weaponsStrArr := strings.Split(*hasWeapons, " ")
	for i := range weaponsStrArr {
		var weaponVal types.WeaponSchedule
		switch weaponsStrArr[i] {
		case string(types.RandommGrizzco), string(types.SingleRandom), string(types.FourRandom), string(types.Set):
			weaponVal = types.WeaponSchedule(weaponsStrArr[i])
		default:
			errs = append(errs, &types.ErrStrWeaponsNotFound{Weapons: weaponsStrArr[i]}, types.NewStackTrace())
			return false, nil, nil, nil, nil, nil, false, nil, "", errs
		}
		weapons = append(weapons, weaponVal)
	}

	var tides []types.Tide
	tidesStrArr := strings.Split(*hasTides, " ")
	for i := range tidesStrArr {
		inTide := types.Tide(tidesStrArr[i])
		switch inTide {
		case types.Ht, types.Lt, types.Nt:
			tides = append(tides, inTide)
		default:
			errs = append(errs, &types.ErrStrTideNotFound{Tide: tidesStrArr[i]}, types.NewStackTrace())
			return false, nil, nil, nil, nil, nil, false, nil, "", errs
		}
	}

	statInkURLNicks := strings.Split(*statInk, " ")
	var statInkServers []*types.Server
	for i := range statInkURLNicks {
		for j := range statInkURLConf {
			if statInkURLConf[j].ShortName == statInkURLNicks[i] {
				statInkServers = append(statInkServers, statInkURLConf[j])
			}
		}
	}

	salmonStatsURLNicks := strings.Split(*salmonStats, " ")
	var salmonStatsServers []*types.Server
	for i := range salmonStatsURLNicks {
		for j := range salmonStatsURLConf {
			if salmonStatsURLConf[j].ShortName == salmonStatsURLNicks[i] {
				salmonStatsServers = append(salmonStatsServers, salmonStatsURLConf[j])
			}
		}
	}

	return *mostRecentBool, stages, hasEvents, tides, weapons, statInkServers, *useSplatnet, salmonStatsServers, *outFile, nil
}

type config struct {
	Cookie             string          `json:"cookie"`
	SessionToken       string          `json:"session_token"`
	UserLang           string          `json:"user_lang"`
	UserId             string          `json:"user_id"`
	StatinkServers     []*types.Server `json:"statink_servers"`
	SalmonstatsServers []*types.Server `json:"salmonstats_servers"`
}

func newConfig() config {
	return config{
		StatinkServers: []*types.Server{
			{
				ShortName: "official",
				APIKey:    "",
				Address:   "https://stat.ink/api/v2/",
			},
		},
		SalmonstatsServers: []*types.Server{
			{
				ShortName: "official",
				Address:   "https://salmon-stats-api.yuki.games/api/",
			},
		},
	}
}

func main() {
	_, err := os.Stat("config.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Panicln(err)
	}
	if errors.Is(err, os.ErrNotExist) {
		configValues := newConfig()
		configJson, err := json.MarshalIndent(configValues, "", "    ")
		if err != nil {
			log.Panicln(err)
		}
		if err := os.WriteFile("config.json", configJson, 0600); err != nil {
			log.Panicln(err)
		}
	}
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Panicln(err)
	}
	configJson, err := ioutil.ReadAll(configFile)
	if err != nil {
		if err := configFile.Close(); err != nil {
			log.Println(err)
		}
		log.Panicln(err)
	}
	if err := configFile.Close(); err != nil {
		log.Panicln(err)
	}
	var configValues config
	if err := json.Unmarshal(configJson, &configValues); err != nil {
		log.Panicln(err)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	recent, stages, hasEvents, tides, weapons, statInkServers, useSplatnet, salmonStatsServers, outFile, errs := getFlags(configValues.StatinkServers, configValues.SalmonstatsServers)
	if len(errs) > 0 {
		log.Panicln(errs)
	}
	if errs := types.CheckForUpdate(client, outFile == ""); len(errs) > 0 {
		log.Panicln(errs)
	}
	if configValues.UserLang == "" {
		configValues.UserLang, errs = setLanguage()
		if errs != nil {
			log.Panicln(errs)
		}
		configJson, err = json.MarshalIndent(configValues, "", "    ")
		if err != nil {
			log.Panicln(err)
		}
		if err := os.WriteFile("config.json", configJson, 0600); err != nil {
			log.Panicln(err)
		}
	}
	dbType := "sqlite"
	db, err := sql.Open(dbType, "cpbs.sqlite")
	if err != nil {
		log.Panicln(err)
	}
	iterators := []core.ShiftIterator{core.NewDBShiftIterator(db, dbType)}
	if useSplatnet {
		sessionToken, cookie, userID, errs := splatnet.GetAllShifts(db, dbType, configValues.SessionToken, configValues.Cookie, configValues.UserLang, configValues.UserId, client, outFile == "")
		if errs != nil {
			log.Panicln(errs)
		}
		configValues.SessionToken, configValues.Cookie, configValues.UserId = *sessionToken, *cookie, *userID
		configJson, err = json.MarshalIndent(configValues, "", "    ")
		if err != nil {
			log.Panicln(err)
		}
		if err := os.WriteFile("config.json", configJson, 0600); err != nil {
			log.Panicln(err)
		}
	}
	for i := range salmonStatsServers {
		if errs := salmonstats.GetAllShifts(db, dbType, configValues.UserId, *salmonStatsServers[i], client, outFile == ""); len(errs) > 0 {
			log.Panicln(errs)
		}
	}
	for i := range statInkServers {
		if errs := statink.GetAllShifts(db, dbType, statInkServers[i], client, outFile == ""); errs != nil {
			log.Panicln(errs)
		}
		configJson, err = json.MarshalIndent(configValues, "", "    ")
		if err != nil {
			log.Panicln(err)
		}
		if err := os.WriteFile("config.json", configJson, 0600); err != nil {
			log.Panicln(err)
		}
	}
	if recent {
		records, errs := core.FindLatest(iterators, hasEvents, tides, client)
		if errs != nil {
			log.Panicln(errs)
		}
		if outFile != "" {
			f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				log.Panicln(err)
			}
			encoder := json.NewEncoder(f)
			encoder.SetIndent("", "    ")
			if err := encoder.Encode(records); err != nil {
				log.Panicln(err)
			}
			return
		}
		if err := json.NewEncoder(os.Stdout).Encode(records); err != nil {
			log.Panicln(err)
		}
		return
	}
	records, errs := core.FindRecords(iterators, stages, hasEvents, tides, weapons, client)
	if errs != nil {
		log.Panicln(errs)
	}
	if outFile != "" {
		f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Panicln(err)
		}
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(records); err != nil {
			log.Panicln(err)
		}
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(records); err != nil {
		log.Panicln(err)
	}
}
