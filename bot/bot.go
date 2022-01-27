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
		log.Println(m.Author.Username + ": !help")
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
			log.Println(m.Author.Username + " : !lastmatch " + args[1])
			s.ChannelMessageSend(m.ChannelID, query.GetLastMatch(args[1]))
		}
		return
	}

	// !live - checks if player is currently in a game
	if args[0] == "!live" {
		if validateName(args) {
			log.Println(m.Author.Username + " : !live " + args[1])
			s.ChannelMessageSend(m.ChannelID, query.IsInGame(args[1]))
		}
		return
	}

	//lookup
	if args[0] == "!lookup" {
		if validateName(args) {
			log.Println(m.Author.Username + " : !lookup " + args[1])
			send, err := query.LookupPlayer(args[1])
			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSendComplex(m.ChannelID, send)
			//query.DeleteImages(filesToDelete)
		}
		return
	}

	if args[0] == "!mastery" {
		if validateName(args) {
			log.Println(m.Author.Username + " : !mastery " + args[1])
			send, err := query.MasteryPlayer(args[1])
			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
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
		log.Panicln("Name not found in args list")
		return false
	}
	return len([]rune(name[1])) > 0
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "```Commands:\n"
	msg = fmt.Sprintf("%s\t%s\n", msg, "!help - shows all available commands")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!live <playername> - Checks to see if the player is in a game")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!lastmatch <playername> - shows the players last match stats")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!lookup <playername> - shows ranked history of player```")
	//msg = fmt.Sprintf("%s\t%s\n", msg, "!mastery <playername> - shows mastery stats of player```")
	s.ChannelMessageSend(m.ChannelID, msg)
}
