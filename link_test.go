package main

import (
	"testing"
)

func TestParsePublicLink(t *testing.T) {
	link := "https://t.me/username/114514"
	want := Link{
		Username: "username",
		MessageId: 114514,
	}
	parser, err := NewLinkParser()
	if err != nil {
		t.Errorf("failed to initialize parser: %v+", err)
	}
	parsed, err := parser.ParseLink(link)
	if err != nil {
		t.Errorf("failed to parse link: %v+", err)
	}
	if *parsed != want {
		t.Errorf("failed to parse public link")
	}
}

func TestParsePublicThreadLink(t *testing.T) {
	link := "https://t.me/username/114514/1919810"
	want := Link{
		Username: "username",
		MessageId: 1919810,
	}
	parser, err := NewLinkParser()
	if err != nil {
		t.Errorf("failed to initialize parser: %v+", err)
	}
	parsed, err := parser.ParseLink(link)
	if err != nil {
		t.Errorf("failed to parse link: %v+", err)
	}
	if *parsed != want {
		t.Log(parsed)
		t.Errorf("failed to parse public thread link")
	}
}

func TestParsePrivateLink(t *testing.T) {
	link := "https://t.me/c/114514/1919810"
	want := Link{
		ChatId: -100114514,
		MessageId: 1919810,
	}
	parser, err := NewLinkParser()
	if err != nil {
		t.Errorf("failed to initialize parser: %v+", err)
	}
	parsed, err := parser.ParseLink(link)
	if err != nil {
		t.Errorf("failed to parse link: %v+", err)
	}
	if *parsed != want {
		t.Errorf("failed to parse public link")
	}
}

func TestParsePrivateThreadLink(t *testing.T) {
	link := "https://t.me/c/114514/1919/810"
	want := Link{
		ChatId: -100114514,
		MessageId: 810,
	}
	parser, err := NewLinkParser()
	if err != nil {
		t.Errorf("failed to initialize parser: %v+", err)
	}
	parsed, err := parser.ParseLink(link)
	if err != nil {
		t.Errorf("failed to parse link: %v+", err)
	}
	if *parsed != want {
		t.Errorf("failed to parse public thread link")
	}
}
