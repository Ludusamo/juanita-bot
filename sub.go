package main

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type SubscriptionMap map[string]chan *discordgo.MessageCreate
type SubscriptionType int

const (
	JuanSubType SubscriptionType = iota
	NewtypeSubType
	JuanNewtypeSubType
	BryantSubType
)

var (
	SubscriptionLocks map[SubscriptionType]*sync.Mutex
	Subscriptions     map[SubscriptionType]SubscriptionMap
)

func init() {
	SubscriptionLocks = make(map[SubscriptionType]*sync.Mutex)
	Subscriptions = make(map[SubscriptionType]SubscriptionMap)
}

func getLock(subType SubscriptionType) *sync.Mutex {
	lock, exists := SubscriptionLocks[subType]
	if !exists {
		lock = &sync.Mutex{}
		SubscriptionLocks[subType] = lock
		Subscriptions[subType] = make(SubscriptionMap)
	}
	return lock
}

func AddSub(subType SubscriptionType, subName string, channel chan *discordgo.MessageCreate) {
	lock := getLock(subType)
	lock.Lock()
	defer lock.Unlock()
	Subscriptions[subType][subName] = channel
}

func RemoveSub(subType SubscriptionType, subName string) {
	lock := getLock(subType)
	lock.Lock()
	defer lock.Unlock()
	delete(Subscriptions[subType], subName)
}

func Notify(subType SubscriptionType, msg *discordgo.MessageCreate) {
	lock := getLock(subType)
	lock.Lock()
	defer lock.Unlock()
	for _, sub := range Subscriptions[subType] {
		sub <- msg
	}
}
