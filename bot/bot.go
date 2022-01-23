package bot

import (
	"Leagly/config" //importing our config package which we have created above
	"Leagly/query"
	"fmt" //to print errors
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
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
			embed := &discordgo.MessageEmbed{
				URL:         "https://www.youtube.com/",
				Color:       000255000,
				Title:       "Lets",
				Description: "Gold III with 68 LP. This season they have a total of 12 wins and 33 losses",
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
						Name:   "```Akshanmnnnnnnnnnnnnn```",
						Value:  "```67 % (2W 1 L)```",
						Inline: true,
					},
					{
						Name:   "```Aurelion Sol```",
						Value:  "```67 % (2W 1 L)```",
						Inline: true,
					},
					{
						Name:   "```Leblanc```",
						Value:  "```67 % (2W 1 L)```",
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
			URL := "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/Teemo.png"
			URL2 := "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/Zoe.png"
			URL3 := "http://ddragon.leagueoflegends.com/cdn/12.2.1/img/champion/Ryze.png"
			err := downloadFile(URL, "Teemo.png")
			if err != nil {
				return
			}
			err = downloadFile(URL2, "Zoe.png")
			if err != nil {
				return
			}
			err = downloadFile(URL3, "Ryze.png")
			if err != nil {
				return
			}

			var imageNames []string
			imageNames = append(imageNames, "Zoe.png")
			imageNames = append(imageNames, "Teemo.png")
			imageNames = append(imageNames, "Ryze.png")

			fileImageName := mergeImages(imageNames)

			file, _ := os.Open(fileImageName)
			file2, _ := os.Open("./assets/Emblem_Challenger.png")

			var files []*discordgo.File
			files = append(files, &discordgo.File{
				Name:        "Emblem_Challenger.png",
				ContentType: "image/png",
				Reader:      file2,
			})
			files = append(files, &discordgo.File{
				Name:        fileImageName,
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
	sp := image.Point{img[0].Bounds().Dx(), 0}
	sp2 := image.Point{img[1].Bounds().Dx(), 0}

	r2 := image.Rectangle{sp, sp.Add(img[1].Bounds().Size())}

	sp3 := image.Point{sp.X + sp2.X, 0}
	r3 := image.Rectangle{sp3, sp3.Add(img[2].Bounds().Size())}

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

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return err
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
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
