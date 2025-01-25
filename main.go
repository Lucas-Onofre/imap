package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	fmt.Println("Starting proccess...")

	if email == "" || password == "" {
		log.Fatal("Error to get email or password")
	}

	c, err := connectToIMAP("imap.gmail.com:993")
	if err != nil {
		log.Fatalf("Error connecting to IMAP server: %v", err)
	}
	defer logoutIMAP(c)

	// server login
	if err := c.Login(email, password); err != nil {
		log.Fatalf("Error in login: %v", err)
	}
	fmt.Println("Logged succesfully.")

	// select mail box (inbox by default)
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatalf("Error selecting inbox: %v", err)
	}
	fmt.Printf("Inbox selected, %d messages.\n", mbox.Messages)


	// fetch last 10 messages
	if err := fetchLastMessages(c, mbox.Messages, 10); err != nil {
		log.Fatalf("Error fetching messages: %v", err)
	}
}

// connect to imap server
func connectToIMAP(server string) (*client.Client, error) {
	c, err := client.DialTLS(server, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to server IMAP.")
	return c, nil
}

// logout imap server
func logoutIMAP(c *client.Client) {
	if err := c.Logout(); err != nil {
		log.Printf("Error on logout: %v", err)
	} else {
		fmt.Println("Succesfully logout.")
	}
}

// fetch and print messages
func fetchLastMessages(c *client.Client, totalMessages uint32, count int) error {
	if totalMessages == 0 {
		fmt.Println("No messages found.")
		return nil
	}

	// set range to fetch
	from := uint32(1)
	to := totalMessages
	if totalMessages > uint32(count) {
		from = totalMessages - uint32(count) + 1
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	// open channel to receive fetched messages 
	messages := make(chan *imap.Message, count)
	go func() {
		if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
			log.Fatalf("Error fetchig messages: %v", err)
		}
	}()

	fmt.Println("Last messages:")
	for msg := range messages {
		if msg.Envelope != nil {
			fmt.Printf("* %v\n", msg.Envelope.Subject)
		} else {
			fmt.Println("* (No subject messages)")
		}
	}
	return nil
}
