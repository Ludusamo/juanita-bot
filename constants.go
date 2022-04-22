package main

import (
	"os"
	"strconv"
)

var (
	Token                 string
	NewtypeThreshold      int
	JuanBotConvoThreshold int
	JuanBotID             string
)

func init() {
	Token = os.Getenv("TOKEN")
	NewtypeThreshold, _ = strconv.Atoi(os.Getenv("NEWTYPE_THRESHOLD"))
	JuanBotConvoThreshold, _ = strconv.Atoi(os.Getenv("JUANBOT_CONVO_THRESHOLD"))
	JuanBotID = os.Getenv("JUANBOT_ID")
}
