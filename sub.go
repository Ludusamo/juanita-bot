package main

import "sync"

var (
	NewtypeSubLock     sync.Mutex
	NewtypeSubscribers map[string]chan string
	JuanBotSubLock     sync.Mutex
	JuanBotSubscribers map[string]chan string
)

func init() {
	NewtypeSubscribers = make(map[string]chan string)
	JuanBotSubscribers = make(map[string]chan string)
}

func AddSub(subType string, subName string, channel chan string) {
	if subType == "newtype" {
		NewtypeSubLock.Lock()
		NewtypeSubscribers[subName] = channel
		NewtypeSubLock.Unlock()
	} else if subType == "juanbot" {
		JuanBotSubLock.Lock()
		JuanBotSubscribers[subName] = channel
		JuanBotSubLock.Unlock()
	}
}

func RemoveSub(subType string, subName string) {
	if subType == "newtype" {
		NewtypeSubLock.Lock()
		delete(NewtypeSubscribers, subName)
		NewtypeSubLock.Unlock()
	} else if subType == "juanbot" {
		JuanBotSubLock.Lock()
		delete(JuanBotSubscribers, subName)
		JuanBotSubLock.Unlock()
	}
}
