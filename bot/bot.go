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

func Initialize(s *discordgo.Session) {
	query.InitializedChampStruct()
	InitializeEmojis(s)
}

func InitializeEmojis(s *discordgo.Session) {
	var emojis [][]*discordgo.Emoji
	emoji, _ := s.GuildEmojis("937465588446539920")
	emoji2, _ := s.GuildEmojis("937453232517693502")
	emoji3, _ := s.GuildEmojis("937481122198200320")
	emoji4, _ := s.GuildEmojis("937537071902503005")
	emoji5, _ := s.GuildEmojis("937482778499485756")
	emojis = append(emojis, emoji)
	emojis = append(emojis, emoji2)
	emojis = append(emojis, emoji3)
	emojis = append(emojis, emoji4)
	emojis = append(emojis, emoji5)
	query.InitEmojis(emojis)
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

	if args[0] == "!test" {
		s.ChannelMessageSend(m.ChannelID, "<:"+args[1]+":"+query.GetEmoji(args[1])+">")
		return
	}
	/* TODO : missing icons from these champs
	<:Viego:>
	<:Sett:>
	<:Senna:>
	<:MonkeyKing:> (Wukong)
	<:Samira:>
	<:Qiyana:>
	<:Velkoz:>
	<:Akshan:>
	<:Gwen:>
	<:Vex:>
	<:Aphelios:>
	<:Seraphine:>
	<:Rell:>
	*/

	if args[0] == "!test2" {
		query.Temp(s, m)
		return
	}

	// !lastmatch - Searches and displays stats from last league game played
	if args[0] == "!lastmatch" {
		if validateName(args) {

			log.Println(m.Author.Username + " : !lastmatch " + args[1])
			send, err := query.GetLastMatch(args[1])
			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSendComplex(m.ChannelID, send)
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
		log.Println("Name not found in args list")
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
