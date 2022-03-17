package guilds

import (
	"log"
)

var DiscordGuilds []*DiscordGuild
var MY_DISCORD_ID = "220732095083839488"

type DiscordGuild struct {
	ID      string
	Region  string
	Region2 string //americas, europe, asia
	Prefix  string
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

func GetGuildPrefix(guildID string) string {
	for _, v := range DiscordGuilds {
		if v.ID == guildID {
			return v.Prefix
		}
	}
	log.Println("Could not find discord server ID: " + guildID + ". Defaulting to >> prefix")
	return ">>"
}

func GetGuildCount() int {
	return len(DiscordGuilds)
}

func HasDebugPermissions(authID string) bool {
	return authID == MY_DISCORD_ID
}
