package main

import (
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"strconv"
	"strings"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	ChannelId string
	NewtypeThreshold int
	JuanBotConvoThreshold int
	JuanBotID string

	NewtypeSubscribers []chan int
	JuanBotSubscribers []chan int
)

func init() {
	Token = os.Getenv("TOKEN")
	ChannelId = os.Getenv("CHANNEL_ID")
	NewtypeThreshold, _ = strconv.Atoi(os.Getenv("NEWTYPE_THRESHOLD"))
	JuanBotConvoThreshold, _ = strconv.Atoi(os.Getenv("JUANBOT_CONVO_THRESHOLD"))
	JuanBotID = os.Getenv("JUANBOT_ID")
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

	juanbotConvoChan := make(chan int, 1)
	JuanBotSubscribers = append(JuanBotSubscribers, juanbotConvoChan)
	go runJuanBotConvo(s, juanbotConvoChan)

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
		fmt.Println("Juan Said Something")
		for _, sub := range JuanBotSubscribers {
			sub <- 1
		}
	}

	if strings.ToLower(m.Content) == "!newtype" {
		fmt.Println("Received newtype")
		for _, sub := range NewtypeSubscribers {
			sub <- 1
		}
	}
}

func runJuanBotConvo(s *discordgo.Session, juanbotConvoChan <-chan int) {
	juanbotStarter := []string{
		"hi stepbro!",
		"help stepbro... I'm stuck! :cold_sweat:",
		":wink:",
	}
	insultReplies := []string{
		"But... I am a bot like you... :cry:",
		"Why are you always so mean to me! :rage:",
		":sob:",
	}

	for {
		for i:= 0; i < JuanBotConvoThreshold; i++ {
			<-juanbotConvoChan
		}
		randIndex := rand.Intn(len(juanbotStarter))
		s.ChannelMessageSend(ChannelId, fmt.Sprintf("<@%s>, %s", JuanBotID, juanbotStarter[randIndex]))
		<-juanbotConvoChan
		randIndex = rand.Intn(len(insultReplies))
		s.ChannelMessageSend(ChannelId, insultReplies[randIndex])
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
