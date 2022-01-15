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

	fmt.Println("Bot is now running")
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
	if m.Content == "!live" {
		if len([]rune(playerName[1])) > 0 {
			fmt.Println(playerName[1])
			s.ChannelMessageSend(m.ChannelID, query.IsInGame(playerName[1]))
		}
		return
	}

	// !challenge <username>
	if strings.HasPrefix(m.Content, "!challenge") && len(strings.Split(m.Content, " ")) == 2 {
		//handleChallenge(s, m)
		return
	}

	// !accepct <username>
	if strings.HasPrefix(m.Content, "!accept") && len(strings.Split(m.Content, " ")) == 2 {
		//handleAcceptChallenge(s, m)
		return
	}

	// !help
	if m.Content == "!help" {
		query.HandleHelp(s, m)
		return
	}

	// help as default
	query.HandleHelp(s, m)
}

func validateName(name string) bool {
	if len([]rune(name)) > 0 {
		return true
	}
	return false
}
