package guilds

import (
	"log"
)

var DiscordGuilds []*DiscordGuild

type DiscordGuild struct {
	ID      string
	Prefix  string
	Prefix2 string //americas, europe, asia
}

func GetGuildPrefix(guildID string) string {
	for _, v := range DiscordGuilds {
		if v.ID == guildID {
			return v.Prefix
		}
	}
	log.Println("Could not find discord server ID: " + guildID + ". Defaulting to NA region")
	return "NA1"
}

func GetGuildPrefix2(guildID string) string {
	for _, v := range DiscordGuilds {
		if v.ID == guildID {
			return v.Prefix2
		}
	}
	log.Println("Could not find discord server ID: " + guildID + ". Defaulting to NA region")
	return "americas"
}
