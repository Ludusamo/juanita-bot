package main

import (
	"fmt"
	"log"
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

	newtypeChan := make(chan *discordgo.MessageCreate, 1)
	AddSub(JuanNewtypeSubType, "Newtype", newtypeChan)
	go RunEventCounter(s, newtypeChan, NewtypeThreshold, NewtypeInteraction)

	downDetectChan := make(chan *discordgo.MessageCreate, 1)
	AddSub(NewtypeSubType, "DownDetectStart", downDetectChan)
	go RunEventCounter(s, downDetectChan, 1, ShitDownDetectorInteraction)

	juanbotConvoChan := make(chan *discordgo.MessageCreate, 1)
	AddSub(JuanSubType, "JuanBotConvo", juanbotConvoChan)
	go RunEventCounter(s, juanbotConvoChan, JuanBotConvoThreshold, JuanBotConvoInteraction)

	for _, c := range ChannelIDs {
		go RunNewtypeTimeout(s, c)
	}

	for _, c := range ChannelIDs {
		go RunYoWord(s, c)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	s.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == JuanBotID {
		log.Println(fmt.Sprintf("Juan Said Something in %s", m.ChannelID))
		Notify(JuanSubType, m)

		if strings.HasPrefix(m.Content, "Bryant the type of guy to") {
			log.Println(fmt.Sprintf("Juan Said a Newtype in %s", m.ChannelID))
			Notify(JuanNewtypeSubType, m)
		}
	}

	if m.Author.ID == BryantID {
		log.Println(fmt.Sprintf("Bryant Said Something in %s", m.ChannelID))
		Notify(BryantSubType, m)
	}

	if strings.ToLower(m.Content) == "!newtype" {
		log.Println(fmt.Sprintf("Received newtype in %s", m.ChannelID))
		Notify(NewtypeSubType, m)
	}
}
