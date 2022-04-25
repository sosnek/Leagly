package bot

import (
	"Leagly/config"
	"Leagly/guilds"
	"Leagly/query"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func live(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "live " + args[1])
		if onCoolDown(m.Author.ID, 3) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.IsInGame(args[1], guild.Region)
		if err != nil {
			log.Println("Error: Discord server ID: " + m.GuildID + "  " + err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m, guild)
	}
}

func lastmatch(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "lastmatch " + args[1])
		send, err := query.GetLastMatch(args[1], guild.Region, guild.Region2)
		if err != nil {
			log.Println("Error: Discord server ID: " + m.GuildID + "  " + err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m, guild)
	}
}

func lookup(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "lookup " + args[1])
		if onCoolDown(m.Author.ID, 5) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.LookupPlayer(args[1], guild.Region, guild.Region2)
		if err != nil {
			log.Println("Error: Discord server ID: " + m.GuildID + "  " + err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m, guild)
	}
}

func mastery(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if validateName(args) {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "mastery " + args[1])
		if onCoolDown(m.Author.ID, 3) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		}
		send, err := query.MasteryPlayer(args[1], guild.Region)
		if err != nil {
			log.Println("Error: Discord server ID: " + m.GuildID + "  " + err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m, guild)
	}
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate, guild guilds.DiscordGuild) {
	s.ChannelTyping(m.ChannelID)
	log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "help")
	s.ChannelMessageSendComplex(m.ChannelID, query.Help(guild.Region, guild.Prefix))
}

func changeRegion(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if validateName(args) {
		if isValidRegion(args[1]) {
			s.ChannelTyping(m.ChannelID)
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "region " + args[1])
			if guild.ID == m.GuildID {
				guild.Region = strings.ToUpper(args[1])
				if guild.Region == "BR1" || guild.Region == "NA1" || guild.Region == "LA1" || guild.Region == "LA2" || guild.Region == "OC1" {
					guild.Region2 = "americas"
				} else if guild.Region == "JP1" || guild.Region == "KR" {
					guild.Region2 = "asia"
				} else {
					guild.Region2 = "europe"
				}
				err := guilds.Update(guilds.DB, guild.ID, guild)
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("Discord server ID: " + m.GuildID + "  Changed region to " + guild.Region + " " + guild.Region2)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Region has been changed to %s for your discord", guild.Region))
			}
		} else {
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " Invalid region")
			s.ChannelMessageSend(m.ChannelID,
				"Invalid region provided. Valid regions are : BR1, EUN1, EUW1, JP1, KR, LA1, LA2, NA1, OC1, RU, TR1")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m, guild)
	}
}

func changePrefix(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	isAdmin, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		log.Println("Discord server ID: " + m.GuildID + " Error getting channel permissions")
		return
	}
	if isAdmin&discordgo.PermissionAdministrator < 1 {
		s.ChannelMessageSend(m.ChannelID, "This is an Admin only command")
		log.Println("Discord server ID: " + m.GuildID + " User does not have channel admin controls. " + m.Author.Username)
		return
	}
	if validateName(args) {
		if len(args[1]) < 10 {
			guild.Prefix = args[1]
			err = guilds.Update(guilds.DB, guild.ID, guild)
			if err != nil {
				log.Println(err)
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Prefix has been changed to "+guild.Prefix)
			log.Println("Discord server ID: " + m.GuildID + "  Changed prefix to " + guild.Prefix)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Prefix must be under 10 characters.")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please follow the command format!")
		handleHelp(s, m, guild)
	}
}

func uptime(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if len(args) < 2 {
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "uptime")
		s.ChannelMessageSendComplex(m.ChannelID, query.UpTime(up_time))
	}
}

func status(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if len(args) < 2 {
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "status")
		s.ChannelMessageSendComplex(m.ChannelID, query.RiotApiStatus(guild.Region))
	}

}

func getGuildCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	if guilds.HasDebugPermissions(m.Author.ID) {
		s.ChannelMessageSendComplex(m.ChannelID, query.GuildCount(guilds.GetGuildCount()))
	}
}

func feedback(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if validateName(args) {
		if onCoolDown(m.Author.ID, 30) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " on cooldown")
			return
		} else {
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "feedback")
			s.ChannelMessageSend("955121671105286175", fmt.Sprintf("From %s, Feedback: %s ", m.Author.Username, args[1]))
			s.ChannelMessageSend(m.ChannelID, "Message has been saved! Thank you for the feedback. :)")
		}
	}
}

func patchNotes(s *discordgo.Session, m *discordgo.MessageCreate, args []string, guild guilds.DiscordGuild) {
	if len(args) == 1 {
		s.ChannelTyping(m.ChannelID)
		log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "patchnotes")
		send, err := query.PatchNotes()
		if err != nil {
			log.Println("Error: Discord server ID: " + m.GuildID + "  " + err.Error())
		}
		s.ChannelMessageSendComplex(m.ChannelID, send)
	} else if len(args) == 2 {
		if args[1] == "toggle" {
			guild.AutoPatchNotes = !guild.AutoPatchNotes

			PatchNotesCh, err := query.Encrypt([]byte(m.ChannelID), []byte(config.EncryptionKey))
			if err != nil {
				log.Println(err)
				return
			}
			guild.PatchNotesCh = PatchNotesCh

			err = guilds.Update(guilds.DB, guild.ID, guild)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Discord server ID: " + m.GuildID + "  " + m.Author.Username + " : " + guild.Prefix + "patchnotes toggle " + strconv.FormatBool(guild.AutoPatchNotes))
			if guild.AutoPatchNotes {
				s.ChannelMessageSend(m.ChannelID, "Auto patch notes have been enabled")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Auto patch notes have been disabled")
			}
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Incorrect format. Patchnotes example: ```>>Patchnotes toggle```")
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
