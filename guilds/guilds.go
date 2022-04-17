package guilds

var DiscordGuilds []*DiscordGuild
var MY_DISCORD_ID = "220732095083839488"

type DiscordGuild struct {
	ID             string
	Region         string
	Region2        string //americas, europe, asia
	Prefix         string
	AutoPatchNotes bool
	PatchNotesCh   string
}

func GetGuild(guildID string) *DiscordGuild {
	for _, v := range DiscordGuilds {
		if v.ID == guildID {
			return v
		}
	}
	return nil
}

func GuildsWithAutoPatchNotes() []string {
	var guildsWithAutoUpdates []string
	for _, v := range DiscordGuilds {
		if v.AutoPatchNotes {
			guildsWithAutoUpdates = append(guildsWithAutoUpdates, v.PatchNotesCh)
		}
	}
	return guildsWithAutoUpdates
}

func GetGuildCount() int {
	return len(DiscordGuilds)
}

func HasDebugPermissions(authID string) bool {
	return authID == MY_DISCORD_ID
}
