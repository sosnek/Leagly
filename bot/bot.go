package bot

import (
	"Leagly/config"
	"Leagly/guilds"
	"Leagly/query"
	"fmt" //to print errors
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var BotId string
var discordUser []*DiscordUser
var up_time = time.Time{}

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

	Initialize(leaglyBot)
	leaglyBot.AddHandler(messageCreate)
	leaglyBot.AddHandler(guildCreate)
	leaglyBot.AddHandler(guildDelete)

	leaglyBot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	err = leaglyBot.Open()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("Leagly is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	leaglyBot.Close()
}

///
///
///
func Initialize(s *discordgo.Session) {
	err := query.InitializedChampStruct()
	guilds.DB, err = guilds.SetupDB()
	if err != nil {
		panic(err)
	}
	query.CreateChampionRatesFile()
	query.InitializeEmojis(s)
	up_time = time.Now()
	query.Version = query.GetLeagueVersion()
	go query.UpdateVersionAsync(s)
	go heartBeat(s)
	s.UpdateGameStatus(0, ">>help | @Leagly")
}

func heartBeat(s *discordgo.Session) {
	uptimeTicker := time.NewTicker(60 * time.Second)
	counter := 0
	for {
		select {
		case <-uptimeTicker.C:
			counter++
			s.ChannelMessageSend("962149630815137832", fmt.Sprintf("```Heartbeat counter %d. time : %s```", counter, time.Now()))
		}
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	_, err := guilds.View(guilds.DB, event.ID)
	if err != nil { //expect an error if the guild already exists in db
		err = guilds.Add(guilds.DB, event.ID, guilds.DiscordGuild{ID: event.ID, Region: "NA1", Region2: "americas", Prefix: ">>", Members: event.Guild.MemberCount})
		if err == nil {
			log.Println("Added guild ID:" + event.Guild.ID + ". Name: " + event.Guild.Name + " Num of users in guild: " + strconv.Itoa(event.Guild.MemberCount))
		} else {
			log.Println(err)
		}
	}
}

func guildDelete(bot *discordgo.Session, event *discordgo.GuildDelete) {
	if event.Unavailable {
		return
	}

	_, err := guilds.View(guilds.DB, event.ID)
	if err != nil {
		log.Println(err)
	} else {
		err = guilds.Delete(guilds.DB, event.ID)
		if err == nil {
			log.Println("Guild ID:" + event.Guild.ID + " Has removed Leagly")
		} else {
			log.Println(err)
		}
	}
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

	if !containsValidCommand(messageContent) {
		return //don't exhaust db search if message doesnt involve valid command
	}

	guild, err := guilds.View(guilds.DB, m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}
	prefix := strings.ToLower(guild.Prefix)
	command := strings.ToLower(args[0])

	if command == prefix+"help" {
		handleHelp(s, m, guild)
		return
	}

	if command == prefix+"region" {
		changeRegion(s, m, args, guild)
		return
	}

	if command == prefix+"live" {
		live(s, m, args, guild)
		return
	}

	if command == prefix+"lastmatch" {
		lastmatch(s, m, args, guild)
		return
	}

	if command == prefix+"lookup" {
		lookup(s, m, args, guild)
		return
	}

	if command == prefix+"mastery" {
		mastery(s, m, args, guild)
		return
	}

	if command == prefix+"prefix" {
		changePrefix(s, m, args, guild)
		return
	}

	if command == prefix+"uptime" {
		uptime(s, m, args, guild)
		return
	}

	if command == prefix+"gc" {
		getGuildCount(s, m)
	}

	if command == prefix+"feedback" {
		feedback(s, m, args, guild)
	}

	if command == prefix+"status" {
		status(s, m, args, guild)
	}

	if command == prefix+"patchnotes" {
		patchNotes(s, m, args, guild)
	}

	for _, v := range m.Mentions {
		if v.ID == s.State.User.ID {
			handleHelp(s, m, guild)
		}
	}
}

func containsValidCommand(msg string) bool {
	roles := []string{"help", "region", "live", "lastmatch", "lookup", "mastery", "prefix", "uptime", "gc", "feedback", "status", "patchnotes"}
	for i := 0; i < len(roles); i++ {
		if strings.Contains(msg, roles[i]) {
			return true
		}
	}
	return false
}
