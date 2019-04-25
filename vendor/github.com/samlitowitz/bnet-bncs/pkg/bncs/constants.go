//go:generate $GOBIN/stringer -type=MessageID
package bncs

const (
	GameProtocol  = 0x01
	BNFTPProtocol = 0x02
)

const (
	HeaderLength = 4 // bytes
)

type MessageID uint8

const (
	SidEnterChat      MessageID = 0x0a
	SidGetChannelList MessageID = 0x0b
	SidJoinChannel    MessageID = 0x0c
	SidChatEvent      MessageID = 0x0f

	SidPing MessageID = 0x25

	SidGetFiletime    MessageID = 0x33
	SidLogonResponse2 MessageID = 0x3a
	SidLogonRealmEx   MessageID = 0x3e

	SidQueryRealms2 MessageID = 0x40
	SidNewsInfo     MessageID = 0x46

	SidAuthInfo  MessageID = 0x50
	SidAuthCheck MessageID = 0x51
)
