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

type RiotStatus struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Locales []string `json:"locales"`
	//Maintenances []interface{} `json:"maintenances"`
	Incidents []struct {
		Titles []struct {
			Content string `json:"content"`
			Locale  string `json:"locale"`
		} `json:"titles"`
		Updates []struct {
			Author       string `json:"author"`
			Translations []struct {
				Content string `json:"content"`
				Locale  string `json:"locale"`
			} `json:"translations"`
		} `json:"updates"`
	} `json:"incidents"`
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
	ChestGranted   bool
	TokensEarned   int
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

//Need to update this every major patch
const CHAMP_DATA = "http://ddragon.leagueoflegends.com/cdn/12.4.1/data/en_US/champion.json"
const CHAMPE_ICONS = "http://ddragon.leagueoflegends.com/cdn/12.4.1/img/champion/"

///
///
///https://na1.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-summoner/
func getMasteryData(accID string, regionPrefix string) Mastery {
	resp, err := http.Get("https://" + regionPrefix + ".api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
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

///https://na1.api.riotgames.com/lol/league/v4/entries/by-summoner/
func getRankedInfo(accID string, regionPrefix string) []*RankedInfo {
	resp, err := http.Get("https://" + regionPrefix + ".api.riotgames.com/lol/league/v4/entries/by-summoner/" + accID + "?api_key=" + config.ApiKey)
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

///https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/
func getLiveGame(summID string, regionPrefix string) LiveGameInfo {
	resp, err := http.Get("https://" + regionPrefix + ".api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + summID + "?api_key=" + config.ApiKey)
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

///https://americas.api.riotgames.com/lol/match/v5/matches/
func getMatch(matchid string, regionPrefix string) MatchResults {
	resp, err := http.Get("https://" + regionPrefix + ".api.riotgames.com/lol/match/v5/matches/" + matchid + "?api_key=" + config.ApiKey)
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

	json.Unmarshal([]byte(match), &matchresults)
	if matchresults == nil {
		log.Println("unmarshal error: error unmarshaling match data")
	}
	return matchresults
}

///https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/
func getMatchID(puuid string, count int, regionPrefix string) ([]string, error) {
	resp, err := http.Get("https://" + regionPrefix + ".api.riotgames.com/lol/match/v5/matches/by-puuid/" + puuid + "/ids?start=0&count=" + (strconv.Itoa(count)) + "&api_key=" + config.ApiKey)
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
		return arr, errors.New("Error unmarshaling MatchID data.")
	}
	return arr, err
}

///https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/
func getAccountInfo(playerName string, regionPrefix string) Summoner {
	resp, err := http.Get("https://" + regionPrefix + ".api.riotgames.com/lol/summoner/v4/summoners/by-name/" + playerName + "?api_key=" + config.ApiKey)
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

func getRiotStatus(regionPrefix string) RiotStatus {
	resp, err := http.Get("https://" + regionPrefix + ".api.riotgames.com/lol/status/v4/platform-data?api_key=" + config.ApiKey)
	var status RiotStatus
	if err != nil {
		log.Println("Unable to get riot status info. Error: " + err.Error())
		return status
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read riot status info. Error: " + err.Error())
		return status
	}
	//Convert the body to type string
	sb := string(body)
	json.Unmarshal([]byte(sb), &status)

	return status
}

func downloadFile(fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(CHAMPE_ICONS + fileName)
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

func InitializedChampStruct() error {
	resp, err := http.Get(CHAMP_DATA)

	if err != nil {
		return errors.New("Unable to get champion struct data. Error: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Unable to read champion struct data. Error: " + err.Error())
	}
	//Convert the body to type string
	sb := string(body)

	var objmap map[string]json.RawMessage
	json.Unmarshal([]byte(sb), &objmap)
	json.Unmarshal(objmap["data"], &champ3)
	return nil
}

func CreateChampionRatesFile() {
	resp, err := http.Get("https://cdn.merakianalytics.com/riot/lol/resources/latest/en-US/championrates.json")
	if err != nil {
		log.Println("Unable to get champion role data. Error: " + err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read champion role data. Error: " + err.Error())
		return
	}
	f, err := os.OpenFile("championRoleRates/championrates.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println("Error overwriting champion role rates file. Error:  " + err.Error())
		return
	}
	f.Write(body)
	if err := f.Close(); err != nil {
		log.Println("Error closing file. Error:  " + err.Error())
	}
}

//Merakianalytics stopped updating their data so for now I will keep a local version that will need manual updating.
func ChampionPositions() *map[string]ChampionRole {

	file, err := ioutil.ReadFile("championRoleRates/championrates.json")
	if err != nil {
		log.Println("Error reading champion role rates file. Error:  " + err.Error())
		return nil
	}
	//Convert the body to type string
	sb := string(file)

	var objmap map[string]json.RawMessage
	var champRole *map[string]ChampionRole
	json.Unmarshal([]byte(sb), &objmap)
	json.Unmarshal(objmap["data"], &champRole)

	return champRole
}
