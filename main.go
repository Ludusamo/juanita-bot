package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

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

	newtypeChan := make(chan string, 1)
	AddSub("newtype", "NewtypeInteraction", newtypeChan)
	go RunEventCounter(s, newtypeChan, NewtypeThreshold, NewtypeInteraction)

	downDetectChan := make(chan string, 1)
	AddSub("newtype", "DownDetect", downDetectChan)
	go RunEventCounter(s, downDetectChan, 1, ShitDownDetectorInteraction)

	juanbotConvoChan := make(chan string, 1)
	AddSub("juanbot", "Convo", juanbotConvoChan)
	go RunEventCounter(s, juanbotConvoChan, JuanBotConvoThreshold, JuanBotConvoInteraction)

	for _, c := range ChannelIDs {
		go RunNewtypeTimeout(s, c)
	}

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

	if m.Author.ID == JuanBotID {
		fmt.Println(fmt.Sprintf("Juan Said Something in %s", m.ChannelID))
		JuanBotSubLock.Lock()
		for _, sub := range JuanBotSubscribers {
			sub <- m.ChannelID
		}
		JuanBotSubLock.Unlock()
	}

	if strings.ToLower(m.Content) == "!newtype" {
		fmt.Println(fmt.Sprintf("Received newtype in %s", m.ChannelID))
		NewtypeSubLock.Lock()
		for _, sub := range NewtypeSubscribers {
			sub <- m.ChannelID
		}
		NewtypeSubLock.Unlock()
	}
}
