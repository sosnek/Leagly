package query

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var champ3 map[string]Champion

//!lastmatch player
func GetLastMatch(playerName string) (result string) {

	accInfo, exists := getAccountInfo(playerName)
	if exists {
		matchID, exist := getMatchID(accInfo.Puuid, 1)
		if exist {
			matchresults := getMatch(matchID[0])
			return formatLastMatchResponse(accInfo.Puuid, matchresults)
		}
		log.Println("Unable to get matchID for: " + playerName)
	}
	return "Sorry, something went wrong"
}

//!live player
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

//!lookup player
func LookupPlayer(playerName string) *discordgo.MessageSend {

	accInfo, exists := getAccountInfo(playerName)
	send := &discordgo.MessageSend{}
	if exists {
		rankedInfo := getRankedInfo(accInfo.Id)
		//masteryStats := getMasteryData(accInfo.Id)
		matchStatsID, exist := getMatchID(accInfo.Puuid, 10)

		if exist {
			var matchStatsSlice []MatchResults
			for n := 0; n < len(matchStatsID); n++ {
				matchStatsSlice = append(matchStatsSlice, getMatch(matchStatsID[n]))
			}

			//matchStatsFormatted := formatMatchStats(matchStatsSlice, accInfo.Puuid)
			//masteryStatsFormatted := formatMasteries(masteryStats)
			//fmt.Println(matchStatsFormatted)
			fileName, rankedType := getRankedAsset(RankedInfo(rankedInfo))
			file, _ := os.Open("./assets/" + fileName)
			embed := formatRankedEmbed(rankedInfo, rankedType, fileName)
			send = createMessageSend(embed, file, fileName)

			return send
		}
	}
	log.Println("Unable to get accInfo for: " + playerName)
	return send
}

func createMessageSend(embed *discordgo.MessageEmbed, file io.Reader, fileName string) *discordgo.MessageSend {
	var files []*discordgo.File
	files = append(files, &discordgo.File{
		Name:        fileName,
		ContentType: "image/png",
		Reader:      file,
	})

	send := &discordgo.MessageSend{
		Embed: embed,
		Files: files,
	}
	return send
}

func GetChampion(champID string) string {
	for k, v := range champ3 {
		if champID == v.Key {
			return k
		}
	}
	return champID
}

func formatMasteries(masteryStats Mastery) string {
	var iterations int
	if len(masteryStats) < 10 {
		iterations = len(masteryStats)
	} else {
		iterations = 5
	}
	masteryStatsFormatted := "```Champion Masteries: \n\n" + fmt.Sprintf("%-30s\t%-30s\t%-30s\t%-30s\n", "CHAMPION", "POINTS", "LEVEL", "LAST TIME CHAMP WAS PLAYED\n")
	for n := 0; n < iterations; n++ {
		masteryStatsFormatted = masteryStatsFormatted + fmt.Sprintf("%-30s\t%-30s\t%-30s\t%-30s", GetChampion(fmt.Sprint(masteryStats[n].ChampionID)),
			strconv.Itoa(masteryStats[n].ChampionPoints), strconv.Itoa(masteryStats[n].ChampionLevel),
			time.Unix(int64((masteryStats[n].LastPlayTime/1000)), 0).UTC().String()+"\n")
	}
	return masteryStatsFormatted + "```"
}

func formatRankedEmbed(rankedInfo RankedInfo, rankedType int, fileName string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:       000255000,
		Title:       rankedInfo[rankedType].SummonerName,
		Description: "Gold III with 68 LP. This season they have a total of 12 wins and 33 losses",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "attachment://" + fileName,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return embed
}

func formatMatchHistoryEmbed(embed *discordgo.MessageEmbed, matchedStats []MatchResults) {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    matchedStats[1].Info.Participants[1].SummonerName,
		IconURL: "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/profileicon/" + strconv.Itoa(matchedStats[1].Info.Participants[1].ProfileIcon) + ".png",
		URL:     "https://na.op.gg/summoner/userName=" + matchedStats[1].Info.Participants[1].SummonerName,
	}
}

//TODO:

//champion win/loss ratio of last 10 games
//champion % played
//champion total KDA of last 10 games
//total KDA of last 10 games
// win % of last 10 games
// role ratio
// using matchStats.Info.Queuetype, only use FLex and solo/duo games
func formatMatchStats(matchedStats []MatchResults, puuid string) string {
	var retString string
	//var playerStats []Participants
	var wins int
	var loss int

	for n := 0; n < len(matchedStats); n++ {
		participant := parseParticipant(puuid, matchedStats[n])
		if participant.Puuid == puuid {
			//playerStats[n].ChampionName = append(playerStats[n].ChampionName, participant.ChampionName)
			if participant.Win {
				wins++
			} else {
				loss++
			}
			break
		}
	}
	//after this is done we should have an array of type Participant that holds stats from the last 10 ranked games they played
	return retString
}

func formatPlayerRankedStats(rankedStats RankedInfo) (formattedRanked string) {
	for n := 0; n < len(rankedStats); n++ {
		if rankedStats[n].QueueType == "RANKED_SOLO_5x5" || rankedStats[n].QueueType == "RANKED_TEAM_5x5 " {
			if rankedStats[n].Tier == "" && rankedStats[n].Rank == "" {
				return "```" + rankedStats[n].SummonerName + " is currently unranked with " + strconv.Itoa(rankedStats[n].Wins) + " wins and " + strconv.Itoa(rankedStats[n].Losses) + " losses.```"
			}
			return "```" + rankedStats[n].SummonerName + " is currently " + rankedStats[n].Tier + " " + rankedStats[n].Rank +
				" and " + strconv.Itoa(rankedStats[n].LeaguePoints) + " LP. This season they have a total of " + strconv.Itoa(rankedStats[n].Wins) + " wins and " + strconv.Itoa(rankedStats[n].Losses) + " losses.```"
		} else {
			return "```" + rankedStats[n].SummonerName + " is not currently ranked.```"
		}

	}
	return "```No ranked data found```"
}

func formatLastMatchResponse(puuid string, matchResults MatchResults) (matchResultsFormatted string) {

	mySummonerStats := parseParticipant(puuid, matchResults)

	var hasWon string
	if mySummonerStats.Win {
		hasWon = "Yes```"
	} else {
		if mySummonerStats.GameEndedInSurrender {
			hasWon = "No, and it was an early surrender... Yikes```"
		} else {
			hasWon = "No```"
		}
	}

	minutes := matchResults.Info.GameDuration / 60
	seconds := matchResults.Info.GameDuration % 60

	resultsFormatted := "```" + mySummonerStats.SummonerName + "'s last game consists of the following stats:" +
		"\nDate: " + time.Unix(int64((matchResults.Info.GameCreation/1000)), 0).UTC().String() +
		"\nGame type: " + matchResults.Info.GameMode +
		"\nGame duration: " + fmt.Sprintf("%02d:%02d", int(minutes), int(seconds)) +
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

func getRankedAsset(rankedStats RankedInfo) (filename string, rankedType int) {
	//var filename string
	//var rankedType
	for n := 0; n < len(rankedStats); n++ {
		if rankedStats[n].QueueType == "RANKED_SOLO_5x5" || rankedStats[n].QueueType == "RANKED_TEAM_5x5 " {
			switch {
			case rankedStats[n].Tier == "IRON":
				return "Emblem_Iron.png", n
			case rankedStats[n].Tier == "BRONZE":
				return "Emblem_Bronze.png", n
			case rankedStats[n].Tier == "SILVER":
				return "Emblem_Silver.png", n
			case rankedStats[n].Tier == "GOLD":
				return "Emblem_Gold.png", n
			case rankedStats[n].Tier == "PLATINUM":
				return "Emblem_Platinum.png", n
			case rankedStats[n].Tier == "DIAMOND":
				return "Emblem_Diamond.png", n
			case rankedStats[n].Tier == "MASTER":
				return "Emblem_Master.png", n
			case rankedStats[n].Tier == "GRANDMASTER":
				return "Emblem_Grandmaster.png", n
			case rankedStats[n].Tier == "CHALLENGER":
				return "Emblem_Challenger.png", n
			}
		}
	}
	//return filename, rankedType
	return "", 0
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
