package bot

import (
	"Leagly/config" //importing our config package which we have created above
	"Leagly/guilds"
	"Leagly/query"
	"fmt" //to print errors
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var BotId string
var goBot *discordgo.Session
var discordUser []*DiscordUser

type DiscordUser struct {
	ID        string
	timestamp time.Time
}

func ConnectToDiscord() {

	leaglyBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	leaglyBot.AddHandler(messageCreate)

	leaglyBot.AddHandler(guildCreate)

	leaglyBot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

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
	emoji9, _ := s.GuildEmojis("946539173597302804")
	emojis = append(emojis, emoji, emoji2, emoji3, emoji4, emoji5, emoji6, emoji7, emoji8, emoji9)
	query.InitEmojis(emojis)
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}
	log.Println(event.Guild.Name)

	guilds.DiscordGuilds = append(guilds.DiscordGuilds, &guilds.DiscordGuild{ID: event.ID, Prefix: "NA"})
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageContent := m.Content

	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}
	//combine names with spaces
	args := createName(strings.Fields(messageContent))
	if len(args) < 1 {
		return
	}

	runes := []rune(args[0])
	if string(runes[0:2]) != config.BotPrefix {
		return
	}

	command := strings.ToLower(args[0])

	// !help
	if command == config.BotPrefix+"help" {
		handleHelp(s, m)
		return
	}

	// !prefix
	if command == config.BotPrefix+"region" {
		changePrefix(s, m, args)
		return
	}

	// !live - checks if player is currently in a game
	if command == config.BotPrefix+"live" {
		live(s, m, args)
		return
	}

	// !lastmatch - Searches and displays stats from last league game played
	if command == config.BotPrefix+"lastmatch" {
		lastmatch(s, m, args)
		return
	}

	if command == config.BotPrefix+"lookup" {
		lookup(s, m, args)
		return
	}

	if command == config.BotPrefix+"mastery" {
		mastery(s, m, args)
		return
	}
}
