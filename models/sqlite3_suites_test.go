// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

func TestUpsert(t *testing.T) {
	t.Run("Chats", testChatsUpsert)

	t.Run("ChatPeers", testChatPeersUpsert)

	t.Run("Messages", testMessagesUpsert)

	t.Run("Peers", testPeersUpsert)
}