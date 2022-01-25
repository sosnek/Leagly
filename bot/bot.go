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
	args := createName(strings.Fields(messageContent))

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
			send, err := query.LookupPlayer(args[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSendComplex(m.ChannelID, send)
			//query.DeleteImages(filesToDelete)
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
			embed := &discordgo.MessageEmbed{
				URL:         "https://www.youtube.com/",
				Color:       000255000,
				Title:       "Lets",
				Description: "Gold III with 68 LP. This season they have a total",
				Image: &discordgo.MessageEmbedImage{
					URL:    "attachment://output.png",
					Width:  64,
					Height: 16,
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Sosnek",
					IconURL: "https://i.imgur.com/AfFp7pu.png",
					URL:     "https://na.op.gg/summoner/userName=lets",
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL:    "attachment://Emblem_Challenger.png",
					Height: 32,
					Width:  32,
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
						Name:   "Akshan",
						Value:  "67%(2W/1L)",
						Inline: true,
					},
					{
						Name:   "Aurelion Sol",
						Value:  "67%(2W/1L)",
						Inline: true,
					},
					{
						Name:   "Leblanc",
						Value:  "67%(2W/1L)",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Some footer text here",
					IconURL: "https://i.imgur.com/AfFp7pu.png",
				},
			}
			//might have to do only 2 champions because :
			/*

				embed width will only be constrained to image width if the image is at least 300px wide
				otherwise the text will expand it, like youve done with your footer
				thumbnails are not factored into this at all, only the big image at the bottom

				embed image size: 400x300
				embed thumbnail size: 80x80
				if the embed image is at least 300 pixels wide after resizing, the embed size will shrink to the size of the image
				*from discohook discord
			*/

			file, _ := os.Open("./output.png")
			file2, _ := os.Open("./assets/Emblem_Challenger.png")

			var files []*discordgo.File
			files = append(files, &discordgo.File{
				Name:        "Emblem_Challenger.png",
				ContentType: "image/png",
				Reader:      file2,
			})
			files = append(files, &discordgo.File{
				Name:        "output.png",
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
}

func createName(args []string) []string {
	for n := 2; n < len(args); n++ {
		args[1] += " " + args[n]
	}
	return args
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
