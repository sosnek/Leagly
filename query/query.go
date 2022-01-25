package query

import (
	"Leagly/config"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Summoner struct {
	Id            string //summoner ID
	AccountId     string
	Puuid         string
	Name          string
	SummonerLevel int
	ProfileIconId int
}

type MatchResults struct {
	Info GameInfo
}

type GameInfo struct {
	GameDuration int
	GameMode     string
	GameCreation int64
	QueueId      int
	Participants []Participants
}

type Participants struct {
	Assists                     int
	ChampionName                string
	TotalDamageDealtToChampions int
	Deaths                      int
	IndividualPosition          string
	GameEndedInSurrender        bool
	Kills                       int
	TotalMinionsKilled          int
	VisionScore                 int
	Win                         bool
	Puuid                       string
	SummonerName                string
	ProfileIcon                 int
	//AP damage taken and AD damage taken ADD IT
}

type LiveGameInfo struct {
	GameStartTime int64
	GameMode      string
	GameType      string
	MapId         int
	Status        Status
}

type Status struct {
	Message     string
	Status_code int
}

type RankedInfo []struct {
	LeagueID     string
	QueueType    string
	Tier         string
	Rank         string
	SummonerName string
	LeaguePoints int
	Wins         int
	Losses       int
}

type Mastery []struct {
	ChampionID     int
	ChampionLevel  int
	ChampionPoints int
	LastPlayTime   int64
}

type Champion struct {
	Key string
}

const RANKED_SOLO = 420
const RANKED_FLEX = 440
const RATE_LIMIT = 30

///
///
///

func getMasteryData(accID string) Mastery {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	masteryStats := string(body)

	var mastery Mastery
	json.Unmarshal([]byte(masteryStats), &mastery)

	return mastery
}

func getRankedInfo(accID string) RankedInfo {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	rankedStats := string(body)

	var rankedInfo RankedInfo
	json.Unmarshal([]byte(rankedStats), &rankedInfo)

	return rankedInfo
}

func getLiveGame(summID string) LiveGameInfo {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + summID + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	liveGame := string(body)

	var liveGameInfo LiveGameInfo
	json.Unmarshal([]byte(liveGame), &liveGameInfo)

	return liveGameInfo
}

func getMatch(matchid string) MatchResults {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/" + matchid + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	match := string(body)

	var matchresults MatchResults
	json.Unmarshal([]byte(match), &matchresults)

	return matchresults
}

func getMatchID(puuid string, count int) ([]string, bool) {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/" + puuid + "/ids?start=0&count=" + (strconv.Itoa(count)) + "&api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var arr []string
	_ = json.Unmarshal([]byte(body), &arr)
	if len(arr) == 0 {
		return arr, false
	}
	return arr, true
}

func getAccountInfo(playerName string) (summoner Summoner, exists bool) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + playerName + "?api_key=" + config.ApiKey)

	if err != nil {
		log.Println(err)
		return summoner, false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return summoner, false
	}
	//Convert the body to type string
	sb := string(body)

	var sum Summoner
	json.Unmarshal([]byte(sb), &sum)

	return sum, true

}

//ToDo : create an enum to map every champion to its key and get the name
func InitializedChampStruct() {
	resp, err := http.Get("http://ddragon.leagueoflegends.com/cdn/12.1.1/data/en_US/champion.json")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)

	var objmap map[string]json.RawMessage
	json.Unmarshal([]byte(sb), &objmap)
	json.Unmarshal(objmap["data"], &champ3) //fuck you :)
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return err
	}
	//Create a empty file
	file, err := os.Create("./championImages/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
