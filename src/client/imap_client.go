package client

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type imapClient struct {
	c *client.Client
}

func NewClient(server string, email string, password string) (*imapClient, error) {
	c, err := client.DialTLS(server, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Login(email, password); err != nil {
		return nil, err
	}

	return &imapClient{c}, nil
}

func (i *imapClient) Logout() error {
	return i.c.Logout()
}

func (i *imapClient) SelectMailBox(mailbox string) (*imap.MailboxStatus, error) {
	mbox, err := i.c.Select(mailbox, false)
	return mbox, err
}

func (i *imapClient) FetchLastMessages(messages uint32, limit int) error {
	return i.fetchLastMessages(messages, limit)
}

func (i *imapClient) fetchLastMessages(totalMessages uint32, count int) error {
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
		if err := i.c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
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
