package bot

import 
(
   "fmt" //to print errors
   "Leagly/config" //importing our config package which we have created above
   "os"
   "os/signal"
   "strings"
   "syscall"
   "log"
   "io/ioutil"
   "net/http"
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
	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// !challengeBot - bot will always accept challenges
	if m.Content == "!challengeBot" {
		//handleChallengeBot(s, m)
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
		testAPIPost(s, m)
		return
	}

	// !httpGet
	if m.Content == "!test" {
		testAPIGet(s, m)
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

//HTTP Get example
func testAPIGet(s *discordgo.Session, m *discordgo.MessageCreate) {

	resp, err := http.Get("https://na1.api.riotgames.com/lol/status/v4/platform-data?api_key=RGAPI-894ce659-63e9-44f8-8297-41a17f3d95cd")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
 	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
	s.ChannelMessageSend(m.ChannelID, sb)

}


//HTTP Post example
func testAPIGet(s *discordgo.Session, m *discordgo.MessageCreate) {

	postBody, _ := json.Marshal(map[string]string{
		"name":  "Toby",
		"email": "Toby@example.com",
	 })
	 responseBody := bytes.NewBuffer(postBody)


	 responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
   resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)
	//Handle Error
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
	//Read the response body
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   sb := string(body)
   log.Printf(sb)

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