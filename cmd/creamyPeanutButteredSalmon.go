package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/salmonstats"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/splatnet"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/statink"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func setLanguage() []error {
	errs := []error{}
	log.Println("Please enter your locale (see readme for list).")

	var locale string
	// Taking input from user
	if _, err := fmt.Scanln(&locale); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return errs
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
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return errs
		}

		_, exists = languageList[locale]
	}
	viper.Set("user_lang", locale)

	if err := viper.WriteConfig(); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return errs
	}
	return nil
}

func getFlags() ([]types.Stage, []types.Event, []types.Tide, []types.WeaponSchedule, []types.Server, bool, []types.Server, string, []error) {
	errs := []error{}
	stagesStr := flag.String("stage", "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris", "To set a specific set of stages.")
	hasEventsStr := flag.String("event", "water_levels rush fog goldie_seeking griller cohock_charge mothership", "To set a specific set of events.")
	hasTides := flag.String("tide", "LT NT HT", "To set a specific set of tides.")
	hasWeapons := flag.String("weapon", "set single_random four_random random_gold", "To restrict to a specific set of weapon types.")
	statInk := flag.String("statink", "", "To read data from stat.ink. Use \"official\" for the server at stat.ink.")
	useSplatnet := flag.Bool("splatnet", false, "To read data from splatnet.")
	salmonStats := flag.String("salmonstats", "", "To read data from salmon-stats. Use \"official\" for the server at salmon-stats-api.yuki.games")
	outFile := flag.String("outfile", "", "To output data to a JSON file.")
	flag.Parse()

	stages := []types.Stage{}
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
			errs = append(errs, fmt.Errorf("stage not found: %s\n", stagesStrArr[i]))
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return nil, nil, nil, nil, nil, false, nil, "", errs
		}
		stages = append(stages, stageRes)
	}
	hasEvents := []types.Event{}
	eventsStrArr := strings.Split(*hasEventsStr, " ")
	for i := range eventsStrArr {
		eventRes, errs2 := types.StringToEvent(eventsStrArr[i])
		if errs2 != nil {
			errs = append(errs, errs2...)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return nil, nil, nil, nil, nil, false, nil, "", errs
		}
		hasEvents = append(hasEvents, *eventRes)
	}
	weapons := []types.WeaponSchedule{}
	weaponsStrArr := strings.Split(*hasWeapons, " ")
	for i := range weaponsStrArr {
		var weaponVal types.WeaponSchedule
		switch weaponsStrArr[i] {
		case string(types.RandommGrizzco), string(types.SingleRandom), string(types.FourRandom), string(types.Set):
			weaponVal = types.WeaponSchedule(weaponsStrArr[i])
		default:
			errs = append(errs, fmt.Errorf("weapon not found: %s\n", weaponsStrArr[i]))
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return nil, nil, nil, nil, nil, false, nil, "", errs
		}
		weapons = append(weapons, weaponVal)
	}

	tides := []types.Tide{}
	tidesStrArr := strings.Split(*hasTides, " ")
	for i := range tidesStrArr {
		inTide := types.Tide(tidesStrArr[i])
		switch inTide {
		case types.Ht, types.Lt, types.Nt:
			tides = append(tides, inTide)
		default:
			errs = append(errs, fmt.Errorf("tide not found: %s\n", tidesStrArr[i]))
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return nil, nil, nil, nil, nil, false, nil, "", errs
		}
	}

	statInkURLNicks := strings.Split(*statInk, " ")
	var statInkURLConf []types.Server
	if err := viper.UnmarshalKey("statink_servers", &statInkURLConf); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return nil, nil, nil, nil, nil, false, nil, "", errs
	}
	statInkServers := []types.Server{}
	for i := range statInkURLNicks {
		for j := range statInkURLConf {
			if statInkURLConf[j].ShortName == statInkURLNicks[i] {
				statInkServers = append(statInkServers, statInkURLConf[j])
			}
		}
	}

	salmonStatsURLNicks := strings.Split(*salmonStats, " ")
	var salmonStatsURLConf []types.Server
	if err := viper.UnmarshalKey("salmonstats_servers", &salmonStatsURLConf); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return nil, nil, nil, nil, nil, false, nil, "", errs
	}
	salmonStatsServers := []types.Server{}
	for i := range salmonStatsURLNicks {
		for j := range salmonStatsURLConf {
			if salmonStatsURLConf[j].ShortName == salmonStatsURLNicks[i] {
				salmonStatsServers = append(salmonStatsServers, salmonStatsURLConf[j])
			}
		}
	}

	return stages, hasEvents, tides, weapons, statInkServers, *useSplatnet, salmonStatsServers, *outFile, nil
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No config file found. One will be created.")
			viper.Set("cookie", "")
			viper.Set("session_token", "")
			viper.Set("user_lang", "")
			viper.Set("user_id", "")
			viper.Set("statink_servers", []types.Server{{
				ShortName: "official",
				APIKey:    "",
				Address:   "https://stat.ink/api/v2/",
			}})
			viper.Set("salmonstats_servers", []types.Server{{
				ShortName: "official",
				Address:   "https://salmon-stats-api.yuki.games/api/",
			}})
			if err := viper.WriteConfigAs("./config.json"); err != nil {
				log.Panicln(err)
			}
		} else {
			// Config file was found but another error was produced
			log.Panicf("Error reading the config file. Error is %v\n", err)
		}
	}
	viper.SetDefault("cookie", "")
	viper.SetDefault("session_token", "")
	viper.SetDefault("user_lang", "")
	viper.SetDefault("user_id", "")
	viper.SetDefault("statink_servers", []types.Server{{
		ShortName: "official",
		APIKey:    "",
		Address:   "https://stat.ink/api/v2/",
	}})
	viper.SetDefault("salmonstats_servers", []types.Server{{
		ShortName: "official",
		Address:   "https://salmon-stats-api.yuki.games/api/",
	}})
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	stages, hasEvents, tides, weapons, statInkServers, useSplatnet, salmonStatsServers, outFile, errs := getFlags()
	if len(errs) > 0 {
		log.Panicln(errs)
	}
	if errs := types.CheckForUpdate(client, outFile == ""); len(errs) > 0 {
		log.Panicln(errs)
	}
	if !(viper.IsSet("user_lang")) || viper.GetString("user_lang") == "" {
		setLanguage()
	}
	iterators := []core.ShiftIterator{}
	if useSplatnet {
		var sessionToken, cookie *string
		var errs []error
		sessionToken, cookie, errs = splatnet.GetAllShifts(viper.GetString("session_token"), viper.GetString("cookie"), viper.GetString("user_lang"), client, outFile == "")
		if errs != nil {
			for err := range errs {
				log.Println(errs[err])
			}
			log.Panicln(nil)
		}
		viper.Set("session_token", sessionToken)
		viper.Set("cookie", cookie)
		if err := viper.WriteConfig(); err != nil {
			log.Panicln(err)
		}
		iter, err := splatnet.LoadFromFileIterator()
		if err != nil {
			log.Panicln(err)
		}
		iterators = append(iterators, iter)
	}
	for i := range salmonStatsServers {
		if errs := salmonstats.GetAllShifts(salmonStatsServers[i], client, outFile == ""); len(errs) > 0 {
			log.Panicln(errs)
		}
		iter, errs := salmonstats.LoadFromFileIterator(salmonStatsServers[i])
		if len(errs) > 0 {
			log.Panicln(errs)
		}
		iterators = append(iterators, iter)
	}
	for i := range statInkServers {
		if errs := statink.GetAllShifts(statInkServers[i], client, outFile == ""); errs != nil {
			log.Panicln(errs)
		}
		iter, errs := statink.LoadFromFileIterator(statInkServers[i])
		if errs != nil {
			log.Panicln(errs)
		}
		iterators = append(iterators, iter)
	}
	records, errs := core.FindRecords(iterators, stages, hasEvents, tides, weapons, client)
	if errs != nil {
		log.Panicln(errs)
	}
	if outFile != "" {
		f, err := os.OpenFile("records.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
