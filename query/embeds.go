package query

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

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
	embeds := []*discordgo.MessageEmbed{embed}
	send := &discordgo.MessageSend{
		Embeds: embeds,
		Files:  files,
	}
	return send
}

func formatePatchNotesEmbed(embed *discordgo.MessageEmbed, patchnotesURL string) *discordgo.MessageEmbed {
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Current Patch",
			Value:  "[This Patch](" + patchnotesURL + ")",
			Inline: true,
		},
		{
			Name:   "Past Updates",
			Value:  "[Patch History](" + PATCH_NOTES_HISTORY_URL + ")",
			Inline: true,
		},
	}
	return embed
}

func formatApiStatusEmbed(embed *discordgo.MessageEmbed, riotStatus RiotStatus, lang string) *discordgo.MessageEmbed {
	if len(riotStatus.Incidents) < 1 {
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "No incidents found.",
				Value:  "\u200b",
				Inline: false,
			},
		}
	}
	embed2 := &discordgo.MessageEmbed{}
	var statusTitle string
	var statusMsg string
	for n := 0; n < len(riotStatus.Incidents); n++ {
		for m := 0; m < len(riotStatus.Incidents[n].Titles); m++ {
			if lang == riotStatus.Incidents[n].Titles[m].Locale {
				statusTitle = riotStatus.Incidents[n].Titles[m].Content
				statusMsg = riotStatus.Incidents[n].Updates[0].Translations[m].Content
			}
		}
		embed2.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   statusTitle,
				Value:  statusMsg,
				Inline: false,
			},
		}
		embed.Fields = append(embed.Fields, embed2.Fields[0])
	}
	return embed
}

///
///
///
func formatHelpEmbed(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "/help",
			Value:  "Shows all available commands",
			Inline: false,
		},
		{
			Name:   "/live <playername>",
			Value:  "Checks to see if the player is in a game",
			Inline: false,
		},
		{
			Name:   "/lastmatch <playername>",
			Value:  "Shows the players last match stats",
			Inline: false,
		},
		{
			Name:   "/lookup <playername>",
			Value:  "Shows ranked history of player",
			Inline: false,
		},
		{
			Name:   "/mastery <playername>",
			Value:  "Shows mastery stats of player",
			Inline: false,
		},
		{
			Name:   "/region <region code>",
			Value:  "Updates the region for your discord server",
			Inline: false,
		},
		{
			Name:   "\u200b",
			Value:  "[Join Leagly Discord](https://discord.gg/bxQRKA8D9g)\n",
			Inline: false,
		},
	}
	return embed
}

///
///
///
func formatMasteriesEmbedFields(embed *discordgo.MessageEmbed, mastery Mastery) *discordgo.MessageEmbed {
	var championString string
	var masteryString string
	var chestString string

	var totalPoints int
	totalChampions := len(mastery)
	var masteryTokens int
	for n := 0; n < len(mastery); n++ {
		totalPoints += mastery[n].ChampionPoints
		masteryTokens += mastery[n].TokensEarned
	}

	for n := 0; n < len(mastery) && n < 10; n++ {
		championString += "> <:" + strconv.Itoa(mastery[n].ChampionID) + ":" + GetEmoji(GetChampion(strconv.Itoa(mastery[n].ChampionID))) + ">" +
			"<:" + "Mastery" + strconv.Itoa(mastery[n].ChampionLevel) + ":" + GetEmoji(GetChampion("Mastery"+strconv.Itoa(mastery[n].ChampionLevel))) + "> " + GetChampion(strconv.Itoa(mastery[n].ChampionID)) + "\n"
		masteryString += strconv.Itoa(mastery[n].ChampionPoints) + "\n"
		chestString += " <:Chest" + strconv.FormatBool(mastery[n].ChestGranted) + ":" + GetEmoji("Chest"+strconv.FormatBool(mastery[n].ChestGranted)) + ">\n"
	}
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Champion",
			Value:  championString,
			Inline: true,
		},
		{
			Name:   "Mastery",
			Value:  masteryString,
			Inline: true,
		},
		{
			Name:   "Chests",
			Value:  chestString,
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "Champions: \n __" + strconv.Itoa(totalChampions) + "__",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "Total Mastery Points:\n__" + strconv.Itoa(totalPoints) + "__",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "Tokens : \n__" + strconv.Itoa(masteryTokens) + "__",
			Inline: true,
		},
	}

	return embed
}

///
///
///
func formatLiveMatchEmbedFields(embed *discordgo.MessageEmbed, rankedPlayers []*RankedInfo, liveGameInfo LiveGameInfo, participant LiveGameParticipants, bannedChampions string) *discordgo.MessageEmbed {
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Bans",
			Value: bannedChampions,
		},
		{
			Name:   "\u200b",
			Value:  "\u200b",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  fmt.Sprintf("<:%s:%s>", "blue_team", GetEmoji("blue_team")) + "Blue Team",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  fmt.Sprintf("<:%s:%s>", "red_team", GetEmoji("red_team")) + "Red Team",
			Inline: true,
		},
	}
	embed2 := &discordgo.MessageEmbed{}
	roles := []string{"<:PositionTop:" + GetEmoji("PositionTop") + ">", "<:PositionJungle:" + GetEmoji("PositionJungle") + ">", "<:PositionMid:" + GetEmoji("PositionMid") + ">", "<:PositionBot:" + GetEmoji("PositionBot") + ">", "<:PositionSupport:" + GetEmoji("PositionSupport") + ">"}
	for k := 0; k < len(rankedPlayers)-5; k++ {
		if rankedPlayers[k+5].Rank == "NA" {
			rankedPlayers[k+5].Rank = " "
		}
		if rankedPlayers[k].Rank == "NA" {
			rankedPlayers[k].Rank = " "
		}
		embed2.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   roles[k],
				Value:  "\u200b",
				Inline: true,
			},
			{
				Name:   "> __<:" + GetChampion(strconv.Itoa(liveGameInfo.Participants[k].ChampionId)) + ":" + GetEmoji(GetChampion(strconv.Itoa(liveGameInfo.Participants[k].ChampionId))) + ">" + rankedPlayers[k].SummonerName + "__",
				Value:  fmt.Sprintf("> <:%s:%s>**%s**   WR: %d%% (%dG)", rankedPlayers[k].Tier, GetEmoji(rankedPlayers[k].Tier), rankedPlayers[k].Rank, (rankedPlayers[k].Wins*100)/(rankedPlayers[k].Wins+rankedPlayers[k].Losses), rankedPlayers[k].Wins+rankedPlayers[k].Losses),
				Inline: true,
			},
			{
				Name:   "> __**<:" + GetChampion(strconv.Itoa(liveGameInfo.Participants[k+5].ChampionId)) + ":" + GetEmoji(GetChampion(strconv.Itoa(liveGameInfo.Participants[k+5].ChampionId))) + ">" + rankedPlayers[k+5].SummonerName + "**__",
				Value:  fmt.Sprintf("> <:%s:%s>**%s**   WR: %d%% (%dG)", rankedPlayers[k+5].Tier, GetEmoji(rankedPlayers[k+5].Tier), rankedPlayers[k+5].Rank, (rankedPlayers[k+5].Wins*100)/(rankedPlayers[k+5].Wins+rankedPlayers[k+5].Losses), rankedPlayers[k+5].Wins+rankedPlayers[k+5].Losses),
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
			Value:  fmt.Sprintf("<:%s:%s>", "blue_team", GetEmoji("blue_team")) + "Blue Team",
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  fmt.Sprintf("<:%s:%s>", "red_team", GetEmoji("red_team")) + "Red Team",
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

	winRate := (totalWins * 100) / (totalWins + totalLoss)

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
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Please share Leagly if you enjoy the app!",
		},
		Timestamp: times.Format(time.RFC3339),
	}
	return embed
}

///
///
///
func formatEmbedAuthor(embed *discordgo.MessageEmbed, playerInfo Summoner, region string) *discordgo.MessageEmbed {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    playerInfo.Name + " [" + region + "]",
		IconURL: "http://ddragon.leagueoflegends.com/cdn/" + Version + "/img/profileicon/" + strconv.Itoa(playerInfo.ProfileIconId) + ".png",
		URL:     "https://na.op.gg/summoner/userName=" + strings.ReplaceAll(playerInfo.Name, " ", "%20"),
	}
	return embed
}

///
///
///
func formatEmbedAuthorLeagly(embed *discordgo.MessageEmbed, name string, iconUrl string) *discordgo.MessageEmbed {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		IconURL: iconUrl,
		URL:     "https://discord.com/api/oauth2/authorize?client_id=930924283599925260&permissions=2147798016&scope=bot%20applications.commands",
	}
	return embed
}
