package tests

import (
	"Leagly/query"
	"fmt"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
)

type TestCase struct {
	value     string
	expected1 string
	expected  *discordgo.MessageSend
	actual    *discordgo.MessageSend
}

var (
	dg    *discordgo.Session // Stores a global discordgo user session
	dgBot *discordgo.Session // Stores a global discordgo bot session

	envOAuth2Token = os.Getenv("DG_OAUTH2_TOKEN") // Token to use when authenticating using OAuth2 token
	envBotToken    = os.Getenv("DGB_TOKEN")       // Token to use when authenticating the bot account
	envGuild       = os.Getenv("DG_GUILD")        // Guild ID to use for tests
	envChannel     = os.Getenv("DG_CHANNEL")      // Channel ID to use for tests
	envAdmin       = os.Getenv("DG_ADMIN")        // User ID of admin user to use for tests
)

// temporary tests
func TestMain(m *testing.M) {
	fmt.Println("Init is being called.")
	if envBotToken != "" {
		if d, err := discordgo.New(envBotToken); err == nil {
			dgBot = d
		}
	}

	if envOAuth2Token == "" {
		envOAuth2Token = os.Getenv("DGU_TOKEN")
	}

	if envOAuth2Token != "" {
		if d, err := discordgo.New(envOAuth2Token); err == nil {
			dg = d
		}
	}

	os.Exit(m.Run())
}

func TestLiveCommand(t *testing.T) {
	t.Run("Should not find data for non-existing profile testaccountbcz11", func(t *testing.T) {
		testCase := TestCase{
			value:    "testaccountbcz11",
			expected: nil,
		}
		testCase.actual, _ = query.IsInGame(testCase.value, "na1")
		if testCase.actual != testCase.expected {
			t.Fail()
		}
	})
}

func TestLastMatchCommand(t *testing.T) {
	t.Run("Test last match data from my inactive smurf", func(t *testing.T) {
		testCase := TestCase{
			value:     "testaccountbcz11",
			expected1: "> __<:Talon:937537143281164298>ruń__",
		}
		testCase.actual, _ = query.GetLastMatch(testCase.value, "na1", "americas")
		if testCase.actual.Embed.Fields[16].Name != testCase.expected1 {
			t.Fail()
		}
	})
}

func TestGettingEmojisFromMessage(t *testing.T) {
	// msg := "test test <:kitty14:811736565172011058> <:kitty4:811736468812595260>"
	// m := &discordgo.Message{
	// 	Content: msg,
	// }
	// emojis := m.GetCustomEmojis()
	// if len(emojis) < 1 {
	// 	t.Error("No emojis found.")
	// 	return
	// }

}
