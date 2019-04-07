package discordbot

import (
	"errors"
	"time"

	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/static"
)

// A DeletableMessage inhertits from message and
// provides functionalities to delete the message
// after a psecific duration.
type DeletableMessage struct {
	*discordgo.Message

	s *discordgo.Session
}

// SendEmbedMessage sends an embeded message to a channel
// and returns a deletableMessage.
//
// Parameters:
//   s      : discordgo Session
//   chanID : ID of the text channel
//   cont   : the embeds description content
//   title  : the embeds title
//   color  : the embeds color (0 -> default color)
func SendEmbedMessage(s *discordgo.Session, chanID, cont, title string, color int) (*DeletableMessage, error) {
	if color == 0 {
		color = static.ColorDefault
	}

	msg, err := s.ChannelMessageSendEmbed(chanID, &discordgo.MessageEmbed{
		Description: cont,
		Title:       title,
		Color:       color,
	})

	return &DeletableMessage{
		Message: msg,
		s:       s,
	}, err
}

// SendEmbedError sends an error embed message to a
// channel, which is defaultly colored red, and returns
// a deletableMessage.
//
// Parameters:
//   s      : discordgo Session
//   chanID : ID of the text channel
//   cont   : the embeds description content
//   title  : the embeds title
func SendEmbedError(s *discordgo.Session, chanID, cont, title string) (*DeletableMessage, error) {
	return SendEmbedMessage(s, chanID, cont, title, static.ColorRed)
}

// Delete deletes the send message.
// If the message or session is nil or if the deletion
// fails ,an error will be returned.
func (d *DeletableMessage) Delete() error {
	if d == nil || d.s == nil {
		return errors.New("message is nil")
	}

	return d.s.ChannelMessageDelete(d.ChannelID, d.ID)
}

// DeleteAfter deletes a send message after
// the specified duration and returns a
// buffered channel which will receive an
// error or nil after delete was executed.
func (d *DeletableMessage) DeleteAfter(t time.Duration) chan error {
	c := make(chan error, 1)
	time.AfterFunc(t, func() {
		c <- d.Delete()
		close(c)
	})
	return c
}
