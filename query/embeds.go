package query

import (
	"Leagly/config"
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
	send := &discordgo.MessageSend{
		Embed: embed,
		Files: files,
	}
	return send
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
				Value:  fmt.Sprintf(">    WR: %d%%", (rankedPlayers[k].Wins*100)/(rankedPlayers[k].Wins+rankedPlayers[k].Losses)),
				Inline: true,
			},
			{
				Name:   "> __**<:" + GetChampion(strconv.Itoa(liveGameInfo.Participants[k+5].ChampionId)) + ":" + GetEmoji(GetChampion(strconv.Itoa(liveGameInfo.Participants[k+5].ChampionId))) + ">" + rankedPlayers[k+5].SummonerName + "**__",
				Value:  fmt.Sprintf(">    WR: %d%%", (rankedPlayers[k+5].Wins*100)/(rankedPlayers[k+5].Wins+rankedPlayers[k+5].Losses)),
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
