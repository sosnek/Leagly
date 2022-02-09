package bot

import (
	"Leagly/config"
	"Leagly/query"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func live(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "live " + args[1])
		send, err := query.IsInGame(args[1])
		if err != nil {
			log.Println("Discord ID: " + m.GuildID + "  " + err.Error())
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	}
}

func lastmatch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "lastmatch " + args[1])
		send, err := query.GetLastMatch(args[1])
		if err != nil {
			log.Println("Discord ID: " + m.GuildID + "  " + err.Error())
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	}
}

func lookup(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "lookup " + args[1])
		if onCoolDown(m.Author.ID, 5) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.LookupPlayer(args[1])
		if err != nil {
			log.Println("Discord ID: " + m.GuildID + "  " + err.Error())
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
		//query.DeleteImages(filesToDelete)
	}
}

func mastery(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		log.Println(m.Author.Username + " : >>mastery " + args[1])
		if onCoolDown(m.Author.ID, 3) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.MasteryPlayer(args[1])
		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	}
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "```Commands:\n"
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>help - shows all available commands")
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>live <playername> - Checks to see if the player is in a game")
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>lastmatch <playername> - shows the players last match stats")
	msg = fmt.Sprintf("%s\t%s\n", msg, ">>lookup <playername> - shows ranked history of player```")
	//msg = fmt.Sprintf("%s\t%s\n", msg, "!mastery <playername> - shows mastery stats of player```")
	s.ChannelMessageSend(m.ChannelID, msg)
}

func onCoolDown(user string, cd float64) float64 {
	for i := range discordUser {
		if discordUser[i].ID == user {
			t := time.Now()
			elapsed := t.Sub(discordUser[i].timestamp)
			if elapsed.Seconds() < cd {
				return elapsed.Seconds()
			} else {
				discordUser[i].timestamp = t
				return 0
			}
		}
	}
	discordUser = append(discordUser, &DiscordUser{ID: user, timestamp: time.Now()})
	return 0
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
