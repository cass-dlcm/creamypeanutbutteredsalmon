package main

import (
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
	"strings"
)

func setLanguage() {
	log.Println("Please enter your locale (see readme for list).")

	var locale string
	// Taking input from user
	if _, err := fmt.Scanln(&locale); err != nil {
		log.Panic(err)
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
			log.Panicln(err)
		}

		_, exists = languageList[locale]
	}
	viper.Set("user_lang", locale)

	if err := viper.WriteConfig(); err != nil {
		log.Panicln(err)
	}
}

func getFlags() ([]types.Stage, []types.Event, []types.Tide, []types.WeaponSchedule, []types.Server, bool, []types.Server) {
	stagesStr := flag.String("stage", "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris", "To set a specific set of stages.")
	hasEventsStr := flag.String("event", "water_levels rush fog goldie_seeking griller cohock_charge mothership", "To set a specific set of events.")
	hasTides := flag.String("tide", "LT NT HT", "To set a specific set of tides.")
	hasWeapons := flag.String("weapon", "set single_random four_random random_gold", "To restrict to a specific set of weapon types.")
	statInk := flag.String("statink", "", "To read data from stat.ink. Use \"official\" for the server at stat.ink.")
	useSplatnet := flag.Bool("splatnet", false, "To read data from splatnet.")
	salmonStats := flag.String("salmonstats", "", "To read data from salmon-stats. Use \"official\" for the server at salmon-stats-api.yuki.games")
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
			log.Panicf("stage not found: %s\n", stagesStrArr[i])
		}
		stages = append(stages, stageRes)
	}
	hasEvents := []types.Event{}
	eventsStrArr := strings.Split(*hasEventsStr, " ")
	for i := range eventsStrArr {
		eventRes, err := types.StringToEvent(eventsStrArr[i])
		if err != nil {
			log.Panicln(err)
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
			log.Panicf("weapon not found: %s\n", weaponsStrArr[i])
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
			log.Panicf("tide not found: %s\n", tidesStrArr[i])
		}
	}

	statInkURLNicks := strings.Split(*statInk, " ")
	var statInkURLConf []types.Server
	if err := viper.UnmarshalKey("statink_servers", &statInkURLConf); err != nil {
		log.Panicln(err)
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
		log.Panicln(err)
	}
	salmonStatsServers := []types.Server{}
	for i := range salmonStatsURLNicks {
		for j := range salmonStatsURLConf {
			if salmonStatsURLConf[j].ShortName == salmonStatsURLNicks[i] {
				salmonStatsServers = append(salmonStatsServers, salmonStatsURLConf[j])
			}
		}
	}

	return stages, hasEvents, tides, weapons, statInkServers, *useSplatnet, salmonStatsServers
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
	if errs := types.CheckForUpdate(client); len(errs) > 0 {
		for i := range errs {
			log.Println(errs[i])
		}
		log.Panicln(nil)
	}
	if !(viper.IsSet("user_lang")) || viper.GetString("user_lang") == "" {
		setLanguage()
	}
	stages, hasEvents, tides, weapons, statInkServers, useSplatnet, salmonStatsServers := getFlags()
	iterators := []core.ShiftIterator{}
	if useSplatnet {
		sessionToken, cookie, errs := splatnet.GetAllShifts(viper.GetString("session_token"), viper.GetString("cookie"), viper.GetString("user_lang"), client)
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
		errs := salmonstats.GetAllShifts(salmonStatsServers[i], client)
		if len(errs) > 0 {
			log.Panicln(errs)
		}
		iter, errs := salmonstats.LoadFromFileIterator(salmonStatsServers[i])
		if len(errs) > 0 {
			log.Panicln(errs)
		}
		iterators = append(iterators, iter)
	}
	for i := range statInkServers {
		errs := statink.GetAllShifts(statInkServers[i], client)
		if errs != nil {
			log.Panicln(errs)
		}
		iter, errs := statink.LoadFromFileIterator(statInkServers[i])
		if errs != nil {
			log.Panicln(errs)
		}
		iterators = append(iterators, iter)
	}
	f, err := os.OpenFile("records.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panicln(err)
	}
	if errs := core.FindRecords(iterators, stages, hasEvents, tides, weapons, client, f); errs != nil {
		log.Panicln(errs)
	}
}
