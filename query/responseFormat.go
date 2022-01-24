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
	Role        string
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
		//masteryStats := getMasteryData(accInfo.Id)
		matchStatsID, exist := getMatchID(accInfo.Puuid, 40)

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
			top3ChampStats := getTop3Champions(playermatchstats)
			if len(top3ChampStats) < 1 {
				return send, errors.New("No match history found for " + playerName)
			}
			top3ChampNames := []string{top3ChampStats[0].Name, top3ChampStats[1].Name, top3ChampStats[2].Name}
			//masteryStatsFormatted := formatMasteries(masteryStats)
			//fmt.Println(matchStatsFormatted)
			fileName, rankedType := getRankedAsset(RankedInfo(rankedInfo))
			description := formatPlayerRankedStats(rankedInfo)
			embed := formatRankedEmbed(rankedInfo, rankedType, fileName, description)
			embed = formatEmbedAuthor(embed, playermatchstats)
			embed = formatPlayerLookupEmbedFields(embed, playermatchstats, top3ChampNames)
			files := formatEmbedImages(embed, top3ChampNames, fileName)
			send = createMessageSend(embed, files)
			return send, err
		}
	}
	log.Println("Unable to get accInfo for: " + playerName)
	return send, err
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

	champ1 := "\u200b"
	champ2 := "\u200b"
	champ3 := "\u200b"
	champ1Name := "\u200b"
	champ2Name := "\u200b"
	champ3Name := "\u200b"

	var champ1PH int
	var champ2PH int
	var champ3PH int
	for j := 0; j < len(playerMatchStats.PlayerChampions); j++ {
		if top3Champs[0] == playerMatchStats.PlayerChampions[j].Name {
			champ1PH = j
		} else if top3Champs[1] == playerMatchStats.PlayerChampions[j].Name {
			champ2PH = j
		} else if top3Champs[2] == playerMatchStats.PlayerChampions[j].Name {
			champ3PH = j
		}
	}

	if top3Champs[0] != "" {
		champ1 = fmt.Sprintf("WR:%d%%\n (%dW/%dL)", ((playerMatchStats.PlayerChampions[champ1PH].Wins * 100) / playerMatchStats.PlayerChampions[champ1PH].GamesPlayed),
			playerMatchStats.PlayerChampions[champ1PH].Wins, playerMatchStats.PlayerChampions[champ1PH].Loss)
		champ1Name = playerMatchStats.PlayerChampions[champ1PH].Name
	}
	if top3Champs[1] != "" {
		champ2 = fmt.Sprintf("WR:%d%%\n (%dW/%dL)", ((playerMatchStats.PlayerChampions[champ2PH].Wins * 100) / playerMatchStats.PlayerChampions[champ2PH].GamesPlayed),
			playerMatchStats.PlayerChampions[champ2PH].Wins, playerMatchStats.PlayerChampions[champ2PH].Loss)
		champ2Name = playerMatchStats.PlayerChampions[champ2PH].Name
	}
	if top3Champs[2] != "" {
		champ3 = fmt.Sprintf("WR:%d%%\n (%dW/%dL)", ((playerMatchStats.PlayerChampions[champ3PH].Wins * 100) / playerMatchStats.PlayerChampions[champ3PH].GamesPlayed),
			playerMatchStats.PlayerChampions[champ3PH].Wins, playerMatchStats.PlayerChampions[champ3PH].Loss)
		champ3Name = playerMatchStats.PlayerChampions[champ3PH].Name
	}
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Place holder\t\t",
			Value: "Place holder\t\t",
		},
		{
			Name:   "\u200b",
			Value:  "\u200b",
			Inline: false,
		},
		{
			Name:   "Primary Role:",
			Value:  playerMatchStats.Role.PreferredRole1, // % picked?
			Inline: true,
		},
		{
			Name:   "Secondary Role:",
			Value:  playerMatchStats.Role.PreferredRole2,
			Inline: true,
		},
		{
			Name:   "\u200b",
			Value:  "\u200b",
			Inline: false,
		},

		{
			Name:   champ1Name,
			Value:  champ1,
			Inline: true,
		},
		{
			Name:   champ2Name,
			Value:  champ2,
			Inline: true,
		},
		{
			Name:   champ3Name,
			Value:  champ3,
			Inline: true,
		},
	}
	embed.Image = &discordgo.MessageEmbedImage{
		URL: "attachment://output.png",
	}
	return embed
}

func formatRankedEmbed(rankedInfo RankedInfo, rankedType int, fileName string, description string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:       000127255,
		Title:       rankedInfo[rankedType].SummonerName + "'s Ranked history",
		Description: description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "attachment://" + fileName,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return embed
}

func formatEmbedAuthor(embed *discordgo.MessageEmbed, playerMatchStats PlayerMatchStats) *discordgo.MessageEmbed {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    playerMatchStats.SummonerName,
		IconURL: "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/profileicon/" + strconv.Itoa(playerMatchStats.ProfileIcon) + ".png",
		URL:     "https://na.op.gg/summoner/userName=" + playerMatchStats.SummonerName,
	}
	return embed
}

func formatEmbedImages(embed *discordgo.MessageEmbed, imageNames []string, rankFileName string) []*discordgo.File {
	URL := "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/"

	for k := 0; k < len(imageNames); k++ {
		if imageNames[k] == "" {
			imageNames[k] = "ph"
		}
	}
	imageNames[0] += ".png"
	imageNames[1] += ".png"
	imageNames[2] += ".png"

	err := downloadFile(URL+imageNames[0], imageNames[0])
	if err != nil {
		log.Fatal("Unable to download file")
	}
	err = downloadFile(URL+imageNames[1], imageNames[1])
	if err != nil {
		log.Fatal("Unable to download file")
	}
	err = downloadFile(URL+imageNames[2], imageNames[2])
	if err != nil {
		log.Fatal("Unable to download file")
	}
	imageNames[0] = "./championImages/" + imageNames[0]
	imageNames[1] = "./championImages/" + imageNames[1]
	imageNames[2] = "./championImages/" + imageNames[2]

	fileImageName := mergeImages(imageNames)

	file, _ := os.Open(fileImageName)
	file2, _ := os.Open("./assets/" + rankFileName)

	var files []*discordgo.File
	files = append(files, &discordgo.File{
		Name:        rankFileName,
		ContentType: "image/png",
		Reader:      file2,
	})
	files = append(files, &discordgo.File{
		Name:        fileImageName,
		ContentType: "image/png",
		Reader:      file,
	})

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
				Role:        participant.IndividualPosition,
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
					playermatchstats.PlayerChampions[k].Role = participant.IndividualPosition
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
						Role:        participant.IndividualPosition,
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
	var pHolder int
	for j := 1; j < len(playerRoles.RoleCount); j++ {
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
			return "SUPPORT"
		}
	}
	return "UNKOWN"
}

func getTop3Champions(playerMatchStats PlayerMatchStats) []*PlayerChampions {
	var playerChampions []*PlayerChampions
	if len(playerMatchStats.PlayerChampions) < 1 {
		return playerChampions
	}
	playerChampion := playerMatchStats.PlayerChampions[0]
	playerChampion2 := &PlayerChampions{
		GamesPlayed: 0,
	}
	playerChampion3 := &PlayerChampions{
		GamesPlayed: 0,
	}

	for k := 1; k < len(playerMatchStats.PlayerChampions); k++ {
		if playerChampion.GamesPlayed <= playerMatchStats.PlayerChampions[k].GamesPlayed {
			playerChampion = playerMatchStats.PlayerChampions[k]
		}
	}
	for k := 0; k < len(playerMatchStats.PlayerChampions); k++ {
		if playerMatchStats.PlayerChampions[k].Name != playerChampion.Name {
			if playerChampion2.GamesPlayed <= playerMatchStats.PlayerChampions[k].GamesPlayed {
				playerChampion2 = playerMatchStats.PlayerChampions[k]
			}
		}
	}
	for k := 0; k < len(playerMatchStats.PlayerChampions); k++ {
		if playerMatchStats.PlayerChampions[k].Name != playerChampion.Name && playerMatchStats.PlayerChampions[k].Name != playerChampion2.Name {
			if playerChampion3.GamesPlayed <= playerMatchStats.PlayerChampions[k].GamesPlayed {
				playerChampion3 = playerMatchStats.PlayerChampions[k]
			}
		}
	}

	playerChampions = append(playerChampions, playerChampion)
	playerChampions = append(playerChampions, playerChampion2)
	playerChampions = append(playerChampions, playerChampion3)
	return playerChampions
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
	sp := image.Point{(img[0].Bounds().Dx() - 20), 0}
	sp2 := image.Point{(img[1].Bounds().Dx() - 20), 0}

	r2 := image.Rectangle{sp, sp.Add(img[1].Bounds().Size())}

	sp3 := image.Point{sp.X + sp2.X, 0}

	r3 := image.Rectangle{sp3, sp3.Add(img[2].Bounds().Size())}
	r3.Max.X = 299
	r := image.Rectangle{image.Point{0, 0}, r3.Max}

	rgba := image.NewRGBA(r)
	draw.Draw(rgba, img[0].Bounds(), img[0], image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img[1], image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r3, img[2], image.Point{0, 0}, draw.Src)

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
				" with " + strconv.Itoa(rankedStats[n].LeaguePoints) + " LP. Season wins/loss: " + strconv.Itoa(rankedStats[n].Wins) + " wins and " + strconv.Itoa(rankedStats[n].Losses) + " losses." + strconv.Itoa((rankedStats[n].Wins*100)/(rankedStats[n].Wins+rankedStats[n].Losses)) + "%"
		} else {
			return "Currently unranked."
		}

	}
	return "```No ranked data found```"
}
