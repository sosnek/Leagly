package bot

import 
(
   "fmt" //to print errors
   "Leagly/config" //importing our config package which we have created above
   "Leagly/query"
   "os"
   "syscall"
   "strings"
   "os/signal"
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

	// !challengeBot - bot will always accept challenges
	if playerName[0] == "!lastmatch" {
		if len([]rune(playerName[1])) > 0 {
			fmt.Println(playerName[1])
			query.GetLastMatch(s, m, playerName[1])
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

	// !httpPost
	if m.Content == "!test2" {
		query.TestAPIPost(s, m)
		return
	}

	// !httpGet
	if m.Content == "!test" {
		query.TestAPIGet(s, m)
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





