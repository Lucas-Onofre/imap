package main

import (
	"fmt"
	"log"
	"os"

	"Lucas-Onofre/go-imap-poc/src/client"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	fmt.Println("Starting proccess...")

	if email == "" || password == "" {
		log.Fatal("Error to get email or password")
	}

	imapClient, err := client.NewClient("imap.gmail.com:993", email, password) 
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer imapClient.Logout()

	mbox, err := imapClient.SelectMailBox("INBOX") 
	if err != nil {
		log.Fatalf("Error selecting inbox: %v", err)
	}
	fmt.Printf("Inbox selected, %d messages.\n", mbox.Messages)


	// fetch last 10 messages
	if err := imapClient.FetchLastMessages(mbox.Messages, 10); err != nil {
		log.Fatalf("Error fetching messages: %v", err)
	}
}