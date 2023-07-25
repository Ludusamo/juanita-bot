package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	juanbotConvoChan := make(chan *discordgo.MessageCreate, 1)
	AddSub(JuanSubType, "JuanBotConvo", juanbotConvoChan)
	go RunEventCounter(s, juanbotConvoChan, JuanBotConvoThreshold, JuanBotConvoInteraction)

	for _, c := range ChannelIDs {
		go RunNewtypeTimeout(s, c)
		go RunYoWord(s, c)
		go ShitDownDetectorInteraction(s, c)
	}

	go waitUntilEricBday(s)

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

func ericThisYearBday() time.Time {
	loc, _ := time.LoadLocation("America/New_York")
	return time.Date(time.Now().Year(), 5, 28, 0, 0, 0, 0, loc)
}

func waitUntilEricBday(s *discordgo.Session) {
	ericBday := ericThisYearBday()
	if time.Now().After(ericBday) {
		ericBday = ericThisYearBday().AddDate(1, 0, 0)
	}
	log.Println(fmt.Sprintf("Eric's Next Birthday is %s", ericBday))
	for {
		if time.Now().Before(ericBday) {
			timeUntilEricsBday := time.Until(ericBday)
	log.Println(fmt.Sprintf("Time until Erics Bday %s", timeUntilEricsBday))
			time.Sleep(timeUntilEricsBday / 2)
		} else {
			for _, c := range ChannelIDs {
				s.ChannelMessageSend(c,
					fmt.Sprintf("Happy birthday <@%s>! :partying_face: Thanks for creating my big bro <@%s>! :kissing_heart:", EricID, JuanBotID))
			}
			ericBday = ericThisYearBday().AddDate(1, 0, 0)
			log.Println(fmt.Sprintf("Eric's Next Birthday is %s", ericBday))
		}
	}
}

func loadSound() (buf [][]byte, err error) {
	file, err := os.Open("eric.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return nil, err
	}

	var opuslen int16

	buffer := make([][]byte, 0)
	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

func voiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.UserID == EricID && v.ChannelID != "" {
		log.Println(fmt.Sprintf("Eric joined channel %s", v.ChannelID))
		vc, err := s.ChannelVoiceJoin(v.GuildID, v.ChannelID, false, false)
		if err != nil {
			fmt.Println(err)
		}
		buffer, load_err := loadSound()
		vc.Speaking(true)
		if load_err == nil {
			for _, buff := range buffer {
				vc.OpusSend <- buff
			}
		}
		vc.Speaking(false)
		vc.Disconnect()
	}
}
