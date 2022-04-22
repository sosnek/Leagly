package bot

import (
	"Leagly/config" //importing our config package which we have created above
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
	leaglyBot.AddHandler(messageCreate)
	leaglyBot.AddHandler(guildCreate)
	leaglyBot.AddHandler(guildDelete)

	leaglyBot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	err = leaglyBot.Open()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	Initialize(leaglyBot)
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
	if err != nil {
		panic(err)
	}
	query.CreateChampionRatesFile()
	InitializeEmojis(s)
	up_time = time.Now()
	query.Version = query.GetLeagueVersion()
	go query.UpdateVersionAsync(s)
	go heartBeat(s)
	s.UpdateGameStatus(0, ">>help | @Leagly")
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

	var exists bool
	for i := 0; i < len(guilds.DiscordGuilds); i++ {
		if event.ID == guilds.DiscordGuilds[i].ID {
			//log.Println("Guild ID:" + event.Guild.ID + ". Name: " + event.Guild.Name + " Already exists. Num of users in guild: " + strconv.Itoa(event.Guild.MemberCount))
			exists = true
		}
	}
	if !exists {
		guilds.DiscordGuilds = append(guilds.DiscordGuilds, &guilds.DiscordGuild{ID: event.ID, Region: "NA1", Region2: "americas", Prefix: ">>"})
		log.Println("Added guild ID:" + event.Guild.ID + ". Name: " + event.Guild.Name + " Num of users in guild: " + strconv.Itoa(event.Guild.MemberCount))
	}
}

func guildDelete(bot *discordgo.Session, event *discordgo.GuildDelete) {
	if event.Unavailable {
		return
	}

	for i := 0; i < len(guilds.DiscordGuilds); i++ {
		if event.ID == guilds.DiscordGuilds[i].ID {
			log.Println("Guild ID:" + event.Guild.ID + ". Name: " + event.Guild.Name + " has removed Leagly.")
			guilds.DiscordGuilds[i] = guilds.DiscordGuilds[len(guilds.DiscordGuilds)-1]
			guilds.DiscordGuilds = guilds.DiscordGuilds[:len(guilds.DiscordGuilds)-1]
			return
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
	guild := guilds.GetGuild(m.GuildID)
	prefix := strings.ToLower(guild.Prefix)

	command := strings.ToLower(args[0])
	// !help
	if command == prefix+"help" {
		handleHelp(s, m)
		return
	}

	// !prefix
	if command == prefix+"region" {
		changeRegion(s, m, args)
		return
	}

	// !live - checks if player is currently in a game
	if command == prefix+"live" {
		live(s, m, args)
		return
	}

	// !lastmatch - Searches and displays stats from last league game played
	if command == prefix+"lastmatch" {
		lastmatch(s, m, args)
		return
	}

	if command == prefix+"lookup" {
		lookup(s, m, args)
		return
	}

	if command == prefix+"mastery" {
		mastery(s, m, args)
		return
	}

	if command == prefix+"prefix" {
		changePrefix(s, m, args)
		return
	}

	if command == prefix+"uptime" {
		uptime(s, m, args)
		return
	}

	if command == prefix+"gc" {
		getGuildCount(s, m)
	}

	if command == prefix+"feedback" {
		feedback(s, m, args)
	}

	if command == prefix+"status" {
		status(s, m, args)
	}

	if command == prefix+"patchnotes" {
		patchNotes(s, m, args)
	}

	for _, v := range m.Mentions {
		if v.ID == s.State.User.ID {
			handleHelp(s, m)
		}
	}
}
