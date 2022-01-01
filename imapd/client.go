package imapd

import (
	"LookIn/utils"
	"crypto/tls"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/patrickmn/go-cache"
	"io"
	"io/ioutil"
	"log"
	"strconv"
)

func LoginRetrieve(ca *cache.Cache, info utils.EmailServer, username string, password string) *cache.Cache {
	var c *client.Client

	var err error

	if info.Proto == "startls" {
		c, err = client.Dial(utils.JoinAddr(info.Addr, info.Port))
		if err != nil {
			log.Fatal(err)
		}
		defer c.Logout()
		if err = c.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			log.Fatal(err)
		}
		if err := c.Login(username, password); err != nil {
			log.Fatal(err)
		}
	} else if info.Proto == "tls" {
		c, err = client.DialTLS(utils.JoinAddr(info.Addr, info.Port), &tls.Config{InsecureSkipVerify: true})
		if err := c.Login(username, password); err != nil {
			log.Fatal(err)
		}
	} else if info.Proto == "plain" {
		c, err = client.Dial(utils.JoinAddr(info.Addr, info.Port))
		if err := c.Login(username, password); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("SHOULD NOT GET HERE")
	}
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(1, mbox.Messages)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal("Fetch " + err.Error())
		}
	}()
	n := 0
	for msg := range messages {
		var entry utils.CacheEntry
		if msg == nil {
			continue
		}
		r := msg.GetBody(&section)
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal("Reader: " + err.Error())
		}
		entry.Header = mr.Header
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal("Next Part: " + err.Error())
			}
			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				b, _ := ioutil.ReadAll(p.Body)

				entry.Body = append(entry.Body, b)
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				entry.Attachment.Exists = true
				entry.Attachment.Filenames = append(entry.Attachment.Filenames, filename)
			}
		}
		ca.Add(strconv.Itoa(n), entry, cache.NoExpiration)
		n += 1
	}
	return ca
}
