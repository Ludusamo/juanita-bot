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
)

func init() {
	Token = os.Getenv("TOKEN")
	NewtypeThreshold, _ = strconv.Atoi(os.Getenv("NEWTYPE_THRESHOLD"))
	JuanBotConvoThreshold, _ = strconv.Atoi(os.Getenv("JUANBOT_CONVO_THRESHOLD"))
	JuanBotID = os.Getenv("JUANBOT_ID")
	ChannelIDs = strings.Split(os.Getenv("CHANNEL_IDS"), ",")
	NewtypeTimeout, _ = strconv.Atoi(os.Getenv("NEWTYPE_TIMEOUT"))
	JuanDeadTimeout, _ = strconv.Atoi(os.Getenv("JUAN_DEAD_TIMEOUT"))
}
