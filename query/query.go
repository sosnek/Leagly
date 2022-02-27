package query

import (
	"Leagly/config"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Summoner *struct {
	Id            string //summoner ID
	AccountId     string
	Puuid         string
	Name          string
	SummonerLevel int
	ProfileIconId int
}

type MatchResults *struct {
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
	NeutralMinionsKilled        int
}

type LiveGameInfo *struct {
	GameStartTime     int64
	GameMode          string
	GameType          string
	MapId             int
	GameQueueConfigId int
	Participants      []LiveGameParticipants
	BannedChampions   []BannedChampions
	Status            Status
}

type LiveGameParticipants struct {
	ChampionId   int
	SummonerName string
	SummonerId   string
	Role         string
	Spell1Id     int
	Spell2Id     int
	TeamId       int
	championRole ChampionRole
}

type Status struct {
	Message     string
	Status_code int
}

type RankedInfo struct {
	LeagueID     string
	QueueType    string
	Tier         string
	Rank         string
	SummonerName string
	LeaguePoints int
	Wins         int
	Losses       int
}

type Mastery []*struct {
	ChampionID     int
	ChampionLevel  int
	ChampionPoints int
	LastPlayTime   int64
}

type PlayerMatchStats struct {
	SummonerName    string
	ProfileIcon     int
	Role            Role
	PlayerChampions []*PlayerChampions
}

type Role struct {
	PreferredRole1 string
	PreferredRole2 string

	RoleCount [5]int

	PreferredRole1WinRate int
	PreferredRole2WinRate int
}

type PlayerChampions struct {
	Name        string
	Wins        int
	Loss        int
	Kills       int
	Deaths      int
	Assists     int
	GamesPlayed int
}
type Champion struct {
	Key string
}

type ChampionRole struct {
	ID       int
	role     string
	Sum1     int
	Sum2     int
	BPH      float32
	skipRole []string
	Top      struct {
		PlayRate float32
	} `json:"TOP"`
	Jungle struct {
		PlayRate float32
	} `json:"JUNGLE"`
	Middle struct {
		PlayRate float32
	} `json:"MIDDLE"`
	Bottom struct {
		PlayRate float32
	} `json:"BOTTOM"`
	Utility struct {
		PlayRate float32
	} `json:"UTILITY"`
}

type BannedChampions struct {
	ChampionID int `json:"championId"`
	TeamID     int `json:"teamId"`
	PickTurn   int `json:"pickTurn"`
}

///
///
///

func getMasteryData(accID string) Mastery {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
	var mastery Mastery
	if err != nil {
		log.Println("Unable to get mastery info. Error: " + err.Error())
		return mastery
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read mastery info. Error: " + err.Error())
		return mastery
	}

	masteryStats := string(body)
	json.Unmarshal([]byte(masteryStats), &mastery)

	return mastery
}

func getRankedInfo(accID string) []*RankedInfo {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
	var rankedInfo []*RankedInfo
	if err != nil {
		log.Println("Unable to get ranked info. Error: " + err.Error())
		return rankedInfo
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read ranked info. Error: " + err.Error())
		return rankedInfo
	}

	rankedStats := string(body)

	json.Unmarshal([]byte(rankedStats), &rankedInfo)
	return rankedInfo
}

func getLiveGame(summID string) LiveGameInfo {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + summID + "?api_key=" + config.ApiKey)
	var liveGameInfo LiveGameInfo
	if err != nil {
		log.Println("Unable to get live game info. Error: " + err.Error())
		return liveGameInfo
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read live game info. Error: " + err.Error())
		return liveGameInfo
	}

	liveGame := string(body)

	json.Unmarshal([]byte(liveGame), &liveGameInfo)
	if liveGameInfo == nil {
		log.Println("Unmarshal error: unable to unmarshal live game data.")
	}

	return liveGameInfo
}

func getMatch(matchid string) MatchResults {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/" + matchid + "?api_key=" + config.ApiKey)
	var matchresults MatchResults
	if err != nil {
		log.Println("Unable to get match info. Error: " + err.Error())
		return matchresults
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read match info. Error: " + err.Error())
		return matchresults
	}

	match := string(body)

	err = json.Unmarshal([]byte(match), &matchresults)
	if matchresults == nil {
		log.Println("unmarshal error: error unmarshaling match data")
	}
	return matchresults
}

func getMatchID(puuid string, count int) ([]string, error) {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/" + puuid + "/ids?start=0&count=" + (strconv.Itoa(count)) + "&api_key=" + config.ApiKey)
	var arr []string
	if err != nil {
		log.Println("Unable to get matchID data. Error: " + err.Error())
		return arr, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read matchID data. Error: " + err.Error())
		return arr, err
	}

	_ = json.Unmarshal([]byte(body), &arr)
	if len(arr) == 0 {
		log.Println("Error unmarshaling MatchID data.")
		return arr, errors.New("Error unmarshaling MatchID data.")
	}
	return arr, err
}

func getAccountInfo(playerName string) Summoner {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + playerName + "?api_key=" + config.ApiKey)
	var sum Summoner
	if err != nil {
		log.Println("Unable to get account info. Error: " + err.Error())
		return sum
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read account info. Error: " + err.Error())
		return sum
	}
	//Convert the body to type string
	sb := string(body)

	json.Unmarshal([]byte(sb), &sum)

	return sum

}

func downloadFile(fileName string) error {
	//Get the response bytes from the url
	URL := "http://ddragon.leagueoflegends.com/cdn/12.4.1/img/champion/"
	response, err := http.Get(URL + fileName)
	if err != nil {
		log.Println("Unable to download file. Error: " + err.Error())
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println("Unable to download file. Status code: " + strconv.Itoa(response.StatusCode))
		return errors.New("http error: unable to download file")
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

func InitializedChampStruct() {
	resp, err := http.Get("http://ddragon.leagueoflegends.com/cdn/12.4.1/data/en_US/champion.json")

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

// func ChampionPositions() *map[string]ChampionRole {
// 	resp, err := http.Get("https://cdn.merakianalytics.com/riot/lol/resources/latest/en-US/championrates.json")
// 	if err != nil {
// 		log.Println("Unable to get champion role data. Error: " + err.Error())
// 		return nil
// 	}

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("Unable to read champion role data. Error: " + err.Error())
// 		return nil
// 	}
// 	//Convert the body to type string
// 	sb := string(body)

// 	var objmap map[string]json.RawMessage
// 	var champRole *map[string]ChampionRole
// 	json.Unmarshal([]byte(sb), &objmap)
// 	json.Unmarshal(objmap["data"], &champRole)

// 	return champRole
// }

//Merakianalytics stopped updating their data so for now I will keep a local version that will need manual updating.
func ChampionPositions() *map[string]ChampionRole {

	file, _ := ioutil.ReadFile("championRoleRates/championrates.json")
	//Convert the body to type string
	sb := string(file)

	var objmap map[string]json.RawMessage
	var champRole *map[string]ChampionRole
	json.Unmarshal([]byte(sb), &objmap)
	json.Unmarshal(objmap["data"], &champRole)

	return champRole
}
