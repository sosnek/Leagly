package query

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var champ3 map[string]Champion

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
func LookupPlayer(playerName string) (send *discordgo.MessageSend, err error) {

	accInfo, exists := getAccountInfo(playerName)
	send = &discordgo.MessageSend{}
	if exists {
		rankedInfo := getRankedInfo(accInfo.Id)
		fileName, rankedType := getRankedAsset(RankedInfo(rankedInfo))
		matchIDs, exist := getMatchID(accInfo.Puuid, MATCH_LIMIT)
		if len(matchIDs) < 1 {
			return send, errors.New("No match history found for " + playerName)
		}

		if exist {
			var matchStatsSlice []MatchResults
			for n, k := 0, 0; n < len(matchIDs) && k < 10; n++ {
				newMatch := getMatch(matchIDs[n])
				if newMatch.Info.QueueId == RANKED_SOLO || newMatch.Info.QueueId == RANKED_FLEX {
					matchStatsSlice = append(matchStatsSlice, newMatch)
					k++
				}
			}

			description := formatPlayerRankedStats(rankedInfo)
			embed := formatRankedEmbed(playerName, rankedType, fileName, description)
			embed = formatEmbedAuthor(embed, accInfo)

			if matchStatsSlice == nil {
				return &discordgo.MessageSend{Embed: embed}, errors.New("No rank history found for " + playerName)
			}

			playermatchstats := formatMatchStats(matchStatsSlice, accInfo.Puuid)
			top3ChampStats := getTop3Champions(playermatchstats)
			var top3ChampNames []string
			for k := 0; k < len(top3ChampStats); k++ {
				top3ChampNames = append(top3ChampNames, top3ChampStats[k].Name)
			}

			//fmt.Println(matchStatsFormatted)

			embed = formatPlayerLookupEmbedFields(embed, playermatchstats, top3ChampNames)
			files := formatEmbedImages(embed, top3ChampNames, fileName)
			send = createMessageSend(embed, files)
			return send, err
		}
	}
	log.Println("Unable to get accInfo for: " + playerName)
	return send, errors.New("Unable to get accInfo for: " + playerName)
}

func MasteryPlayer(playerName string) (send *discordgo.MessageSend, err error) {
	accInfo, exists := getAccountInfo(playerName)
	send = &discordgo.MessageSend{}
	if exists {
		masteryStats := getMasteryData(accInfo.Id)
		//masteryStatsFormatted := formatMasteries(masteryStats)
		fmt.Println(masteryStats)
	}
	log.Println("Unable to get accInfo for: " + playerName)
	return send, errors.New("Unable to get accInfo for: " + playerName)
}

func DeleteImages(fileNames []string) {
	for k := 0; k < len(fileNames); k++ {
		e := os.Remove(fileNames[k])
		if e != nil {
			log.Fatal(e)
		}
	}
}

func createMessageSend(embed *discordgo.MessageEmbed, files []*discordgo.File) *discordgo.MessageSend {
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

	for j := 0; j < len(top3Champs); j++ {
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

func formatRankedEmbed(playerName string, rankedType int, fileName string, description string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:       000127255,
		Title:       playerName,
		Description: description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "attachment://" + fileName,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return embed
}

func formatEmbedAuthor(embed *discordgo.MessageEmbed, playerInfo Summoner) *discordgo.MessageEmbed {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    playerInfo.Name,
		IconURL: "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/profileicon/" + strconv.Itoa(playerInfo.ProfileIconId) + ".png",
		URL:     "https://na.op.gg/summoner/userName=" + strings.ReplaceAll(playerInfo.Name, " ", "%20"),
	}
	return embed
}

func formatEmbedImages(embed *discordgo.MessageEmbed, imageNames []string, rankFileName string) []*discordgo.File {
	var files []*discordgo.File
	file2, _ := os.Open("./assets/" + rankFileName)
	files = append(files, &discordgo.File{
		Name:        rankFileName,
		ContentType: "image/png",
		Reader:      file2,
	})

	URL := "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/"
	for n := 0; n < len(imageNames); n++ {
		imageNames[n] += ".png"
		if _, err := os.Stat("./championImages/" + imageNames[n]); errors.Is(err, os.ErrNotExist) {
			err = downloadFile(URL+imageNames[n], imageNames[n])
			if err != nil {
				log.Fatal("Unable to download file")
				return files
			}
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
	}

	return files
}

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
				set[playermatchstats.PlayerChampions[k].Name] = struct{}{}
				if playermatchstats.PlayerChampions[k].Name == participant.ChampionName {
					playermatchstats.PlayerChampions[k].Kills += participant.Kills
					playermatchstats.PlayerChampions[k].Deaths += participant.Deaths
					playermatchstats.PlayerChampions[k].Assists += participant.Assists
					playermatchstats.PlayerChampions[k].Wins += win
					playermatchstats.PlayerChampions[k].Loss += loss
					playermatchstats.PlayerChampions[k].GamesPlayed++
					break
				}
				if _, ok := set[participant.ChampionName]; ok {
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

func getFavouriteRole(playerRoles Role, ignore int) int {
	largest := 0
	pHolder := -1
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

func getTop3Champions(playerMatchStats PlayerMatchStats) []*PlayerChampions {
	var playerChampions []*PlayerChampions
	if len(playerMatchStats.PlayerChampions) < 1 {
		return playerChampions
	}
	var playerChampion []*PlayerChampions
	for k := 0; k < 3; k++ {
		playerChampion = append(playerChampion, &PlayerChampions{
			GamesPlayed: 0,
		})
	}

	for k := 1; k < len(playerMatchStats.PlayerChampions); k++ {
		if playerChampion[0].GamesPlayed <= playerMatchStats.PlayerChampions[k].GamesPlayed {
			playerChampion[0] = playerMatchStats.PlayerChampions[k]
		}
	}
	for k := 0; k < len(playerMatchStats.PlayerChampions); k++ {
		if playerMatchStats.PlayerChampions[k].Name != playerChampion[0].Name {
			if playerChampion[1].GamesPlayed <= playerMatchStats.PlayerChampions[k].GamesPlayed {
				playerChampion[1] = playerMatchStats.PlayerChampions[k]
			}
		}
	}
	for k := 0; k < len(playerMatchStats.PlayerChampions); k++ {
		if playerMatchStats.PlayerChampions[k].Name != playerChampion[0].Name && playerMatchStats.PlayerChampions[k].Name != playerChampion[1].Name {
			if playerChampion[2].GamesPlayed <= playerMatchStats.PlayerChampions[k].GamesPlayed {
				playerChampion[2] = playerMatchStats.PlayerChampions[k]
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

func formatPlayerRankedStats(rankedStats RankedInfo) string {
	for n := 0; n < len(rankedStats); n++ {
		if rankedStats[n].QueueType == "RANKED_SOLO_5x5" || rankedStats[n].QueueType == "RANKED_TEAM_5x5 " {
			return rankedStats[n].Tier + " " + rankedStats[n].Rank +
				" with " + strconv.Itoa(rankedStats[n].LeaguePoints) + " LP. Season W/L: " + strconv.Itoa(rankedStats[n].Wins) + " wins and " + strconv.Itoa(rankedStats[n].Losses) + " losses. WR: " + strconv.Itoa((rankedStats[n].Wins*100)/(rankedStats[n].Wins+rankedStats[n].Losses)) + "%"
		} else {
			return "Currently unranked."
		}

	}
	return "```No ranked data found```"
}
