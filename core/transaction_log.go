package main

import (
	"log"
	"iota"
	"os"
	"fmt"
	"bufio"
)

type TransactionLoggerInterface interface{
	WritePUT(key string,value string)
	WriteDELETE(key string)
	Err() <-chan error
 	ReadEvents() (<-chan Event, <-chan error)
 	Run()
}

type FileTransactionLogger struct{
	events chan<-Event
	errors <-chan error
	lastSequence uint64
	file *os.file 
}

func (l *FileTransactionLogger) WriteDELETE(key string){
	l.events <-Event{EventType: EventDelete, Key: key}
}

func (l *FileTransactionLogger) WritePUT(key string,value string){
	l.events <- Event{EventType: EventPUT, Key: key, Value: value}
}

func (l *FileTransactionLogger) Err()<-chan error{
	return l.errors
}

func (l *FileTransactionLogger) Run(){
	events := make(chan Event,16)
	l.events = events

	errors := make(chan error,1)
	l.errors = errors

	go func(){
		for e := range events{
			l.lastSequence++
			_,err := fmt.Fprintf(
				l.file,
 				"%d\t%d\t%s\t%s\n",
 				l.lastSequence, e.EventType, e.Key, e.Value
			)

			if err!=nil{
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger)ReadEvents()(<-chan Event,<-chan error){
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error,1)

	go func(){
		var e Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan(){
			line := scanner.Text()
			if err := fmt.Sscanf(line,"%d\t%d\t%s\t%s",&e.Sequence,&e.EventType,&e.Key,&e.Value); err!=nil{
				outError <- fmt.Errorf("Input Parse error: %w",err)
				return
			}
			// Sanity check! Are the sequence numbers in increasing order?
			if e.lastSequence >= e.Sequence{
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}
			l.lastSequence = e.Sequence
			outEvent <- e 

		}
		if err := scanner.Err();err!=nil{
			outError <-fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()
	return outEvent,outError
}

type Event struct{
	Sequence uint64
	EventType EventType
	Key string
	Value string
}

type EventType byte 
const (
	_ = iota
	EventDelete EventType = iota
	EventPUT 
)

func NewFileTransactionLogger(filename string)(FileTransactionLogger,error){
	file,err := os.OpenFile(filename,os.O_RDWR|os.O_APPEND|os.O_CREATE,0755)
	if err!=nil{
		return nil,fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &FileTransactionLogger{file: file},nil
}

