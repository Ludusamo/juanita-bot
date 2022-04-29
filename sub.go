package main

import (
	"sync"
)

type SubscriptionMap map[string]chan string
type SubscriptionType int

const (
	JuanSubType SubscriptionType = iota
	NewtypeSubType
	JuanNewtypeSubType
)

var (
	SubscriptionLocks map[SubscriptionType]*sync.Mutex
	Subscriptions     map[SubscriptionType]SubscriptionMap
)

func init() {
	SubscriptionLocks = make(map[SubscriptionType]*sync.Mutex)
	Subscriptions = make(map[SubscriptionType]SubscriptionMap)
}

func AddSub(subType SubscriptionType, subName string, channel chan string) {
	lock, exists := SubscriptionLocks[subType]
	if !exists {
		lock = &sync.Mutex{}
		SubscriptionLocks[subType] = lock
		Subscriptions[subType] = make(SubscriptionMap)
	}
	lock.Lock()
	defer lock.Unlock()
	Subscriptions[subType][subName] = channel
}

func RemoveSub(subType SubscriptionType, subName string) {
	lock, exists := SubscriptionLocks[subType]
	if !exists {
		lock = &sync.Mutex{}
		SubscriptionLocks[subType] = lock
		Subscriptions[subType] = make(SubscriptionMap)
	}
	lock.Lock()
	defer lock.Unlock()
	delete(Subscriptions[subType], subName)
}
