package bot

import (
	"Leagly/config"
	"Leagly/guilds"
	"Leagly/query"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func live(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "live " + args[1])
		if onCoolDown(m.Author.ID, 3) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.IsInGame(args[1], guilds.GetGuildRegion(m.GuildID))
		if err != nil {
			log.Println("Discord server ID: " + m.GuildID + "  " + err.Error())
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m)
	}
}

func lastmatch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "lastmatch " + args[1])
		send, err := query.GetLastMatch(args[1], guilds.GetGuildRegion(m.GuildID), guilds.GetGuildRegion2(m.GuildID))
		if err != nil {
			log.Println("Discord server ID: " + m.GuildID + "  " + err.Error())
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m)
	}
}

func lookup(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "lookup " + args[1])
		if onCoolDown(m.Author.ID, 5) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.LookupPlayer(args[1], guilds.GetGuildRegion(m.GuildID), guilds.GetGuildRegion2(m.GuildID))
		if err != nil {
			log.Println("Discord server ID: " + m.GuildID + "  " + err.Error())
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m)
	}
}

func mastery(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "mastery " + args[1])
		if onCoolDown(m.Author.ID, 3) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.MasteryPlayer(args[1], guilds.GetGuildRegion(m.GuildID))
		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m)
	}
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelTyping(m.ChannelID)
	log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "help")

	s.ChannelMessageSendComplex(m.ChannelID, query.Help(guilds.GetGuildRegion(m.GuildID)))
}

func changeRegion(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if validateName(args) {
		if isValidRegion(args[1]) {
			s.ChannelTyping(m.ChannelID)
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + config.BotPrefix + "region " + args[1])
			for _, v := range guilds.DiscordGuilds {
				if v.ID == m.GuildID {
					v.Region = strings.ToUpper(args[1])
					if v.Region == "BR1" || v.Region == "NA1" || v.Region == "LA1" || v.Region == "LA2" {
						v.Region2 = "americas"
					} else if v.Region == "JP1" || v.Region == "OCE" || v.Region == "KR" {
						v.Region2 = "asia"
					} else {
						v.Region2 = "europe"
					}
					log.Println("Discord server ID: " + m.GuildID + "  Changed region to " + v.Region + " " + v.Region2)
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Region has been changed to %s for your discord", v.Region))
				}
			}
		} else {
			s.ChannelMessageSend(m.ChannelID,
				"Invalid region provided. Valid regions are : BR1, EUN1, EUW1, JP1, KR, LA1, LA2, NA1, OC1, RU, TR1")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m)
	}
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

func isValidRegion(region string) bool {
	regions := [11]string{"BR1", "EUN1", "EUW1", "JP1", "KR", "LA1", "LA2", "NA1", "OC1", "RU", "TR1"}
	for i := 0; i < len(regions); i++ {
		if regions[i] == strings.ToUpper(region) {
			return true
		}
	}
	return false
}

///Some summoner names can have spaces in them
/// This method will combine each name piece into a whole string
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
