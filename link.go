package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type LinkParser struct {
	rPublicLink       *regexp.Regexp
	rPublicThreadLink *regexp.Regexp

	rPrivateLink       *regexp.Regexp
	rPrivateThreadLink *regexp.Regexp
}

type Link struct {
	Username  string
	ChatId    int64
	MessageId int64
}

func NewLinkParser() (*LinkParser, error) {
	rPublicLink, err := regexp.Compile(`t\.me\/(\w{4,})\/(\d+)`)
	if err != nil {
		return nil, err
	}
	rPublicThreadLink, err := regexp.Compile(`t\.me\/(\w{4,})\/(\d+)\/(\d+)`)
	if err != nil {
		return nil, err
	}
	rPrivateLink, err := regexp.Compile(`t\.me\/c\/(\d+)\/(\d+)`)
	if err != nil {
		return nil, err
	}
	rPrivateThreadLink, err := regexp.Compile(`t\.me\/c\/(\d+)\/(\d+)\/(\d+)`)
	if err != nil {
		return nil, err
	}
	return &LinkParser{
		rPublicLink:        rPublicLink,
		rPublicThreadLink:  rPublicThreadLink,
		rPrivateLink:       rPrivateLink,
		rPrivateThreadLink: rPrivateThreadLink,
	}, nil
}

func (p *LinkParser) ParseLink(link string) (*Link, error) {
	match := p.rPublicThreadLink.FindStringSubmatch(link)
	if len(match) >= 4 {
		msgId, err := strconv.ParseInt(match[3], 10, 64)
		if err == nil {
			return &Link{
				Username:  match[1],
				MessageId: msgId,
			}, nil
		}
	}

	match = p.rPublicLink.FindStringSubmatch(link)
	if len(match) >= 3 {
		msgId, err := strconv.ParseInt(match[2], 10, 64)
		if err == nil {
			return &Link{
				Username:  match[1],
				MessageId: msgId,
			}, nil
		}
	}

	match = p.rPrivateThreadLink.FindStringSubmatch(link)
	if len(match) >= 4 {
		chatId, err1 := strconv.ParseInt(strings.Join([]string{"-100", match[1]}, ""), 10, 64)
		msgId, err2 := strconv.ParseInt(match[3], 10, 64)
		if err1 == nil && err2 == nil {
			return &Link{
				ChatId:  chatId,
				MessageId: msgId,
			}, nil
		}
	}
	
	match = p.rPrivateLink.FindStringSubmatch(link)
	if len(match) >= 3 {
		chatId, err1 := strconv.ParseInt(strings.Join([]string{"-100", match[1]}, ""), 10, 64)
		msgId, err2 := strconv.ParseInt(match[2], 10, 64)
		if err1 == nil && err2 == nil {
			return &Link{
				ChatId:  chatId,
				MessageId: msgId,
			}, nil
		}
	}
	
	return nil, errors.New("not a link")
}
