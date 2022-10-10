package bot

import (
	"Leagly/guilds"
	"Leagly/query"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func registerCommands(s *discordgo.Session) {
	var commands []*discordgo.ApplicationCommand
	commands = append(commands, registerRegionCommand(s), registerHelpCommand(s), registerLiveCommand(s), registerLastmatchCommand(s),
		registerLookupCommand(s), registerMasteryCommand(s), registerUptimeCommand(s), registerGCCommand(s), registerWhoCommand(s),
		registerFeedbackCommand(s), registerStatusCommand(s), registerPatchnotesCommand(s))

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "930923025111580683", commands) //Dev ID 930923025111580683

	if err != nil {
		panic("Could not register commands. " + err.Error())
	}
}

// Command registration start
func registerLiveCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "live",
		Description: "Live game stats of summoner.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "summoner",
				Description: "Summoner IGN",
				Required:    true,
			},
		},
	}
	return command
}

func registerLastmatchCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "lastmatch",
		Description: "Last match game stats of summoner.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "summoner",
				Description: "Summoner IGN",
				Required:    true,
			},
		},
	}
	return command
}

func registerLookupCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "lookup",
		Description: "Last 10 ranked game stats of summoner.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "summoner",
				Description: "Summoner IGN",
				Required:    true,
			},
		},
	}
	return command
}

func registerMasteryCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "mastery",
		Description: "Mastery stats of summoner.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "summoner",
				Description: "Summoner IGN",
				Required:    true,
			},
		},
	}
	return command
}

func registerPatchnotesCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "patchnotes",
		Description: "League of Legends patchnotes",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "auto-patch-notes",
				Description: "Enable automatic patch notes updates to this channel.",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "toggle",
						Value: "toggle",
					},
				},
			},
		},
	}
	return command
}

func registerHelpCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "List of Leagly commands.",
	}
	return command
}

func registerUptimeCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "uptime",
		Description: "Time since Leagly started.",
	}
	return command
}

func registerGCCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	dmPermission := true
	var defaultMemberPermissions int64 = discordgo.PermissionManageServer
	command := &discordgo.ApplicationCommand{
		Name:                     "gc",
		Description:              "How many discord servers leagly is in.",
		DMPermission:             &dmPermission,
		DefaultMemberPermissions: &defaultMemberPermissions,
	}
	return command
}

func registerWhoCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "who",
		Description: "Debug info of guild.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "server ID",
				Required:    true,
			},
		},
	}
	return command
}

func registerFeedbackCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "feedback",
		Description: "Give feedback to the Leagly developer!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "What would you like to say?",
				Required:    true,
			},
		},
	}
	return command
}

func registerStatusCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "status",
		Description: "Status of riot servers for your current region.",
	}
	return command
}

func registerRegionCommand(s *discordgo.Session) *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "region",
		Description: "Update your league of legends region for Leagly.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "code",
				Description: "Region Code",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "North America",
						Value: "NA1",
					},
					{
						Name:  "Brazil",
						Value: "BR1",
					},
					{
						Name:  "Europe North",
						Value: "EUN1",
					},
					{
						Name:  "Europe West",
						Value: "EUW1",
					},
					{
						Name:  "Japan",
						Value: "JP1",
					},
					{
						Name:  "Korea",
						Value: "KR",
					},
					{
						Name:  "Latin America 1",
						Value: "LA1",
					},
					{
						Name:  "Latin America 2",
						Value: "LA2",
					},
					{
						Name:  "Oceania",
						Value: "OC1",
					},
					{
						Name:  "East Europe",
						Value: "RU",
					},
					{
						Name:  "Turkey",
						Value: "TR1",
					},
				},
				Required: true,
			},
		},
	}
	return command
}

// Command registration end

///
///
func live(s *discordgo.Session, interaction *discordgo.InteractionCreate, summoner string, guild guilds.DiscordGuild) {
	sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Please wait...")
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "live " + summoner)
	if onCoolDown(interaction.Member.User.ID, 3) > 0 {
		s.ChannelMessageSend(interaction.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
		log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " on cooldown")
		return
	}
	send, err := query.IsInGame(summoner, guild.Region)
	if err != nil {
		log.Println("Error: Discord server ID: " + interaction.GuildID + "  " + err.Error())
		//dont return because we still want to send an error embed
	}
	sendInteractionEdit(s, interaction.Interaction, send)
}

///
///
func lastmatch(s *discordgo.Session, interaction *discordgo.InteractionCreate, summoner string, guild guilds.DiscordGuild) {
	sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Please wait...")
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "lastmatch " + summoner)
	send, err := query.GetLastMatch(summoner, guild.Region, guild.Region2)
	if err != nil {
		log.Println("Error: Discord server ID: " + interaction.GuildID + "  " + err.Error())
	}
	sendInteractionEdit(s, interaction.Interaction, send)
}

///
///
func lookup(s *discordgo.Session, interaction *discordgo.InteractionCreate, summoner string, guild guilds.DiscordGuild) {
	sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Please wait...")
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "lookup " + summoner)
	if onCoolDown(interaction.Member.User.ID, 5) > 0 {
		s.ChannelMessageSend(interaction.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
		log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " on cooldown")
		return
	}
	send, err := query.LookupPlayer(summoner, guild.Region, guild.Region2)
	if err != nil {
		log.Println("Error: Discord server ID: " + interaction.GuildID + "  " + err.Error())
	}
	sendInteractionEdit(s, interaction.Interaction, send)
}

///
///
func mastery(s *discordgo.Session, interaction *discordgo.InteractionCreate, summoner string, guild guilds.DiscordGuild) {
	sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Please wait...")
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "mastery " + summoner)
	if onCoolDown(interaction.Member.User.ID, 3) > 0 {
		s.ChannelMessageSend(interaction.ChannelID, "You're currently on cooldown. Please wait a few seconds.")
		log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " on cooldown")
		return
	}
	send, err := query.MasteryPlayer(summoner, guild.Region)
	if err != nil {
		log.Println("Error: Discord server ID: " + interaction.GuildID + "  " + err.Error())
	}
	sendInteractionEdit(s, interaction.Interaction, send)
}

///
///
func handleHelp(s *discordgo.Session, interaction *discordgo.InteractionCreate, guild guilds.DiscordGuild) {
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "help")
	sendInteractionRespond(s, interaction, query.Help(guild.Region), "")
}

///
///
func changeRegion(s *discordgo.Session, interaction *discordgo.InteractionCreate, region string, guild guilds.DiscordGuild) {
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "region " + region)
	if guild.ID == interaction.GuildID {
		guild.Region = strings.ToUpper(region)
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
		log.Println("Discord server ID: " + interaction.GuildID + "  Changed region to " + guild.Region + " " + guild.Region2)
		sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, fmt.Sprintf("Region has been changed to %s for your discord", guild.Region))
	}
}

///
///
func uptime(s *discordgo.Session, interaction *discordgo.InteractionCreate, guild guilds.DiscordGuild) {
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "uptime")
	sendInteractionRespond(s, interaction, query.UpTime(up_time), "")
}

///
///
func status(s *discordgo.Session, interaction *discordgo.InteractionCreate, guild guilds.DiscordGuild) {
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "status")
	sendInteractionRespond(s, interaction, query.RiotApiStatus(guild.Region), "")
}

///
///
func getGuildCount(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "gc")
	sendInteractionRespond(s, interaction, query.GuildCount(guilds.GetGuildCount()), "")
}

///
///
func getGuildDebugInfo(s *discordgo.Session, interaction *discordgo.InteractionCreate, guild string) {
	log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "who")
	if guilds.HasDebugPermissions(interaction.Member.User.ID) {
		s.UpdateListeningStatus("/help") //test
		guild, err := guilds.View(guilds.DB, guild)
		if err != nil {
			log.Println(err)
		} else {
			sendInteractionRespond(s, interaction, query.GuildDebugInfo(guild), "")
		}
	} else {
		sendInteractionRespond(s, interaction, query.ErrorCreate("Sorry, you don't have permission for this command"), "")
	}
}

///
///
func feedback(s *discordgo.Session, interaction *discordgo.InteractionCreate, feedback string, guild guilds.DiscordGuild) {
	if onCoolDown(interaction.Member.User.ID, 30) > 0 {
		sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "You're currently on cooldown. Please wait a few seconds.")
		log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " on cooldown")
		return
	} else {
		_, err := s.ChannelMessageSend("955121671105286175", fmt.Sprintf("From %s, Feedback: %s ", interaction.Member.User.Username, feedback))
		if err != nil {
			sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Sorry, something went wrong :(")
			log.Println("Error sending feedback. Discord server ID: " + interaction.GuildID + "  " + err.Error())
		} else {
			sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Message has been saved! Thank you for the feedback. :)")
			log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "feedback")
		}
	}
}

///
///
func patchNotes(s *discordgo.Session, interaction *discordgo.InteractionCreate, toggle string, guild guilds.DiscordGuild) {
	if toggle == "" {
		s.ChannelTyping(interaction.ChannelID)
		log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "patchnotes")
		send, err := query.PatchNotes()
		if err != nil {
			log.Println("Error: Discord server ID: " + interaction.GuildID + "  " + err.Error())
		}
		sendInteractionRespond(s, interaction, send, "")
	} else {
		log.Println("Discord server ID: " + interaction.GuildID + "  " + interaction.Member.User.Username + " : " + "patchnotes toggle")
		guild.AutoPatchNotes = !guild.AutoPatchNotes
		guild.PatchNotesCh = interaction.ChannelID

		err := guilds.Update(guilds.DB, guild.ID, guild)
		if err != nil {
			log.Println(err)
			return
		}

		if guild.AutoPatchNotes {
			sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Auto patch notes have been enabled")
		} else {
			sendInteractionRespond(s, interaction, &discordgo.MessageSend{}, "Auto patch notes have been disabled")
		}
		log.Println("Discord server ID: " + interaction.GuildID + " auto patchnotes have been set to " + strconv.FormatBool(guild.AutoPatchNotes))
	}
}

///
///This should also be removed with discords updated message intents
func sendDiscordMessageComplex(s *discordgo.Session, m *discordgo.MessageCreate, send *discordgo.MessageSend) {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, send)
	if err != nil {
		log.Println("Error sending embed. Discord server ID: " + m.GuildID + "  " + err.Error())
	}
	//until aug 21, send one action items to discord servers asking to update permissions
	_, err = s.GuildApplicationCommandsPermissions(s.State.User.ID, m.GuildID)
	if err != nil {
		s.ChannelMessageSendComplex(m.ChannelID, query.ApplicationCommandsWarningAction(m.GuildID))
		log.Println("application.commands scope not enabled. Discord server ID: " + m.GuildID)
	}
}

///
///
func sendInteractionRespond(s *discordgo.Session, interaction *discordgo.InteractionCreate, send *discordgo.MessageSend, content string) {
	err := s.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:         content,
			Embeds:          send.Embeds,
			Files:           send.Files,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
	if err != nil {
		log.Println("Error sending embed. Discord server ID: " + interaction.GuildID + "  " + err.Error())
		_, err = s.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		if err != nil {
			log.Println("Uh Oh! Error sending interaction follow-up: Discord server ID: " + interaction.GuildID + "  " + err.Error())
		}
		return
	}
}

///
///
func sendInteractionEdit(s *discordgo.Session, interaction *discordgo.Interaction, send *discordgo.MessageSend) {
	str := "Done!"
	contentPtr := &str
	_, err := s.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{
		Content: contentPtr,
		Embeds:  &send.Embeds,
		Files:   send.Files,
	})
	if err != nil {
		log.Println("Error sending interaction embed: Discord server ID: " + interaction.GuildID + "  " + err.Error())
		_, err = s.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{ //embed error message?
			Content: "Something went wrong",
		})
		if err != nil {
			log.Println("Uh Oh! Error sending interaction follow-up: Discord server ID: " + interaction.GuildID + "  " + err.Error())
		}
		return
	}
}

///
///
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

///Some summoner names can have spaces in them
/// This method will combine each name piece into a whole string
func createName(args []string) []string {
	for n := 2; n < len(args); n++ {
		args[1] += " " + args[n]
	}
	return args
}
