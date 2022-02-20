package query

import (
	"Leagly/config"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var emojis [][]*discordgo.Emoji
var champ3 map[string]Champion

//game codes
const URF = 900
const NORMAL = 400
const RANKED_SOLO = 420
const RANKED_FLEX = 440
const ARAM = 450

//Lookup match limit maxium last 30 matches
const MATCH_LIMIT = 30
const NUM_OF_RANK_GAMES = 10

///
///
///
func Help() *discordgo.MessageSend {
	embed := formatRankedEmbed("", "a", "Here is a list of the available commands for Leagly bot:", 16777215, time.Now())
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    "Leagly Bot",
		IconURL: "http://ddragon.leagueoflegends.com/cdn/12.4.1/img/profileicon/1630.png",
		URL:     "https://discord.com/oauth2/authorize?client_id=930924283599925260&permissions=1074056192&scope=bot",
	}
	embed = formatHelpEmbed(embed)
	return createMessageSend(embed, []*discordgo.File{})
}

//!live player
func IsInGame(playerName string) (send *discordgo.MessageSend, err error) {
	accInfo := getAccountInfo(playerName)
	send = &discordgo.MessageSend{}
	if accInfo != nil {
		liveGameInfo := getLiveGame(accInfo.Id)
		if liveGameInfo == nil {
			return send, errors.New("sorry, something went wrong")
		}
		if liveGameInfo.Status.Status_code == 0 {
			//var liveGameParticipants []LiveGameParticipants

			//liveGameParticipants = append(liveGameParticipants, liveGameInfo.Participants[2], liveGameInfo.Participants[0], liveGameInfo.Participants[4], liveGameInfo.Participants[3], liveGameInfo.Participants[1], liveGameInfo.Participants[7], liveGameInfo.Participants[5], liveGameInfo.Participants[9], liveGameInfo.Participants[8], liveGameInfo.Participants[6])
			//liveGameInfo.Participants = liveGameParticipants
			determineRoles(liveGameInfo.Participants)
			getTime := time.Now().UTC()
			elapsed := getTime.Sub(time.Unix(int64((liveGameInfo.GameStartTime / 1000)), 0).UTC())
			Gametime := fmt.Sprintf("%02d:%02d", (int(elapsed.Seconds()) / 60), (int(elapsed.Seconds()) % 60))

			participant := parseLiveParticipant(accInfo.Id, liveGameInfo)
			rankPlayers := formatRankedPlayers(liveGameInfo)
			//get bans as well
			//bannedChampions := getBannedChampsID(liveGameInfo)
			champion := GetChampion(strconv.Itoa(participant.ChampionId))
			err = getChampionFile(champion + ".png")
			if err != nil {
				return send, errors.New("Error getting champion file for" + champion)
			}
			//Playing as to banned champs
			embed := formatRankedEmbed(playerName+" Is currently in a "+getMatchType(liveGameInfo.GameQueueConfigId), champion+".png", "Playing as "+champion+". Time: "+Gametime, 71, time.Now())
			embed = formatEmbedAuthor(embed, accInfo)
			files := formatEmbedImages([]string{}, "./championImages/", champion+".png")
			embed = formatLiveMatchEmbedFields(embed, rankPlayers, liveGameInfo, participant)
			send = createMessageSend(embed, files)
			return send, nil
		}
		embed := &discordgo.MessageEmbed{
			Color:       3080243,
			Title:       playerName,
			Description: "Not currently in-game",
			Timestamp:   time.Now().Format(time.RFC3339),
		}
		embed = formatEmbedAuthor(embed, accInfo)
		send = createMessageSend(embed, []*discordgo.File{})
		return send, nil
	}
	return send, errors.New("sorry, something went wrong")
}

//!lastmatch player
func GetLastMatch(playerName string) (send *discordgo.MessageSend, err error) {
	accInfo := getAccountInfo(playerName)
	send = &discordgo.MessageSend{}
	if accInfo != nil {
		matchID, err := getMatchID(accInfo.Puuid, 1)
		if err != nil {
			return send, errors.New("Error getting match results for " + playerName)
		}
		if len(matchID) < 1 {
			return send, errors.New("No match history found for " + playerName)
		}
		matchresults := getMatch(matchID[0])
		if matchresults == nil {
			return send, errors.New("Error getting match results for " + playerName)
		}
		participant := parseParticipant(accInfo.Puuid, matchresults)
		fileName := participant.ChampionName + ".png"
		err = getChampionFile(fileName)
		if err != nil {
			return send, err
		}
		embed := formatRankedEmbed(getMatchType(matchresults.Info.QueueId)+". Time: "+fmt.Sprintf("%02d:%02d", int(matchresults.Info.GameDuration/60), int(matchresults.Info.GameDuration%60)), fileName, formatItems(participant), getEmbedColour(participant.Win), time.Unix(int64((matchresults.Info.GameCreation)/1000) + +int64(matchresults.Info.GameDuration), 0).Local())
		files := formatEmbedImages([]string{}, "./championImages/", fileName)
		embed = formatEmbedAuthor(embed, accInfo)
		embed = formatLastMatchEmbedFields(embed, matchresults, accInfo.Puuid)
		send = createMessageSend(embed, files)
		return send, nil
	}
	return send, errors.New("Sorry something went wrong getting lastmatch info for " + playerName)

}

//!lookup player
func LookupPlayer(playerName string) (send *discordgo.MessageSend, err error) {
	accInfo := getAccountInfo(playerName)
	send = &discordgo.MessageSend{}
	if accInfo != nil {
		rankedInfo := getRankedInfo(accInfo.Id)
		if rankedInfo == nil {
			return send, errors.New("Error getting match results for " + playerName)
		}
		fileName := getRankedAsset(rankedInfo)
		matchIDs, err := getMatchID(accInfo.Puuid, MATCH_LIMIT) // Request MATCH_LIMIT amount of match ID's to be later filtered out for ranked ones
		if err != nil {
			return send, errors.New("Error getting match results for " + playerName)
		}
		if len(matchIDs) < 1 {
			return send, errors.New("No match history found for " + playerName)
		}

		var matchStatsSlice []MatchResults
		for n, k := 0, 0; n < len(matchIDs) && k < NUM_OF_RANK_GAMES; n++ { // Get 10 games
			newMatch := getMatch(matchIDs[n])
			if newMatch == nil {
				return send, errors.New("Error getting match results for " + playerName)
			}
			if newMatch.Info.QueueId == RANKED_SOLO || newMatch.Info.QueueId == RANKED_FLEX { // But only if they are ranked_solo or ranked_flex games
				matchStatsSlice = append(matchStatsSlice, newMatch)
				k++
			}
		}

		description := formatPlayerRankedStats(rankedInfo)
		embed := formatRankedEmbed(playerName, fileName, description, 000127255, time.Now())
		embed = formatEmbedAuthor(embed, accInfo)

		if matchStatsSlice == nil {
			//Player has a rank, but no ranked matches within the last 30 games
			files := formatEmbedImages([]string{}, "./assets/", fileName)
			send = createMessageSend(embed, files)
			return send, nil
		}

		playermatchstats := formatMatchStats(matchStatsSlice, accInfo.Puuid)
		top3ChampStats := getTop3Champions(playermatchstats)
		var top3ChampNames []string
		for k := 0; k < len(top3ChampStats); k++ {
			top3ChampNames = append(top3ChampNames, top3ChampStats[k].Name)
		}

		embed = formatPlayerLookupEmbedFields(embed, playermatchstats, top3ChampNames)
		files := formatEmbedImages(top3ChampNames, "", fileName)
		send = createMessageSend(embed, files)
		return send, nil

	}
	return send, errors.New("Unable to get accInfo for: " + playerName)
}

///
///
///
func MasteryPlayer(playerName string) (send *discordgo.MessageSend, err error) {
	accInfo := getAccountInfo(playerName)
	send = &discordgo.MessageSend{}
	if accInfo != nil {
		rankedInfo := getRankedInfo(accInfo.Id)
		if rankedInfo == nil {
			return send, errors.New("Error getting match results for " + playerName)
		}
		fileName := getRankedAsset(rankedInfo)
		masteryStats := getMasteryData(accInfo.Id)
		if masteryStats == nil {
			return send, errors.New("Error getting masteries for " + playerName)
		}
		embed := formatRankedEmbed("Champion Masteries", fileName, "", 16747032, time.Now())
		embed = formatEmbedAuthor(embed, accInfo)
		files := formatEmbedImages([]string{}, "./assets/", fileName)
		embed = formatMasteriesEmbedFields(embed, masteryStats)
		send = createMessageSend(embed, files)
		return send, nil
	}
	return send, errors.New("Unable to get accInfo for: " + playerName)
}

func getEmbedColour(hasWon bool) int {
	if hasWon {
		return 28672 //Green
	}
	return 10747904 // Red
}

///
///
///
func createMessageSend(embed *discordgo.MessageEmbed, files []*discordgo.File) *discordgo.MessageSend {
	send := &discordgo.MessageSend{
		Embed: embed,
		Files: files,
	}
	return send
}

func InitEmojis(emoji [][]*discordgo.Emoji) {
	emojis = emoji
}

// A champion object is created on program start up containing all the names of all the champions and their ID's. Use this method to retrieve a name by ID
func GetChampion(champID string) string {
	for k, v := range champ3 {
		if champID == v.Key {
			return k // K is the champion name
		}
	}
	return champID
}

// An emoji object is created on program start up containing all the names and ID's of the emojis Leagly has access to.
func GetEmoji(emojiName string) string {
	for i := range emojis {
		for x := range emojis[i] {
			if emojis[i][x].Name == emojiName {
				return emojis[i][x].ID
			}
		}
	}
	return ""
}

///
///
///
func determineRoles(liveGameParticipants []LiveGameParticipants) {
	champPlayRates := ChampionPositions()
	var liveGameParticipantsBlue []LiveGameParticipants
	var liveGameParticipantsRed []LiveGameParticipants
	for i := 0; i < len(liveGameParticipants); i++ {
		if liveGameParticipants[i].TeamId == 100 {
			liveGameParticipantsBlue = append(liveGameParticipantsBlue, liveGameParticipants[i])
		} else {
			liveGameParticipantsRed = append(liveGameParticipantsRed, liveGameParticipants[i])
		}
	}
	// 1. Determine each role by champion ID for each team
	// 2. Must not have duplicate roles in each team, so we need to have a back up role(s)
	// 3. return an object with unique roles and the participant champs

	// filter out champs that are not part of the match
	// Ex: var playRatesRed []ChampionRole, playRatesRed = append(playRatesRed, getPlayRateChamps(liveGameParticipantsRed)
	// Ex: var playRatesBlue []ChampionRole, playRatesBlue = append(playRatesBlue, getPlayRateChamps(liveGameParticipantsBlue)
	champPlayRatesBlueTeam := getCurrentRoles(champPlayRates, liveGameParticipantsBlue)
	champPlayRatesRedTeam := getCurrentRoles(champPlayRates, liveGameParticipantsRed)

	// Now we restructure the slice
	// Ex: playRatesRed = reorderRoles(playRatesRed)
	// Ex: playRatesBlue = reorderRoles(playRatesBlue)
	determineRoleByPlayRate(champPlayRatesBlueTeam)
	determineRoleByPlayRate(champPlayRatesRedTeam)
	//then return the original reordered slice
	//var reorderedRoles []LiveGameParticipants
	//reorderedRoles = append(reorderedRoles, playRatesRed[0], playRatesRed[1]... playRatesBlue[3], playRatesBlue[4])
	//return reorderedRoles

}

func determineRoleByPlayRate(champPlayRates []ChampionRole) {
	roles := []string{"TOP", "JUNGLE", "MIDDLE", "BOTTOM", "UTILITY"}
	var ph int
	var bpr float32
	var prHolder []float32
	for k := 0; k < len(champPlayRates); k++ {
		prHolder = append(prHolder, champPlayRates[k].Top.PlayRate, champPlayRates[k].Jungle.PlayRate, champPlayRates[k].Middle.PlayRate, champPlayRates[k].Bottom.PlayRate, champPlayRates[k].Utility.PlayRate)
		for n := 0; n < len(prHolder); n++ {
			if bpr < prHolder[n] {
				bpr = prHolder[n]
				champPlayRates[k].Pos = roles[n]
				champPlayRates[k].PH = bpr
			}
		}
		prHolder = prHolder[:0] //still could have duplicates
		bpr = 0
		fmt.Println(fmt.Sprintf("Role: %s Rate: %f Champ: %s", champPlayRates[k].Pos, champPlayRates[k].PH, GetChampion(strconv.Itoa(champPlayRates[k].ID))))
	}
	fmt.Println(GetChampion(strconv.Itoa(champPlayRates[ph].ID)))
}

func getCurrentRoles(champPlayRates *map[string]ChampionRole, liveGameParticipants []LiveGameParticipants) []ChampionRole {
	var roles []ChampionRole
	for k, v := range *champPlayRates {
		if strconv.Itoa(liveGameParticipants[0].ChampionId) == k {
			roles = append(roles, v)
			roles[len(roles)-1].ID = liveGameParticipants[0].ChampionId
		} else if strconv.Itoa(liveGameParticipants[1].ChampionId) == k {
			roles = append(roles, v)
			roles[len(roles)-1].ID = liveGameParticipants[1].ChampionId
		} else if strconv.Itoa(liveGameParticipants[2].ChampionId) == k {
			roles = append(roles, v)
			roles[len(roles)-1].ID = liveGameParticipants[2].ChampionId
		} else if strconv.Itoa(liveGameParticipants[3].ChampionId) == k {
			roles = append(roles, v)
			roles[len(roles)-1].ID = liveGameParticipants[3].ChampionId
		} else if strconv.Itoa(liveGameParticipants[4].ChampionId) == k {
			roles = append(roles, v)
			roles[len(roles)-1].ID = liveGameParticipants[4].ChampionId
		}

	}
	return roles
}

///
///
///
func formatHelpEmbed(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   config.BotPrefix + "help",
			Value:  "Shows all available commands",
			Inline: false,
		},
		{
			Name:   config.BotPrefix + "live <playername>",
			Value:  "Checks to see if the player is in a game",
			Inline: false,
		},
		{
			Name:   config.BotPrefix + "lastmatch <playername>",
			Value:  "Shows the players last match stats",
			Inline: false,
		},
		{
			Name:   config.BotPrefix + "lookup <playername>",
			Value:  "Shows ranked history of player",
			Inline: false,
		},
		{
			Name:   config.BotPrefix + "mastery <playername>",
			Value:  "Shows mastery stats of player",
			Inline: false,
		},
	}
	return embed
}

///
///
///
func formatMasteriesEmbedFields(embed *discordgo.MessageEmbed, mastery Mastery) *discordgo.MessageEmbed {
	embed2 := &discordgo.MessageEmbed{}
	for n := 0; n+1 < len(mastery) && n < 10; n++ {
		embed2.Fields = []*discordgo.MessageEmbedField{
			{
				Name: "> <:" + strconv.Itoa(mastery[n].ChampionID) + ":" + GetEmoji(GetChampion(strconv.Itoa(mastery[n].ChampionID))) + ">" +
					"<:" + "Mastery" + strconv.Itoa(mastery[n].ChampionLevel) + ":" + GetEmoji(GetChampion("Mastery"+strconv.Itoa(mastery[n].ChampionLevel))) + ">",
				Value: "> <:" + strconv.Itoa(mastery[n+1].ChampionID) + ":" + GetEmoji(GetChampion(strconv.Itoa(mastery[n+1].ChampionID))) + ">" +
					"<:" + "Mastery" + strconv.Itoa(mastery[n+1].ChampionLevel) + ":" + GetEmoji(GetChampion("Mastery"+strconv.Itoa(mastery[n+1].ChampionLevel))) + ">",
				Inline: true,
			},
			{
				Name:   GetChampion(strconv.Itoa(mastery[n].ChampionID)),
				Value:  "**" + GetChampion(strconv.Itoa(mastery[n+1].ChampionID)) + "**",
				Inline: true,
			},
			{
				Name:   "__" + strconv.Itoa(mastery[n].ChampionPoints) + "__",
				Value:  "__" + "**" + strconv.Itoa(mastery[n+1].ChampionPoints) + "**__",
				Inline: true,
			},
		}
		embed.Fields = append(embed.Fields, embed2.Fields[0], embed2.Fields[1], embed2.Fields[2])
		n++
	}
	return embed
}

///
///
///
func formatLiveMatchEmbedFields(embed *discordgo.MessageEmbed, rankedPlayers []*RankedInfo, liveGameInfo LiveGameInfo, participant LiveGameParticipants) *discordgo.MessageEmbed {
	embed.Fields = []*discordgo.MessageEmbedField{
		// {
		// 	Name:   "CS",
		// 	Value:  fmt.Sprintf("```%d```", participant.TotalMinionsKilled+participant.NeutralMinionsKilled),
		// 	Inline: true,
		// },
		// {
		// 	Name:   "DMG Dealt",
		// 	Value:  fmt.Sprintf("```%d```", participant.TotalDamageDealtToChampions),
		// 	Inline: true,
		// },
		// {
		// 	Name:   "DMG Taken",
		// 	Value:  fmt.Sprintf("```%d```", participant.TotalDamageTaken),
		// 	Inline: true,
		// },
		{
			Name:   "\u200b",
			Value:  "\u200b",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "Blue Team",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "Red Team",
			Inline: true,
		},
	}
	embed2 := &discordgo.MessageEmbed{}
	roles := []string{"<:PositionTop:" + GetEmoji("PositionTop") + ">", "<:PositionJungle:" + GetEmoji("PositionJungle") + ">", "<:PositionMid:" + GetEmoji("PositionMid") + ">", "<:PositionBot:" + GetEmoji("PositionBot") + ">", "<:PositionSupport:" + GetEmoji("PositionSupport") + ">"}
	for k := 0; k < len(rankedPlayers)-5; k++ {
		embed2.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   roles[k],
				Value:  "\u200b",
				Inline: true,
			},
			{
				Name:   "> __<:" + GetChampion(strconv.Itoa(liveGameInfo.Participants[k].ChampionId)) + ":" + GetEmoji(GetChampion(strconv.Itoa(liveGameInfo.Participants[k].ChampionId))) + ">" + rankedPlayers[k].SummonerName + "__",
				Value:  fmt.Sprintf(">    WR: %d  ", rankedPlayers[k].Wins/(rankedPlayers[k].Wins+rankedPlayers[k].Losses)),
				Inline: true,
			},
			{
				Name:   "> __**<:" + GetChampion(strconv.Itoa(liveGameInfo.Participants[k+5].ChampionId)) + ":" + GetEmoji(GetChampion(strconv.Itoa(liveGameInfo.Participants[k+5].ChampionId))) + ">" + rankedPlayers[k+5].SummonerName + "**__",
				Value:  fmt.Sprintf(">    WR: %d  ", rankedPlayers[k].Wins/(rankedPlayers[k+5].Wins+rankedPlayers[k+5].Losses)),
				Inline: true,
			},
		}
		embed.Fields = append(embed.Fields, embed2.Fields[0], embed2.Fields[1], embed2.Fields[2])
	}
	return embed
}

///
///
///
func formatLastMatchEmbedFields(embed *discordgo.MessageEmbed, matchResults MatchResults, puuid string) *discordgo.MessageEmbed {
	participant := parseParticipant(puuid, matchResults)
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "CS",
			Value:  fmt.Sprintf("```%d```", participant.TotalMinionsKilled+participant.NeutralMinionsKilled),
			Inline: true,
		},
		{
			Name:   "DMG Dealt",
			Value:  fmt.Sprintf("```%d```", participant.TotalDamageDealtToChampions),
			Inline: true,
		},
		{
			Name:   "DMG Taken",
			Value:  fmt.Sprintf("```%d```", participant.TotalDamageTaken),
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "\u200b",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "Blue Team",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "Red Team",
			Inline: true,
		},
	}
	embed2 := &discordgo.MessageEmbed{}
	roles := []string{"<:PositionTop:" + GetEmoji("PositionTop") + ">", "<:PositionJungle:" + GetEmoji("PositionJungle") + ">", "<:PositionMid:" + GetEmoji("PositionMid") + ">", "<:PositionBot:" + GetEmoji("PositionBot") + ">", "<:PositionSupport:" + GetEmoji("PositionSupport") + ">"}
	for k := 0; k < len(matchResults.Info.Participants)-5; k++ {
		embed2.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   roles[k],
				Value:  "\u200b",
				Inline: true,
			},
			{
				Name:   "> __<:" + matchResults.Info.Participants[k].ChampionName + ":" + GetEmoji(matchResults.Info.Participants[k].ChampionName) + ">" + matchResults.Info.Participants[k].SummonerName + "__",
				Value:  fmt.Sprintf(">    %d / %d / %d ", matchResults.Info.Participants[k].Kills, matchResults.Info.Participants[k].Deaths, matchResults.Info.Participants[k].Assists),
				Inline: true,
			},
			{
				Name:   "> __**<:" + matchResults.Info.Participants[k+5].ChampionName + ":" + GetEmoji(matchResults.Info.Participants[k+5].ChampionName) + ">" + matchResults.Info.Participants[k+5].SummonerName + "**__",
				Value:  fmt.Sprintf(">    %d / %d / %d ", matchResults.Info.Participants[k+5].Kills, matchResults.Info.Participants[k+5].Deaths, matchResults.Info.Participants[k+5].Assists),
				Inline: true,
			},
		}
		embed.Fields = append(embed.Fields, embed2.Fields[0], embed2.Fields[1], embed2.Fields[2])
	}
	return embed
}

///
///
///
func formatPlayerLookupEmbedFields(embed *discordgo.MessageEmbed, playerMatchStats PlayerMatchStats, top3Champs []string) *discordgo.MessageEmbed {

	if len(top3Champs) < 1 {
		return embed
	}
	champData := []string{"\u200b", "\u200b", "\u200b", "\u200b", "\u200b", "\u200b"}

	var totalkills int
	var totalDeaths int
	var totalWins int
	var totalLoss int
	var totalAssist int
	for k := 0; k < len(playerMatchStats.PlayerChampions); k++ {
		totalkills += playerMatchStats.PlayerChampions[k].Kills
		totalDeaths += playerMatchStats.PlayerChampions[k].Deaths
		totalWins += playerMatchStats.PlayerChampions[k].Wins
		totalLoss += playerMatchStats.PlayerChampions[k].Loss
		totalAssist += playerMatchStats.PlayerChampions[k].Assists
	}

	// colourCypher := "css"
	winRate := (totalWins * 100) / (totalWins + totalLoss)
	// if winRate < 50 {
	// 	colourCypher = "diff "
	// }
	pickRate1, pickRate2 := getRole(playerMatchStats)
	KDA := fmt.Sprintf("```%dG %dW %dL (%d%%)\t \t %.1f / %.1f / %.1f KDA```", totalWins+totalLoss, totalWins, totalLoss, winRate,
		float64(totalkills/(totalWins+totalLoss)), float64(totalDeaths/(totalWins+totalLoss)), float64(totalAssist/(totalWins+totalLoss)))
	KDA1 := fmt.Sprintf("  Pick Rate %d%%", (pickRate1*100)/(totalWins+totalLoss))
	KDA2 := fmt.Sprintf("  Pick Rate %d%%", (pickRate2*100)/(totalWins+totalLoss))

	for j := 0; j < len(top3Champs); j++ { // This loop will iterate over the match history object that contains combined duplicate champion data. Creates unique data such as KDA per champion and win rates
		for k := 0; k < len(playerMatchStats.PlayerChampions); k++ {
			if top3Champs[j] == playerMatchStats.PlayerChampions[k].Name {
				if playerMatchStats.PlayerChampions[k].Deaths == 0 {
					playerMatchStats.PlayerChampions[k].Deaths = 1
				}
				champData[j] = fmt.Sprintf("WR:%d%% %dW|%dL\nKDA: %.1f", ((playerMatchStats.PlayerChampions[k].Wins * 100) / playerMatchStats.PlayerChampions[k].GamesPlayed),
					playerMatchStats.PlayerChampions[k].Wins, playerMatchStats.PlayerChampions[k].Loss, (float64(playerMatchStats.PlayerChampions[k].Kills+playerMatchStats.PlayerChampions[k].Assists))/(float64(playerMatchStats.PlayerChampions[k].Deaths)))
				champData[j+3] = playerMatchStats.PlayerChampions[k].Name
			}
		}
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "\u200b",
			Value:  "\u200b",
			Inline: false,
		},
		{
			Name:  "Past 10 games stats: \t\t",
			Value: KDA,
		},
		{
			Name:   "Primary Role:",
			Value:  "```" + playerMatchStats.Role.PreferredRole1 + KDA1 + "```",
			Inline: true,
		},
		{
			Name:   "Secondary Role:",
			Value:  "```" + playerMatchStats.Role.PreferredRole2 + KDA2 + "```",
			Inline: true,
		},
		{
			Name:  "\u200b",
			Value: "Top 3 Recent Champions",
		},
		{
			Name:   champData[3],
			Value:  "```" + champData[0] + "```",
			Inline: true,
		},
		{
			Name:   champData[4],
			Value:  "```" + champData[1] + "```",
			Inline: true,
		},
		{
			Name:   champData[5],
			Value:  "```" + champData[2] + "```",
			Inline: true,
		},
	}
	embed.Image = &discordgo.MessageEmbedImage{
		URL: "attachment://output.png",
	}
	return embed
}

///
///
///
func formatRankedEmbed(playerName string, fileName string, description string, colour int, times time.Time) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:       colour,
		Title:       playerName,
		Description: description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "attachment://" + fileName,
		},
		Timestamp: times.Format(time.RFC3339),
	}
	return embed
}

///
///
///
func formatEmbedAuthor(embed *discordgo.MessageEmbed, playerInfo Summoner) *discordgo.MessageEmbed {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    playerInfo.Name,
		IconURL: "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/profileicon/" + strconv.Itoa(playerInfo.ProfileIconId) + ".png",
		URL:     "https://na.op.gg/summoner/userName=" + strings.ReplaceAll(playerInfo.Name, " ", "%20"),
	}
	return embed
}

///
///
///
func getChampionFile(filename string) (err error) {
	URL := "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/"
	if _, err := os.Stat("./championImages/" + filename); errors.Is(err, os.ErrNotExist) {
		errs := downloadFile(URL+filename, filename) //champion icons are only downloaded if they don't exist in the "championImages" directory
		return errs
	}
	return err
}

///
///
///
func formatItems(participant Participants) string {

	res := fmt.Sprintf("Items: <:%d:%s> <:%d:%s> <:%d:%s> <:%d:%s> <:%d:%s> <:%d:%s>",
		participant.Item0, GetEmoji(strconv.Itoa(participant.Item0)),
		participant.Item1, GetEmoji(strconv.Itoa(participant.Item1)),
		participant.Item2, GetEmoji(strconv.Itoa(participant.Item2)),
		participant.Item3, GetEmoji(strconv.Itoa(participant.Item3)),
		participant.Item4, GetEmoji(strconv.Itoa(participant.Item4)),
		participant.Item5, GetEmoji(strconv.Itoa(participant.Item5)))

	space := regexp.MustCompile(`<:0:>`)
	return space.ReplaceAllString(res, " ")
}

///
///
///
func formatEmbedImages(imageNames []string, relativePath string, rankFileName string) []*discordgo.File {
	var files []*discordgo.File

	for n := 0; n < len(imageNames); n++ {
		imageNames[n] += ".png"
		err := getChampionFile(imageNames[n])
		if err != nil {
			return files
		}
		imageNames[n] = "./championImages/" + imageNames[n]
	}

	if len(imageNames) > 0 {
		fileImageName := mergeImages(imageNames)
		file, _ := os.Open(fileImageName)
		files = append(files, &discordgo.File{
			Name:        fileImageName,
			ContentType: "image/png",
			Reader:      file,
		})
		file2, _ := os.Open("./assets/" + rankFileName)
		files = append(files, &discordgo.File{
			Name:        rankFileName,
			ContentType: "image/png",
			Reader:      file2,
		})
	} else {
		file, _ := os.Open(relativePath + rankFileName) // actually champion image for lastmatch & live
		files = append(files, &discordgo.File{
			Name:        rankFileName,
			ContentType: "image/png",
			Reader:      file,
		})
	}

	return files
}

//This method iterates through the bulk matchresult struct and combine the select players game data by champion.
//A new struct will be returned that contains match results by unique champion
func formatMatchStats(matchedStats []MatchResults, puuid string) PlayerMatchStats {
	var win int
	var loss int
	var playermatchstats PlayerMatchStats
	set := make(map[string]struct{})

	for n := 0; n < len(matchedStats); n++ {
		participant := parseParticipant(puuid, matchedStats[n])
		playermatchstats.ProfileIcon = participant.ProfileIcon
		playermatchstats.SummonerName = participant.SummonerName

		if participant.Win {
			win++
		} else {
			loss++
		}

		if len(playermatchstats.PlayerChampions) == 0 {
			playermatchstats.PlayerChampions = append(playermatchstats.PlayerChampions, &PlayerChampions{
				Name:        participant.ChampionName,
				Kills:       participant.Kills,
				Deaths:      participant.Deaths,
				Assists:     participant.Assists,
				Wins:        win,
				Loss:        loss,
				GamesPlayed: 1,
			})
		} else {
			counter := len(playermatchstats.PlayerChampions)
			for k := 0; k < counter; k++ {
				set[playermatchstats.PlayerChampions[k].Name] = struct{}{} //created a map to keep track of the champions that have been so far looked through
				if playermatchstats.PlayerChampions[k].Name == participant.ChampionName {
					playermatchstats.PlayerChampions[k].Kills += participant.Kills
					playermatchstats.PlayerChampions[k].Deaths += participant.Deaths
					playermatchstats.PlayerChampions[k].Assists += participant.Assists
					playermatchstats.PlayerChampions[k].Wins += win
					playermatchstats.PlayerChampions[k].Loss += loss
					playermatchstats.PlayerChampions[k].GamesPlayed++
					break
				}
				if _, ok := set[participant.ChampionName]; ok { //if our map doesn't contain the champion the loop is iterating over, we append it to the object
				} else {
					playermatchstats.PlayerChampions = append(playermatchstats.PlayerChampions, &PlayerChampions{
						Name:        participant.ChampionName,
						Kills:       participant.Kills,
						Deaths:      participant.Deaths,
						Assists:     participant.Assists,
						Wins:        win,
						Loss:        loss,
						GamesPlayed: 1,
					})
					set[playermatchstats.PlayerChampions[len(playermatchstats.PlayerChampions)-1].Name] = struct{}{}
				}
			}
		}

		win = 0
		loss = 0
		if participant.IndividualPosition == "TOP" {
			playermatchstats.Role.RoleCount[0]++
		} else if participant.IndividualPosition == "JUNGLE" {
			playermatchstats.Role.RoleCount[1]++
		} else if participant.IndividualPosition == "MIDDLE" {
			playermatchstats.Role.RoleCount[2]++
		} else if participant.IndividualPosition == "BOTTOM" {
			playermatchstats.Role.RoleCount[3]++
		} else {
			playermatchstats.Role.RoleCount[4]++
		}
	}

	pHolder := getFavouriteRole(playermatchstats.Role, -1)
	playermatchstats.Role.PreferredRole1 = getFavouriteRoleName(pHolder)
	playermatchstats.Role.PreferredRole2 = getFavouriteRoleName(getFavouriteRole(playermatchstats.Role, pHolder))

	return playermatchstats
}

// the following 3 methods are my lazy way of determining role and favourite role
func getFavouriteRole(playerRoles Role, ignore int) int {
	largest := 0
	pHolder := -1 // set to -1 to ensure secondary role isn't duplicated. (will be skipped in the method below)
	for j := 0; j < len(playerRoles.RoleCount); j++ {
		if j == ignore {
			continue
		}
		if largest < playerRoles.RoleCount[j] {
			pHolder = j
			largest = playerRoles.RoleCount[j]
		}
	}
	return pHolder
}

///
///
///
func getFavouriteRoleName(pHolder int) string {

	switch pHolder {
	case 0:
		{
			return "TOP"
		}
	case 1:
		{
			return "JUNGLE"
		}
	case 2:
		{
			return "MIDDLE"
		}
	case 3:
		{
			return "BOTTOM"
		}
	case 4:
		{
			return "UTILITY"
		}
	}
	return "UNKOWN"
}

///
///
///
func getRole(role PlayerMatchStats) (int, int) {
	roles := []string{"TOP", "JUNGLE", "MIDDLE", "BOTTOM", "UTILITY"}
	var pr int
	var pr2 int
	for k := 0; k < len(roles); k++ {
		if role.Role.PreferredRole1 == roles[k] {
			pr = role.Role.RoleCount[k]
		}
		if role.Role.PreferredRole2 == roles[k] {
			pr2 = role.Role.RoleCount[k]
		}
	}
	return pr, pr2
}

///
///
///
func getTop3Champions(playerMatchStats PlayerMatchStats) []*PlayerChampions {
	var playerChampions []*PlayerChampions
	if len(playerMatchStats.PlayerChampions) < 1 { // We return here if no champion data exists, but we still want to create an embed object later on
		return playerChampions
	}
	var playerChampion []*PlayerChampions
	for k := 0; k < 3; k++ {
		playerChampion = append(playerChampion, &PlayerChampions{
			GamesPlayed: 0,
		})
	}
	// At this point duplicate champions in match history has been combined. Now the program will choose the top 3 champions determined by games played

	for i := 0; i < len(playerMatchStats.PlayerChampions); i++ {
		if playerChampion[0].GamesPlayed < playerMatchStats.PlayerChampions[i].GamesPlayed {
			playerChampion[0] = playerMatchStats.PlayerChampions[i]
			i = 0
		}
		if playerMatchStats.PlayerChampions[i].Name != playerChampion[0].Name {
			if playerChampion[1].GamesPlayed < playerMatchStats.PlayerChampions[i].GamesPlayed {
				playerChampion[1] = playerMatchStats.PlayerChampions[i]
				i = 0
			}
		}
		if playerMatchStats.PlayerChampions[i].Name != playerChampion[0].Name && playerMatchStats.PlayerChampions[i].Name != playerChampion[1].Name {
			if playerChampion[2].GamesPlayed < playerMatchStats.PlayerChampions[i].GamesPlayed {
				playerChampion[2] = playerMatchStats.PlayerChampions[i]
				i = 0
			}
		}
	}

	for k := 0; k < 3; k++ {
		if playerChampion[k].GamesPlayed > 0 {
			playerChampions = append(playerChampions, playerChampion[k])
		}
	}
	return playerChampions
}

// Ranked icon images are locally stored. This method is used to determine which ranked icon image we need.
func getRankedAsset(rankedStats []*RankedInfo) string {
	for n := 0; n < len(rankedStats); n++ {
		if rankedStats[n].QueueType == "RANKED_SOLO_5x5" { //Player can have 2 different ranks in random order. We want to prioritize the solo rank
			rank := getRankFile(rankedStats[n].Tier)
			if rank == "" {
				continue
			}
			return rank
		}
	}
	for n := 0; n < len(rankedStats); n++ {
		if rankedStats[n].QueueType == "RANKED_TEAM_5x5" || rankedStats[n].QueueType == "RANKED_FLEX_SR" {
			rank := getRankFile(rankedStats[n].Tier)
			if rank == "" {
				continue
			}
			return rank
		}
	}
	return "UNRANKED.png"
}

func getRankFile(rank string) string {
	switch {
	case rank == "IRON":
		return "Emblem_Iron.png"
	case rank == "BRONZE":
		return "Emblem_Bronze.png"
	case rank == "SILVER":
		return "Emblem_Silver.png"
	case rank == "GOLD":
		return "Emblem_Gold.png"
	case rank == "PLATINUM":
		return "Emblem_Platinum.png"
	case rank == "DIAMOND":
		return "Emblem_Diamond.png"
	case rank == "MASTER":
		return "Emblem_Master.png"
	case rank == "GRANDMASTER":
		return "Emblem_Grandmaster.png"
	case rank == "CHALLENGER":
		return "Emblem_Challenger.png"
	}
	return ""
}

///
///
///
func getMatchType(queueType int) string {
	if queueType == RANKED_SOLO || queueType == RANKED_FLEX {
		return "Summoners Rift Ranked"
	} else if queueType == NORMAL {
		return "Summoners Rift Normal"
	} else if queueType == URF {
		return "Summoners Rift URF"
	} else if queueType == ARAM {
		return "Howling Abyss ARAM"
	}
	return "Custom Game"
}

// When calling the Riot API for match data, we get a large json object with match data of all 10 players.
// This method is used to filter out each player and only returns an object of the one we're looking for
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

/// same method as above but could not reuse because livegame is pulled into a different data object
///
///
func parseLiveParticipant(sumID string, liveGameInfo LiveGameInfo) LiveGameParticipants {
	var i int
	for n := 0; n < len(liveGameInfo.Participants); n++ {
		if sumID == liveGameInfo.Participants[n].SummonerId {
			i = n
			break
		}
	}
	return liveGameInfo.Participants[i]
}

// Because discord embeds only support 2x2 images at a maximum, I decided to use a method
// that combines 3 images into one to be use in a 1x3 format. Unfortunately discord also
// has limitations on image size. Embed will be constrained if the image is greater than 300px
// As a result, i limit 3 images combined to be at a maximum of 299 px in length :D
func mergeImages(imageName []string) string {

	var imgFile []*os.File
	var img []image.Image
	for n := 0; n < len(imageName); n++ {
		imgFile1, err := os.Open(imageName[n])
		if err != nil {
			fmt.Println(err)
		}
		imgFile = append(imgFile, imgFile1)
		img1, _, err := image.Decode(imgFile1)
		if err != nil {
			fmt.Println(err)
		}
		img = append(img, img1)
	}

	var sp3 image.Point
	var sp image.Point
	if len(img) == 3 {
		sp = image.Point{(img[0].Bounds().Dx() - 20), 0}
		sp3 = image.Point{sp.X + sp.X, 0}
	} else if len(img) == 2 {
		sp = image.Point{(img[0].Bounds().Dx()), 0}
		sp3 = image.Point{sp.X, 0}
	} else {
		sp = image.Point{(img[0].Bounds().Dx()), 0}
		sp3 = image.Point{0, 0}
	}
	r2 := image.Rectangle{sp, sp.Add(img[0].Bounds().Size())}

	r3 := image.Rectangle{sp3, sp3.Add(img[0].Bounds().Size())} //all images are same size anyways
	if len(img) == 3 {
		r3.Max.X = 299 //Discord embeed width will be constrained if the image is 300px in width or greater
	} else if len(img) == 2 {
		r3.Max.X = sp.X + sp.X
	} else {
		r3.Max.X = sp.X
	}
	r := image.Rectangle{image.Point{0, 0}, r3.Max}

	rgba := image.NewRGBA(r)
	draw.Draw(rgba, img[0].Bounds(), img[0], image.Point{0, 0}, draw.Src)
	if len(img) == 2 {
		draw.Draw(rgba, r2, img[1], image.Point{0, 0}, draw.Src)
	} else if len(img) == 3 {
		draw.Draw(rgba, r2, img[1], image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r3, img[2], image.Point{0, 0}, draw.Src)
	}

	out, err := os.Create("./output.png")
	if err != nil {
		fmt.Println(err)
	}
	var opt jpeg.Options
	opt.Quality = 80

	jpeg.Encode(out, rgba, &opt)
	return "./output.png"
}

///
///
///
func formatRankedPlayers(liveGameInfo LiveGameInfo) []*RankedInfo {
	var rankedPlayers []*RankedInfo
	for i := 0; i < len(liveGameInfo.Participants); i++ {
		rankHistory := getRankedInfo(liveGameInfo.Participants[i].SummonerId)
		for n := 0; n < len(rankHistory); n++ {
			if rankHistory[n].QueueType == "RANKED_SOLO_5x5" {
				rankedPlayers = append(rankedPlayers, rankHistory[n])
				break
			}
		}
		if len(rankedPlayers) <= i {
			rankedPlayers = append(rankedPlayers, &RankedInfo{Tier: "NA", Rank: "NA", Losses: 1, SummonerName: liveGameInfo.Participants[i].SummonerName})
		}
	}
	return rankedPlayers
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

///
///
///
func formatPlayerRankedStats(rankedStats []*RankedInfo) string {
	for n := 0; n < len(rankedStats); n++ {
		if rankedStats[n].QueueType == "RANKED_TFT_PAIRS" {
			continue
		}
		if rankedStats[n].QueueType == "RANKED_SOLO_5x5" || rankedStats[n].QueueType == "RANKED_TEAM_5x5 " || rankedStats[n].QueueType == "RANKED_FLEX_SR" {
			for k := 0; k < len(rankedStats); k++ { //Look again because we want to prioritize solo duo rank over flex rank
				if rankedStats[k].QueueType == "RANKED_SOLO_5x5" {
					return rankedStats[k].Tier + " " + rankedStats[k].Rank +
						" with " + strconv.Itoa(rankedStats[k].LeaguePoints) + " LP. Season W/L: " + strconv.Itoa(rankedStats[k].Wins) + " wins and " + strconv.Itoa(rankedStats[k].Losses) + " losses. WR: " + strconv.Itoa((rankedStats[k].Wins*100)/(rankedStats[k].Wins+rankedStats[k].Losses)) + "%"
				}
			}
			return rankedStats[n].Tier + " " + rankedStats[n].Rank +
				" with " + strconv.Itoa(rankedStats[n].LeaguePoints) + " LP. Season W/L: " + strconv.Itoa(rankedStats[n].Wins) + " wins and " + strconv.Itoa(rankedStats[n].Losses) + " losses. WR: " + strconv.Itoa((rankedStats[n].Wins*100)/(rankedStats[n].Wins+rankedStats[n].Losses)) + "%"
		} else {
			return "```Currently unranked.```"
		}
	}
	return "```No ranked data found```"
}
