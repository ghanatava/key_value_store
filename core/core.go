package main

import (
	"errors"
	"sync"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m : make(map[string]string)}


func Put (key string, value string) error{
	store.Lock()
	defer store.Unlock()
	store.m[key] = value
	return nil
}

var ErrorNoSuchKey = errors.New("no such key")

func Get(key string)(string,error){
	store.RLock()

	value,ok := store.m[key]

	store.RUnlock()

	if !ok {
		return "",ErrorNoSuchKey
	}
	return value,nil
}

func Delete(key string)error{
	store.Lock()
	defer store.Unlock()
	delete(store.m,key)
	return nil
}
