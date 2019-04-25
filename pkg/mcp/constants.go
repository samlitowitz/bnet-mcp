//go:generate $GOBIN/stringer -type=MessageID
package mcp

type MessageID uint8

const (
	McpStartup MessageID = iota + 1
	McpCharCreate
	McpCreateGame
	McpJoinGame
	McpGameList
	McpGameInfo
	McpCharLogon
	// 0x08
	// 0x09
	McpCharDelete
	// 0x0b
	// 0x0c
	// 0x0d
	// 0x0e
	// 0x0f
	// 0x10
	McpRequestLadderData
	McpMOTD
	McpCancelGameCreate
	McpCreateQueue
	// 0x15
	McpCharRank
	McpCharList
	McpCharUpgrade
	McpCharList2
)
