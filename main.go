package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Tobotobo/msgbox"
	pop3 "github.com/knadh/go-pop3"
)

var (
	email    string
	password string
)

func main() {
	flag.StringVar(&email, "email", "", "email")
	flag.StringVar(&password, "password", "", "password")
	flag.Parse()

	// Initialize the client.
	p := pop3.New(pop3.Opt{
		Host:       "outlook.office365.com",
		Port:       995,
		TLSEnabled: true,
	})

	// Create a new connection. POP3 connections are stateful and should end
	// with a Quit() once the opreations are done.
	c, err := p.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// Authenticate.
	if err := c.Auth(email, password); err != nil {
		log.Fatal(err)
	}

	// Print the total number of messages and their size.
	count, size, _ := c.Stat()
	fmt.Println("total messages=", count, "size=", size)

	// Pull the list of all message IDs and their sizes.
	msgs, _ := c.List(0)
	for _, m := range msgs {
		fmt.Println("id=", m.ID, "size=", m.Size)
	}

	// Pull all messages on the server. Message IDs go from 1 to N.
	for id := count; id >= 1; id-- {
		m, _ := c.Retr(id)

		fmt.Println(id, "=", m.Header.Get("subject"), "@", m.Header.Get("from"))

		if strings.Contains(m.Header.Get("from"), "eie.notice@polyu.edu.hk") {
			msgbox.Show(m.Header.Get("subject"))
			if err != nil {
				panic(err)
			}
		}

		// To read the multi-part e-mail bodies, see:
		// https://github.com/emersion/go-message/blob/master/example_test.go#L12
	}

	// Delete all the messages. Server only executes deletions after a successful Quit()
	// for id := 1; id <= count; id++ {
	// 	c.Dele(id)
	// }
}
