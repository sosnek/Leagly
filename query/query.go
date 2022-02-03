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
	TotalDamageTaken            int
	Item0                       int
	Item1                       int
	Item2                       int
	Item3                       int
	Item4                       int
	Item5                       int
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

///
///
///

func getMasteryData(accID string) Mastery {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Println("Unable to get mastery info. Error: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read mastery info. Error: " + err.Error())
	}

	masteryStats := string(body)

	var mastery Mastery
	json.Unmarshal([]byte(masteryStats), &mastery)

	return mastery
}

func getRankedInfo(accID string) RankedInfo {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Println("Unable to get ranked info. Error: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read ranked info. Error: " + err.Error())
	}

	rankedStats := string(body)

	var rankedInfo RankedInfo
	json.Unmarshal([]byte(rankedStats), &rankedInfo)

	return rankedInfo
}

func getLiveGame(summID string) LiveGameInfo {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + summID + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Println("Unable to get live game info. Error: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read live game info. Error: " + err.Error())
	}

	liveGame := string(body)

	var liveGameInfo LiveGameInfo
	json.Unmarshal([]byte(liveGame), &liveGameInfo)

	return liveGameInfo
}

func getMatch(matchid string) MatchResults {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/" + matchid + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Println("Unable to get match info. Error: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read match info. Error: " + err.Error())
	}

	match := string(body)

	var matchresults MatchResults
	json.Unmarshal([]byte(match), &matchresults)

	return matchresults
}

func getMatchID(puuid string, count int) []string {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/" + puuid + "/ids?start=0&count=" + (strconv.Itoa(count)) + "&api_key=" + config.ApiKey)
	var arr []string
	if err != nil {
		log.Println("Unable to get matchID data. Error: " + err.Error())
		return arr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read matchID data. Error: " + err.Error())
		return arr
	}

	_ = json.Unmarshal([]byte(body), &arr)
	if len(arr) == 0 {
		log.Println("Error unmarshaling MatchID data.")
		return arr
	}
	return arr
}

func getAccountInfo(playerName string) (summoner Summoner, exists bool) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + playerName + "?api_key=" + config.ApiKey)

	if err != nil {
		log.Println("Unable to get account info. Error: " + err.Error())
		return summoner, false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read account info. Error: " + err.Error())
		return summoner, false
	}
	//Convert the body to type string
	sb := string(body)

	var sum Summoner
	json.Unmarshal([]byte(sb), &sum)

	return sum, true

}

func InitializedChampStruct() {
	resp, err := http.Get("http://ddragon.leagueoflegends.com/cdn/12.1.1/data/en_US/champion.json")

	if err != nil {
		log.Println("Unable to get champion struct data. Error: " + err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read champion struct data. Error: " + err.Error())
		return
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
		log.Println("Unable to download file. Error: " + err.Error())
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println("Unable to download file. Status code: " + strconv.Itoa(response.StatusCode) + " Error: " + err.Error())
		return err
	}
	//Create a empty file
	file, err := os.Create("./championImages/" + fileName)
	if err != nil {
		log.Println("Error creating champion image directore or file. Error: " + err.Error())
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("Error copying image data to file. Error: " + err.Error())
		return err
	}

	return nil
}
