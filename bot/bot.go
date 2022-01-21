package bot

import (
	"Leagly/config" //importing our config package which we have created above
	"Leagly/query"
	"fmt" //to print errors
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var BotId string
var goBot *discordgo.Session

func ConnectToDiscord() {

	leaglyBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	leaglyBot.AddHandler(messageCreate)

	leaglyBot.Identify.Intents = discordgo.IntentsGuildMessages

	err = leaglyBot.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Leagly is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	leaglyBot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageContent := m.Content
	args := strings.Fields(messageContent)

	if len(args) < 1 {
		return
	}

	// !help
	if m.Content == "!help" {
		handleHelp(s, m)
		return
	}

	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// !lastmatch - Searches and displays stats from last league game played
	if args[0] == "!lastmatch" {
		if validateName(args) {
			s.ChannelMessageSend(m.ChannelID, query.GetLastMatch(args[1]))
		}
		return
	}

	// !live - checks if player is currently in a game
	if args[0] == "!live" {
		if validateName(args) {
			s.ChannelMessageSend(m.ChannelID, query.IsInGame(args[1]))
		}
		return
	}

	//lookup
	if args[0] == "!lookup" {
		if validateName(args) {
			s.ChannelMessageSendComplex(m.ChannelID, query.LookupPlayer(args[1]))
		}
		return
	}

	if args[0] == "!test" {
		if validateName(args) {
			s.ChannelMessageSend(m.ChannelID, query.GetChampion(args[1]))
		}
		return
	}

	if args[0] == "!g" {
		if validateName(args) {
			file, _ := os.Open("./assets/Emblem_Challenger.png")
			embed := &discordgo.MessageEmbed{
				Color:       000255000,
				Title:       "Lets",
				Description: "Gold III with 68 LP. This season they have a total of 12 wins and 33 losses",
				Image: &discordgo.MessageEmbedImage{
					URL:    "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/Akshan.png",
					Width:  32,
					Height: 32,
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Sosnek",
					IconURL: "https://i.imgur.com/AfFp7pu.png",
					URL:     "https://na.op.gg/summoner/userName=lets",
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "attachment://Emblem_Challenger.png",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "GOLD III\t\t",
						Value: "W/L Ratio : 52%\t\t",
					},
					{
						Name:   "Primary Role:",
						Value:  "ADC",
						Inline: true,
					},
					{
						Name:   "\u200b",
						Value:  "\u200b",
						Inline: false,
					},
					{
						Name:   "Akshan\t\t",
						Value:  "67 % (2W 1 L)\t\t",
						Inline: true,
					},
					{
						Name:   "Neeko\t\t",
						Value:  "67 % (2W 1 L)\t\t",
						Inline: true,
					},
					{
						Name:   "Leblanc\t\t",
						Value:  "67 % (2W 1 L)\t\t",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Some footer text here",
					IconURL: "https://i.imgur.com/AfFp7pu.png",
				},
			}

			var files []*discordgo.File
			files = append(files, &discordgo.File{
				Name:        "Emblem_Challenger.png",
				ContentType: "image/png",
				Reader:      file,
			})

			send := &discordgo.MessageSend{
				Embed: embed,
				Files: files,
			}
			s.ChannelMessageSendComplex(m.ChannelID, send)
		}
		return
	}

	if args[0] == "!em" {
		if validateName(args) {
			var tmp discordgo.MessageEmbed
			var tmp2 discordgo.MessageEmbedAuthor
			var tmp3 discordgo.MessageEmbedThumbnail
			var tmp4 discordgo.MessageEmbedField
			var tmp5 discordgo.MessageEmbedField
			var tmp11 discordgo.MessageEmbedField
			var tmp9 discordgo.MessageEmbedField
			var tmp10 discordgo.MessageEmbedField
			var tmp6 discordgo.MessageEmbedImage
			var tmp7 discordgo.MessageEmbedFooter
			var tmp12 discordgo.File
			var tmp13 discordgo.MessageAttachment
			tmp13.Filename = "test.jpng"
			tmp13.URL = "./bot/test.png"
			m.Attachments = append(m.Attachments, &tmp13)
			var tmp14 discordgo.MessageSend

			file, _ := os.Open("./bot/test.png")
			tmp12.Reader = file
			tmp12.Name = "test.png"
			tmp12.ContentType = "image/png"
			tmp14.File = &tmp12
			s.ChannelMessageSendComplex(m.ChannelID, &tmp14)
			//msg, err := s.ChannelFileSend(m.ChannelID, "test.jpg", file)
			//fmt.Println(msg, err)
			tmp7.Text = "Some footer text"
			tmp7.IconURL = "https://i.imgur.com/AfFp7pu.png"
			//tmp6.URL = "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/Akshan.png"
			tmp6.URL = "attachment://test.png"
			tmp5.Inline = true
			tmp5.Name = fmt.Sprintf("%-30s", "Neeko")
			tmp5.Value = fmt.Sprintf("%-30s", "77% (3W 4 L)")
			tmp9.Name = "10.7 / 8.0 / 17.9"
			tmp9.Value = "3:58:1 (53%)"
			tmp10.Name = "\u200B"
			tmp10.Value = "\u200B"
			tmp4.Inline = true
			tmp4.Name = fmt.Sprintf("%-30s", "Akshan")
			tmp4.Value = fmt.Sprintf("%-30s", "67% (2W 1L)")
			tmp11.Name = fmt.Sprintf("%-30s", "Neeko")
			tmp11.Value = fmt.Sprintf("%-30s", "100% (1W 0L)")
			tmp11.Inline = true
			tmp3.URL = "attachment://test.png"
			tmp2.IconURL = "https://i.imgur.com/AfFp7pu.png"
			tmp2.URL = "https://na.op.gg/summoner/userName=lets"
			tmp2.Name = "Lets"
			tmp.Color = 000255000
			tmp.Description = "Gold III with 68 LP. This season they have a total of 12 wins and 33 losses"
			tmp.URL = "https://na.op.gg/summoner/userName=lets"
			tmp.Image = &tmp6
			tmp.Author = &tmp2
			tmp.Thumbnail = &tmp3
			var tmp8 []*discordgo.MessageEmbedField
			tmp8 = append(tmp8, &tmp9)
			tmp8 = append(tmp8, &tmp10)
			tmp8 = append(tmp8, &tmp4)
			tmp8 = append(tmp8, &tmp5)
			tmp8 = append(tmp8, &tmp11)
			tmp.Fields = tmp8
			s.ChannelMessageSendEmbed(m.ChannelID, &tmp)
			//s.ChannelMessageSend(m.ChannelID, "test")
		}
		return
	}

}

func validateName(name []string) bool {
	if len(name) < 2 {
		return false
	}
	return len([]rune(name[1])) > 0
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "```Commands:\n"
	msg = fmt.Sprintf("%s\t%s\n", msg, "!help - shows all available commands")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!live <playername> - Checks to see if the player is in a game")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!lastmatch <playername> - shows the players last match stats")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!lookup <playername> - shows ranked history + mastery stats of player```")

	s.ChannelMessageSend(m.ChannelID, msg)
}
