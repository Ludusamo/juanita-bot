package main

import (
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	ChannelId string
	NewtypeThreshold int
	NewtypeSubscribers []chan int
)

func init() {
	Token = os.Getenv("TOKEN")
	ChannelId = os.Getenv("CHANNEL_ID")
	NewtypeThreshold, _ = strconv.Atoi(os.Getenv("NEWTYPE_THRESHOLD"))
}

func main() {
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session", err)
		return
	}
	s.AddHandler(messageCreate)

	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	newtypeChan := make(chan int, 1)
	NewtypeSubscribers = append(NewtypeSubscribers, newtypeChan)
	go runNewtype(s, newtypeChan)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	s.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.ToLower(m.Content) == "!newtype" {
		fmt.Println("Received newtype")
		for _, sub := range NewtypeSubscribers {
			sub <- 1
		}
	}
}

func runNewtype(s *discordgo.Session, newtypeChan <-chan int) {
	for {
		for i:= 0; i < NewtypeThreshold; i++ {
			<-newtypeChan
		}
		_, err := s.ChannelMessageSend(ChannelId, "!newtype")
		if err != nil {
			fmt.Println("error sending DM message: ", err)
			s.ChannelMessageSend(ChannelId, "Failed to send you a DM. " + "Did you disable DM in your privacy settings?")
		}
	}
}
