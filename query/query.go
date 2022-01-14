package query

import 
(
   "fmt" //to print errors
   "encoding/json"
   "bytes"
   "log"
   "io/ioutil"
   "net/http"
   "strconv"
   "Leagly/config"
   "github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin . 
)

type Summoner struct {
	Id string
	AccountId string
	Puuid string
	Name string
	SummonerLevel int
}

type Info struct {
	GameDuration int
	GameMode string
	Participants Participants
}

type Participants struct {
	Assists int
	ChampionName string
	TotalDamageDealtToChampions int
	Deaths int
	GameEndedInSurrender bool
	Kills int
	TotalMinionsKilled int
	VisionScore int
	Win bool

}


func GetLastMatch(s *discordgo.Session, m *discordgo.MessageCreate, playerName string) {
	sb := GetAccountInfo(playerName)
	//error checking here?

	matchID := GetLastMatchID(sb.Puuid)
	GetMatch(matchID)

	s.ChannelMessageSend(m.ChannelID, playerName + " has an account level of: " + strconv.Itoa(sb.SummonerLevel))
}


func GetMatch(matchid string) {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/"+ matchid + "?api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}

	match := string(body)
	fmt.Printf(match)
}



func GetLastMatchID(puuid string) string {
	resp, err := http.Get("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/" + puuid + "/ids?start=0&count=1&api_key=" + config.ApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}

	var arr []string
    _ = json.Unmarshal([]byte(body), &arr)
fmt.Printf("\n")
	fmt.Printf("Unmarshaled: %v", arr)
	fmt.Printf("\n")
	return arr[0]
}


func GetAccountInfo(playerName string) Summoner {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + playerName + "?api_key=" + config.ApiKey)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
 	//Convert the body to type string
	sb := string(body)
	fmt.Printf(sb)


	var summoner Summoner	
	json.Unmarshal([]byte(sb), &summoner)

	return summoner
}

func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "commands:\n"
	msg = fmt.Sprintf("%s\t%s\n", msg, "!help - shows all available commands")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challengeBot - bot will always accept challenges")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challenges - shows all your open challenges")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!challenge <username> - challenge another user")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!accepct <username> - accept a challenge from another user")
	msg = fmt.Sprintf("%s\t%s\n", msg, "!leaderboard - shows the leaderboard")

	s.ChannelMessageSend(m.ChannelID, msg)
}





//HTTP Get example
func TestAPIGet(s *discordgo.Session, m *discordgo.MessageCreate) {

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
func TestAPIPost(s *discordgo.Session, m *discordgo.MessageCreate) {

	postBody, _ := json.Marshal(map[string]string{
		"name":  "Toby",
		"email": "Toby@example.com",
	 })
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