package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

func dumpGuild(_ *cobra.Command, args []string) {
	guildId := args[0]
	guild, err := session.Guild(guildId)
	if err != nil {
		panic(err)
	}
	log.Printf("Dumping guild %s (%s)", guild.Name, guild.ID)

	chans, err := session.GuildChannels(guildId)
	if err != nil {
		panic(err)
	}
	log.Printf("Channels: %d", len(chans))

	for _, c := range chans {
		dumpChannel(guild, c)
	}
}

func getFileName(guildName string, channelName string) string {
	t := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s_%s_%s.json", guildName, channelName, t)
}

func dumpChannel(guild *discordgo.Guild, c *discordgo.Channel) {
	log.Printf("Dumping channel #%s (%s)", c.Name, c.ID)

	firstMsg, err := session.ChannelMessages(c.ID, 1, "", "", "")
	if err != nil {
		log.Printf("Failed to access #%s (%s): %s", c.Name, c.ID, err)
		return
	}
	if len(firstMsg) == 0 {
		log.Printf("Failed to access #%s (%s): Channel is empty", c.Name, c.ID)
		return
	}

	fpath := getFileName(guild.Name, c.Name)
	f, err := os.Create(fpath)
	if err != nil {
		panic(err)
	}
	wr := bufio.NewWriter(f)
	defer wr.Flush()

	jwr := json.NewEncoder(wr)

	var beforeID string
	var total int64 = 0
	for i := 1; true; i++ {
		msgs, err := session.ChannelMessages(c.ID, 100, beforeID, "", "")
		if err != nil {
			log.Printf("Failed to dump #%s (%s): %s", c.Name, c.ID, err)
			return
		}

		if len(msgs) == 0 {
			return
		}

		for _, msg := range msgs {
			mini := miniMsg {
				ID: msg.ID,
				Content: msg.Content,
				Timestamp: msg.Timestamp,
				Type: msg.Type,
				Author: miniUser {
					ID: msg.Author.ID,
					Username: msg.Author.Username,
					Bot: msg.Author.Bot,
				},
			}
			err := jwr.Encode(mini)
			if err != nil {
				panic(err)
			}
		}

		total += int64(len(msgs))

		log.Printf("Dumping #%s (%s): Got page %03d, Total: % 8d", c.Name, c.ID, i, total)

		beforeID = msgs[len(msgs)-1].ID
	}

	defer f.Close()
}
