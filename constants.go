package main

import (
	"os"
	"strconv"
	"strings"
)

var (
	Token                 string
	NewtypeThreshold      int
	JuanBotConvoThreshold int
	JuanBotID             string
	ChannelIDs            []string
	NewtypeTimeout        int
	JuanDeadTimeout       int
	BryantID              string
	YoWordTimeout         int
	EricID                string
)

func init() {
	Token = os.Getenv("TOKEN")
	NewtypeThreshold, _ = strconv.Atoi(os.Getenv("NEWTYPE_THRESHOLD"))
	JuanBotConvoThreshold, _ = strconv.Atoi(os.Getenv("JUANBOT_CONVO_THRESHOLD"))
	JuanBotID = os.Getenv("JUANBOT_ID")
	ChannelIDs = strings.Split(os.Getenv("CHANNEL_IDS"), ",")
	NewtypeTimeout, _ = strconv.Atoi(os.Getenv("NEWTYPE_TIMEOUT"))
	JuanDeadTimeout, _ = strconv.Atoi(os.Getenv("JUAN_DEAD_TIMEOUT"))
	BryantID = os.Getenv("BRYANT_ID")
	YoWordTimeout, _ = strconv.Atoi(os.Getenv("YO_WORD_TIMEOUT"))
	EricID = os.Getenv("ERIC_ID")
}
