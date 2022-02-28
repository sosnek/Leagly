package guilds

import (
	"log"
)

var DiscordGuilds []*DiscordGuild

type DiscordGuild struct {
	ID      string
	Region  string
	Region2 string //americas, europe, asia
}

func GetGuildRegion(guildID string) string {
	for _, v := range DiscordGuilds {
		if v.ID == guildID {
			return v.Region
		}
	}
	log.Println("Could not find discord server ID: " + guildID + ". Defaulting to NA region")
	return "NA1"
}

func GetGuildRegion2(guildID string) string {
	for _, v := range DiscordGuilds {
		if v.ID == guildID {
			return v.Region2
		}
	}
	log.Println("Could not find discord server ID: " + guildID + ". Defaulting to NA region")
	return "americas"
}
