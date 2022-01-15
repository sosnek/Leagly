package query

import (
	"Leagly/config"
	"encoding/json"
	"fmt" //to print errors
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

type Summoner struct {
	Id            string
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

func GetLastMatch(playerName string) (result string) {
	accInfo := getAccountInfo(playerName)
	//error checking here?

	matchID, exists := getLastMatchID(accInfo.Puuid)
	if exists {
		matchresults := getMatch(matchID)
		lastMatchResultsFormatted := formatLastMatchResponse(accInfo.Puuid, matchresults)

		return lastMatchResultsFormatted
	} else {
		log.Println("Unable to get matchID for: " + playerName)
		return "Sorry, something went wrong"
	}

}

func IsInGame(playerName string) (result string) {

	accInfo := getAccountInfo(playerName)
	getLiveGame(accInfo.Id)
}

///
///
///

func getLiveGame(summID string) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + summID + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	liveGame := string(body)
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

func getAccountInfo(playerName string) Summoner {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + playerName + "?api_key=" + config.ApiKey)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)

	var summoner Summoner
	json.Unmarshal([]byte(sb), &summoner)

	return summoner
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
		"\nGame duration: " + strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds) +
		"\nChampion: " + mySummonerStats.ChampionName +
		"\nRole:" + mySummonerStats.IndividualPosition +
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

func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "commands:\n"
	msg = fmt.Sprintf("%s\t%s\n", msg, "!help - shows all available commands")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challengeBot - bot will always accept challenges")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challenges - shows all your open challenges")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challenge <username> - challenge another user")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!accepct <username> - accept a challenge from another user")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!leaderboard - shows the leaderboard")

	s.ChannelMessageSend(m.ChannelID, msg)
}
