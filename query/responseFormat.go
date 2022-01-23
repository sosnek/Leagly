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

type PlayerMatchStats struct {
	SummonerName    string
	ProfileIcon     int
	Role            Role
	PlayerChampions PlayerChampions
}

type Role struct {
	PreferredRole1 string
	PrefferedRole2 string

	PreferredRole1PickRate int
	PreferredRole2PickRate int

	PreferredRole1WinRate int
	PreferredRole2WinRate int
}

type PlayerChampions []struct {
	Name        string
	Wins        int
	Loss        int
	Kills       int
	Deaths      int
	Assists     int
	GamesPlayed int
}

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
		matchStatsID, exist := getMatchID(accInfo.Puuid, 20)

		if exist {
			var matchStatsSlice []MatchResults
			for n, k := 0, 0; n < len(matchStatsID) && k < 10; n++ {
				newMatch := getMatch(matchStatsID[n])
				if newMatch.Info.QueueId == RANKED_SOLO || newMatch.Info.QueueId == RANKED_FLEX {
					matchStatsSlice = append(matchStatsSlice, newMatch)
					k++
				}
			}

			playermatchstats := formatMatchStats(matchStatsSlice, accInfo.Puuid)
			//masteryStatsFormatted := formatMasteries(masteryStats)
			//fmt.Println(matchStatsFormatted)
			fileName, rankedType := getRankedAsset(RankedInfo(rankedInfo))
			file, _ := os.Open("./assets/" + fileName)
			description := formatPlayerRankedStats(rankedInfo)
			embed := formatRankedEmbed(rankedInfo, rankedType, fileName, description)
			embed = formatMatchHistoryEmbed(embed, playermatchstats)
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

func formatRankedEmbed(rankedInfo RankedInfo, rankedType int, fileName string, description string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:       000255000,
		Title:       rankedInfo[rankedType].SummonerName,
		Description: description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "attachment://" + fileName,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return embed
}

func formatMatchHistoryEmbed(embed *discordgo.MessageEmbed, playerMatchStats PlayerMatchStats) *discordgo.MessageEmbed {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    playerMatchStats.SummonerName,
		IconURL: "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/profileicon/" + strconv.Itoa(playerMatchStats.ProfileIcon) + ".png",
		URL:     "https://na.op.gg/summoner/userName=" + playerMatchStats.SummonerName,
	}
	return embed
}

//TODO:

//champion win/loss ratio of last 10 games
//champion % played
//champion total KDA of last 10 games
//total KDA of last 10 games
// win % of last 10 games
// role ratio
// using matchStats.Info.Queuetype, only use FLex and solo/duo games PlayerMatchStats
func formatMatchStats(matchedStats []MatchResults, puuid string) PlayerMatchStats {
	//var playerStats []Participants

	var playermatchstats PlayerMatchStats //Incorrect figure it out tomorrow @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

	for n := 0; n < len(matchedStats); n++ {
		participant := parseParticipant(puuid, matchedStats[n])
		playermatchstats.ProfileIcon = participant.ProfileIcon
		playermatchstats.SummonerName = participant.SummonerName
		playermatchstats.PlayerChampions[n].Name = participant.ChampionName
		playermatchstats.PlayerChampions[n].Kills = participant.Kills
		playermatchstats.PlayerChampions[n].Deaths = participant.Deaths
		playermatchstats.PlayerChampions[n].Assists = participant.Assists
		if participant.Win {
			playermatchstats.PlayerChampions[n].Wins++
		} else {
			playermatchstats.PlayerChampions[n].Loss++
		}
		playermatchstats.PlayerChampions[n].GamesPlayed++
		if participant.IndividualPosition == "TOP" {

		}

	}
	return playermatchstats
}

// func countChamps() {

// 	duplicate_frequency := make(map[string]int)

// 	for _, item := range list {
// 		// check if the item/element exist in the duplicate_frequency map

// 		_, exist := duplicate_frequency[item]

// 		if exist {
// 			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
// 		} else {
// 			duplicate_frequency[item] = 1 // else start counting from 1
// 		}
// 	}
// 	return duplicate_frequency
// }

func getRankedAsset(rankedStats RankedInfo) (filename string, rankedType int) {
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
	return "UNRANKED.png", 0
}

func parseParticipant(puuid string, matchresults MatchResults) Participants {
	var i int
	for n := 0; n < len(matchresults.Info.Participants); n++ {
		if puuid == matchresults.Info.Participants[n].Puuid {
			i = n
			break
		}
	}
	return matchresults.Info.Participants[i]
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

func formatPlayerRankedStats(rankedStats RankedInfo) string {
	for n := 0; n < len(rankedStats); n++ {
		if rankedStats[n].QueueType == "RANKED_SOLO_5x5" || rankedStats[n].QueueType == "RANKED_TEAM_5x5 " {
			return rankedStats[n].Tier + " " + rankedStats[n].Rank +
				" with " + strconv.Itoa(rankedStats[n].LeaguePoints) + " LP. This season they have a total of " + strconv.Itoa(rankedStats[n].Wins) + " wins and " + strconv.Itoa(rankedStats[n].Losses) + " losses.```"
		} else {
			return "Currently unranked."
		}

	}
	return "```No ranked data found```"
}
