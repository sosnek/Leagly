package query

import (
	"Leagly/config"
	"encoding/json"
	"fmt" //to print errors
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	//discordgo package from the repo of bwmarrin .
)

type Summoner struct {
	Id            string //summoner ID
	AccountId     string
	Puuid         string
	Name          string
	SummonerLevel int
}

type MatchResults struct {
	Metadata string
	Info     GameInfo
}

type GameInfo struct {
	GameDuration int
	GameMode     string
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
}

type LiveGameInfo struct {
	GameStartTime int64 //might be useless. TODO: use gamestart time instead
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

func GetLastMatch(playerName string) (result string) {
	accInfo, exists := getAccountInfo(playerName)
	//error checking here?

	if exists {
		matchID, exist := getLastMatchID(accInfo.Puuid)
		if exist {
			matchresults := getMatch(matchID)
			lastMatchResultsFormatted := formatLastMatchResponse(accInfo.Puuid, matchresults)
			return lastMatchResultsFormatted
		}
		log.Println("Unable to get matchID for: " + playerName)
	}
	return "Sorry, something went wrong"
}

func IsInGame(playerName string) (result string) {
	accInfo, exists := getAccountInfo(playerName)

	if exists {
		liveGameInfo := getLiveGame(accInfo.Id)

		if liveGameInfo.Status.Status_code == 0 {
			getTime := time.Now().UTC()
			elapsed := getTime.Sub(time.Unix(int64((liveGameInfo.GameStartTime / 1000)), 0).UTC())
			return fmt.Sprintf(playerName+" is currently in a game. Game time: %02d:%02d", (int(elapsed.Seconds()) / 60), (int(elapsed.Seconds()) % 60))
		}
		return playerName + " is not currently in-game."
	}
	return "Sorry, something went wrong"
}

func LookupPlayer(playerName string) (result string) {
	accInfo, exists := getAccountInfo(playerName)
	if exists {
		rankedStats := getRankedStats(accInfo.Id)
		return fmt.Sprintf(rankedStats[0].Tier + rankedStats[0].Rank)
	}

	log.Println("Unable to get accInfo for: " + playerName)
	return "Sorry, something went wrong"
}

///
///
///

func getRankedStats(accID string) RankedInfo {
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

func getLastMatchID(puuid string) (matchID string, exists bool) {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/" + puuid + "/ids?start=0&count=1&api_key=" + config.ApiKey)
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
		return " ", false
	}
	return arr[0], true
}

func getAccountInfo(playerName string) (summoner Summoner, exists bool) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + playerName + "?api_key=" + config.ApiKey)

	if err != nil {
		log.Fatalln(err)
		return summoner, false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return summoner, false
	}
	//Convert the body to type string
	sb := string(body)

	var sum Summoner
	json.Unmarshal([]byte(sb), &sum)

	return sum, true
}

func formatLastMatchResponse(puuid string, matchResults MatchResults) (matchResultsFormatted string) {

	mySummonerStats := parseParticipant(puuid, matchResults)

	var hasWon string
	if mySummonerStats.Win {
		hasWon = "Yes"
	} else {
		if mySummonerStats.GameEndedInSurrender {
			hasWon = "No, and it was an early surrender... Yikes"
		} else {
			hasWon = "No"
		}
	}

	minutes := matchResults.Info.GameDuration / 60
	seconds := matchResults.Info.GameDuration % 60

	resultsFormatted := mySummonerStats.SummonerName + "'s last game consists of the following stats:" +
		"\nGame type: " + matchResults.Info.GameMode +
		"\nGame duration: %02d" + strconv.Itoa(minutes) + "%02d:" + strconv.Itoa(seconds) +
		"\nChampion: " + mySummonerStats.ChampionName +
		"\nRole: " + mySummonerStats.IndividualPosition +
		"\nKills: " + strconv.Itoa(mySummonerStats.Kills) +
		"\nDeaths: " + strconv.Itoa(mySummonerStats.Deaths) +
		"\nAssists: " + strconv.Itoa(mySummonerStats.Assists) +
		"\nTotal DMG: " + strconv.Itoa(mySummonerStats.TotalDamageDealtToChampions) +
		"\nCS: " + strconv.Itoa(mySummonerStats.TotalMinionsKilled) +
		"\nVision score: " + strconv.Itoa(mySummonerStats.VisionScore) +
		"\nDid they win? " + hasWon

	return resultsFormatted
}

func parseParticipant(puuid string, matchresults MatchResults) Participants {
	var i int
	for n := 0; n < len(matchresults.Info.Participants); n++ {
		if puuid == matchresults.Info.Participants[n].Puuid {
			i = n
		}
	}
	return matchresults.Info.Participants[i]
}
