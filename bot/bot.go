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
	playerName := strings.Fields(messageContent)
	fmt.Println(playerName, len(playerName))

	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// !lastmatch - Searches and displays stats from last league game played
	if playerName[0] == "!lastmatch" {
		if validateName(playerName[1]) {
			s.ChannelMessageSend(m.ChannelID, query.GetLastMatch(playerName[1]))
		}
		return
	}

	// !live - checks if player is currently in a game
	if playerName[0] == "!live" {
		if validateName(playerName[1]) {
			s.ChannelMessageSend(m.ChannelID, query.IsInGame(playerName[1]))
		}
		return
	}

	// !help
	if m.Content == "!help" {
		handleHelp(s, m)
		return
	}

	// help as default
	handleHelp(s, m)
}

func validateName(name string) bool {
	return len([]rune(name)) > 0
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "commands:\n"
	msg = fmt.Sprintf("%s\t%s\n", msg, "!help - shows all available commands")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challengeBot - bot will always accept challenges")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challenges - shows all your open challenges")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challenge <username> - challenge another user")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!accepct <username> - accept a challenge from another user")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!leaderboard - shows the leaderboard")

	s.ChannelMessageSend(m.ChannelID, msg)
}
