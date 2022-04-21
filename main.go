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
	JuanBotSubscribers []chan string
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

	juanbotConvoChan := make(chan string, 1)
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
		fmt.Println(fmt.Sprintf("Juan Said Something in %s", m.ChannelID))
		for _, sub := range JuanBotSubscribers {
			sub <- m.ChannelID
		}
	}

	if strings.ToLower(m.Content) == "!newtype" {
		fmt.Println("Received newtype")
		for _, sub := range NewtypeSubscribers {
			sub <- 1
		}
	}
}

func runJuanBotConvo(s *discordgo.Session, juanbotConvoChan <-chan string) {
	juanbotStarter := []string{
		"Hi stepbro!",
		"Help stepbro... I'm stuck! :cold_sweat:",
		":wink:",
	}
	insultReplies := []string{
		"But... I am a bot like you... :cry:",
		"Why are you always so mean to me! :rage:",
		":sob:",
	}

	channelCount := make(map[string]int)

	for {
		channelId := <-juanbotConvoChan
		if count, ok := channelCount[channelId]; ok {
			channelCount[channelId] = count + 1
		} else {
			channelCount[channelId] = 1
		}
		if channelCount[channelId] == JuanBotConvoThreshold {
			randIndex := rand.Intn(len(juanbotStarter))
			s.ChannelMessageSend(channelId, fmt.Sprintf("%s <@%s>", juanbotStarter[randIndex], JuanBotID))
			<-juanbotConvoChan
			randIndex = rand.Intn(len(insultReplies))
			s.ChannelMessageSend(channelId, insultReplies[randIndex])
			channelCount[channelId] = 0
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
