package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

var cmdAnalyze = cobra.Command{
	Use:   "analyze",
	Short: "Analyze dumps",
}

var cmdUserFrequency = cobra.Command{
	Use:   "user-frequency <file>",
	Short: "Messages per user",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		runAnalysis(&analysisUserFrequency{
			userById: make(map[string]miniUser),
			msgCount: make(map[string]int64),
		}, args[0])
	},
}

var cmdBotRatio = cobra.Command{
	Use:   "bot-ratio <file>",
	Short: "Human/bot ratio",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		runAnalysis(new(analyisBotRatio), args[0])
	},
}

func runAnalysis(a analysis, fpath string) {
	f, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	for i := 1; scan.Scan(); i++ {
		line := scan.Bytes()
		var msg miniMsg
		err = json.Unmarshal(line, &msg)
		if err != nil {
			panic(err)
		}
		a.Add(&msg)

		if i%1000 == 0 {
			short := a.Short()
			if short != "" {
				fmt.Printf("Scanning (%5d): %s\n", i, a.Short())
			} else {
				fmt.Printf("Scanning (%5d)\n", i)
			}
		}
	}
	fmt.Println(a.String())
}

type analysis interface {
	Add(msg *miniMsg)
	Short() string
	String() string
}

type analyisBotRatio struct {
	BotMsgs   int64
	HumanMsgs int64
}

func (a *analyisBotRatio) Add(msg *miniMsg) {
	if msg.Author.Bot {
		a.BotMsgs++
	} else {
		a.HumanMsgs++
	}
}

func (a *analyisBotRatio) Short() string {
	return a.String()
}

func (a *analyisBotRatio) String() string {
	return fmt.Sprintf("%.4f (%6d By Bots : %6d By Humans)",
		float64(a.BotMsgs)/(float64(a.BotMsgs)+float64(a.HumanMsgs)),
		a.BotMsgs, a.HumanMsgs)
}

type analysisUserFrequency struct {
	userById map[string]miniUser
	msgCount map[string]int64
}

func (a *analysisUserFrequency) Add(msg *miniMsg) {
	if _, ok := a.userById[msg.Author.ID]; !ok {
		a.userById[msg.Author.ID] = msg.Author
	}
	a.msgCount[msg.Author.ID]++
}

func (a *analysisUserFrequency) Short() string {
	return ""
}

func (a *analysisUserFrequency) String() string {
	type row struct {
		User  miniUser `json:"user"`
		Count int64    `json:"count"`
	}

	var rows []row
	for id, count := range a.msgCount {
		rows = append(rows, row{a.userById[id], count})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Count > rows[j].Count
	})

	var buf bytes.Buffer
	for _, r := range rows {
		res, _ := json.Marshal(r)
		buf.WriteByte('\n')
		buf.Write(res)
	}
	return buf.String()
}
