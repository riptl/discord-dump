package main

import "github.com/bwmarrin/discordgo"

type miniUser struct {
	ID string `json:"id"`
	Username string `json:"string"`
	Bot bool `json:"bot,omitempty"`
}

type miniMsg struct {
	ID string `json:"id"`
	Content string `json:"content"`
	Timestamp discordgo.Timestamp `json:"timestamp"`
	Type discordgo.MessageType `json:"type"`
	Author miniUser `json:"author"`
}
