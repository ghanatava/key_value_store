package main 

import (
	"testing"
	"errors"
)

func TestPut(t *testing.T){
	const key = "create-key"
	const value = "create-value"

	var val interface{}
	var contains bool

	defer delete(store,key)
	 //sanity check

	 _,contains = store[key]
	 if contains {
		t.Error("key already exists")
	 }

	 //error should be nil
	 err := Put(key,value)
	 if err!=nil{
		t.Error(err)
	 }

	 val,contains = store[key]
	 if !contains{
		t.Error("creation failed")
	 }

	 if val != value{
		t.Error("val/value mismatch")
	 }
}

func TestGet(t *testing.T){
	const key = "read-key"
	const value = "read-value"
	var val interface{}
	var err error
	defer delete(store,key)
	//testing to get a non-existent key
	val,err = Get(key)
	if err == nil{
		t.Error("expected ErrorNoSuchKey")
	}
	if !errors.Is(err,ErrorNoSuchKey){
		t.Error("unexpected error")
	}

	//testing to get an existing key

	store[key] = value

	val,err = Get(key)

	if err != nil{
		t.Error("unexpected error")
	}

	if val != value{
		t.Error("var/value mismatch")
	}
}

func TestDelete(t *testing.T){
	const key = "delete-key"
	const value = "delete-val"

	var contains bool

	defer delete(store,key)

	//testing deletion
	store[key] = value

	_,contains = store[key]
	if !contains{
		t.Error("Key creation failed")
	}
	
	Delete(key)
	_,contains = store[key]
	if contains{
		t.Error("could not delete key")
	}
	
}