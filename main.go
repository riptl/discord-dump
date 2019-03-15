package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"os"
)

var token string
var isBot bool
var session *discordgo.Session

var app = cobra.Command{
	Use:     "discord-dump",
	Version: "1.1.0",
}

var cmdDump = cobra.Command{
	Use:   "dump",
	Short: "Dump Discord resources",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if token == "" {
			fmt.Fprintln(os.Stderr, "--token flag required")
			os.Exit(1)
		}

		var err error
		session, err = discordgo.New()
		session.Token = token
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

/*var cmdChannel = cobra.Command {
	Use: "channel <channel_id>",
	Args: cobra.ExactArgs(1),
	Run: dumpChannel,
}*/

var cmdGuild = cobra.Command{
	Use:   "guild <guild_id>",
	Short: "Dump a Guild",
	Args:  cobra.ExactArgs(1),
	Run:   dumpGuild,
}

func main() {
	//app.AddCommand(&cmdChannel)
	app.AddCommand(&cmdDump)
	cmdDump.AddCommand(&cmdGuild)

	app.AddCommand(&cmdAnalyze)
	cmdAnalyze.AddCommand(&cmdBotRatio)
	cmdAnalyze.AddCommand(&cmdUserFrequency)

	{
		pf := cmdDump.PersistentFlags()
		pf.StringVar(&token, "token", "", "User/Bot Token")
		pf.BoolVar(&isBot, "bot", false, "Token is a bot token (default user)")
	}

	err := app.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
