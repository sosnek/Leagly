package bot

import (
	"Leagly/config" //importing our config package which we have created above
	"Leagly/query"
	"fmt" //to print errors
	"log"
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
		log.Println(err)
		panic(err)
	}

	leaglyBot.AddHandler(messageCreate)

	leaglyBot.Identify.Intents = discordgo.IntentsGuildMessages

	err = leaglyBot.Open()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	Initialize(leaglyBot)
	fmt.Println("Leagly is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	leaglyBot.Close()
}

///
///
///
func Initialize(s *discordgo.Session) {
	query.InitializedChampStruct()
	InitializeEmojis(s)
}

///
///
///
func InitializeEmojis(s *discordgo.Session) {
	var emojis [][]*discordgo.Emoji
	emoji, _ := s.GuildEmojis("937465588446539920")
	emoji2, _ := s.GuildEmojis("937453232517693502")
	emoji3, _ := s.GuildEmojis("937481122198200320")
	emoji4, _ := s.GuildEmojis("937537071902503005")
	emoji5, _ := s.GuildEmojis("937482778499485756")
	emoji6, _ := s.GuildEmojis("938569984748163112")
	emoji7, _ := s.GuildEmojis("938569677326671913")
	emoji8, _ := s.GuildEmojis("938569400724910110")
	emojis = append(emojis, emoji)
	emojis = append(emojis, emoji2)
	emojis = append(emojis, emoji3)
	emojis = append(emojis, emoji4)
	emojis = append(emojis, emoji5)
	emojis = append(emojis, emoji6)
	emojis = append(emojis, emoji7)
	emojis = append(emojis, emoji8)
	query.InitEmojis(emojis)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageContent := m.Content

	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := createName(strings.Fields(messageContent))
	if len(args) < 1 {
		return
	}

	if messageContent[0:2] != config.BotPrefix {
		return
	}
	command := args[0]
	// !help
	if command == config.BotPrefix+"help" {
		log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + ": " + config.BotPrefix + "help")
		handleHelp(s, m)
		return
	}

	// if args[0] == "!test" {
	// 	s.ChannelMessageSend(m.ChannelID, "<:"+args[1]+":"+query.GetEmoji(args[1])+">")
	// 	return
	// }

	// !lastmatch - Searches and displays stats from last league game played
	if command == config.BotPrefix+"lastmatch" {
		if validateName(args) {

			log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "lastmatch " + args[1])
			send, err := query.GetLastMatch(args[1])
			if err != nil {
				log.Println("Discord ID: " + m.GuildID + "  " + err.Error())
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSendComplex(m.ChannelID, send)
		}
		return
	}

	// !live - checks if player is currently in a game
	if command == config.BotPrefix+"live" {
		if validateName(args) {
			log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "live " + args[1])
			send, err := query.IsInGame(args[1])
			if err != nil {
				log.Println("Discord ID: " + m.GuildID + "  " + err.Error())
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSendComplex(m.ChannelID, send)
		}
		return
	}

	//lookup
	if command == config.BotPrefix+"lookup" {
		if validateName(args) {
			log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "lookup " + args[1])
			send, err := query.LookupPlayer(args[1])
			if err != nil {
				log.Println("Discord ID: " + m.GuildID + "  " + err.Error())
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSendComplex(m.ChannelID, send)
			//query.DeleteImages(filesToDelete)
		}
		return
	}

	// if command == config.BotPrefix+"mastery" {
	// 	if validateName(args) {
	// 		log.Println(m.Author.Username + " : >>mastery " + args[1])
	// 		send, err := query.MasteryPlayer(args[1])
	// 		if err != nil {
	// 			log.Println(err)
	// 			s.ChannelMessageSend(m.ChannelID, err.Error())
	// 		}
	// 		s.ChannelMessageSendComplex(m.ChannelID, send)
	// 	}
	// 	return
	// }
}

func createName(args []string) []string {
	for n := 2; n < len(args); n++ {
		args[1] += " " + args[n]
	}
	return args
}

func validateName(name []string) bool {
	if len(name) < 2 {
		log.Println("Name not found in args list")
		return false
	}
	return len([]rune(name[1])) > 0
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "```Commands:\n"
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>help - shows all available commands")
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>live <playername> - Checks to see if the player is in a game")
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>lastmatch <playername> - shows the players last match stats")
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>lookup <playername> - shows ranked history of player```")
	//msg = fmt.Sprintf("%s\t%s\n", msg, "!mastery <playername> - shows mastery stats of player```")
	s.ChannelMessageSend(m.ChannelID, msg)
}
