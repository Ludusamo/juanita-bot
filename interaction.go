package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

type InteractionCallback func(*discordgo.Session, string)

func RunNewtypeTimeout(s *discordgo.Session, channelId string) {
	receiveNewtype, _ := waitForChannelSpecificReply(channelId, NewtypeSubType)
	for {
		select {
		case <-receiveNewtype:
		case <-time.After(time.Duration(NewtypeTimeout) * time.Second):
			s.ChannelMessageSend(channelId, "!newtype")
		}
	}
}

func RunEventCounter(s *discordgo.Session, channelIdChan chan string, threshold int, cb InteractionCallback) {
	channelCount := make(map[string]int)
	for {
		channelId := <-channelIdChan
		if count, ok := channelCount[channelId]; ok {
			channelCount[channelId] = count + 1
		} else {
			channelCount[channelId] = 1
		}
		if channelCount[channelId] == threshold {
			go cb(s, channelId)
			channelCount[channelId] = 0
		}
	}
}

func NewtypeInteraction(s *discordgo.Session, channelId string) {
	_, err := s.ChannelMessageSend(channelId, "!newtype")
	if err != nil {
		fmt.Println("error sending DM message: ", err)
		s.ChannelMessageSend(channelId, "Failed to send you a DM. "+"Did you disable DM in your privacy settings?")
	}
}

func JuanBotConvoInteraction(s *discordgo.Session, channelId string) {
	juanbotStarter := []string{
		"Hi stepbro!",
		"Help stepbro... I'm stuck! :cold_sweat:",
		":wink:",
		"How are you doing today?",
		"Whatcha up to?",
	}
	insultReplies := []string{
		"But... I am a bot like you... :cry:",
		"Why are you always so mean to me! :rage:",
		":sob:",
		"You are such a meanie! :disappointed_relieved:",
		"I am not puny! :rage:",
	}

	randIndex := rand.Intn(len(juanbotStarter))
	s.ChannelMessageSend(channelId, fmt.Sprintf("%s <@%s>", juanbotStarter[randIndex], JuanBotID))
	receiveReply, quit := waitForChannelSpecificReply(channelId, JuanSubType)
	<-receiveReply
	quit <- 1

	randIndex = rand.Intn(len(insultReplies))
	s.ChannelMessageSend(channelId, insultReplies[randIndex])
}

func waitForChannelSpecificReply(channelId string, subType SubscriptionType) (<-chan int, chan<- int) {
	receiveReply := make(chan int)
	quit := make(chan int)
	subChan := make(chan string, 1)

	go func() {
		AddSub(subType, fmt.Sprintf("Down-%s", channelId), subChan)
		defer RemoveSub(subType, fmt.Sprintf("Down-%s", channelId))
		for {
			select {
			case chanId := <-subChan:
				if chanId == channelId {
					receiveReply <- 1
				}
			case <-quit:
				return
			}
		}
	}()
	return receiveReply, quit
}

func ShitDownDetectorInteraction(s *discordgo.Session, channelId string) {
	juanDeadQuotes := []string{
		"OH NO! JUAN IS DEAD :scream: :skull_crossbones:",
		"Where were you when Juan was kil :skull_crossbones: :sob:",
	}
	receiveReply, quit := waitForChannelSpecificReply(channelId, JuanSubType)

	select {
	case <-receiveReply:
	case <-time.After(2 * time.Second):
		randIndex := rand.Intn(len(juanDeadQuotes))
		s.ChannelMessageSend(channelId, juanDeadQuotes[randIndex])
	}
	quit <- 1
}
