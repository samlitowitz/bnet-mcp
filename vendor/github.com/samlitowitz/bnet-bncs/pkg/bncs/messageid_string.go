// Code generated by "stringer -type=MessageID"; DO NOT EDIT.

package bncs

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SidEnterChat-10]
	_ = x[SidGetChannelList-11]
	_ = x[SidJoinChannel-12]
	_ = x[SidChatEvent-15]
	_ = x[SidPing-37]
	_ = x[SidGetFiletime-51]
	_ = x[SidLogonResponse2-58]
	_ = x[SidLogonRealmEx-62]
	_ = x[SidQueryRealms2-64]
	_ = x[SidNewsInfo-70]
	_ = x[SidAuthInfo-80]
	_ = x[SidAuthCheck-81]
}

const (
	_MessageID_name_0 = "SidEnterChatSidGetChannelListSidJoinChannel"
	_MessageID_name_1 = "SidChatEvent"
	_MessageID_name_2 = "SidPing"
	_MessageID_name_3 = "SidGetFiletime"
	_MessageID_name_4 = "SidLogonResponse2"
	_MessageID_name_5 = "SidLogonRealmEx"
	_MessageID_name_6 = "SidQueryRealms2"
	_MessageID_name_7 = "SidNewsInfo"
	_MessageID_name_8 = "SidAuthInfoSidAuthCheck"
)

var (
	_MessageID_index_0 = [...]uint8{0, 12, 29, 43}
	_MessageID_index_8 = [...]uint8{0, 11, 23}
)

func (i MessageID) String() string {
	switch {
	case 10 <= i && i <= 12:
		i -= 10
		return _MessageID_name_0[_MessageID_index_0[i]:_MessageID_index_0[i+1]]
	case i == 15:
		return _MessageID_name_1
	case i == 37:
		return _MessageID_name_2
	case i == 51:
		return _MessageID_name_3
	case i == 58:
		return _MessageID_name_4
	case i == 62:
		return _MessageID_name_5
	case i == 64:
		return _MessageID_name_6
	case i == 70:
		return _MessageID_name_7
	case 80 <= i && i <= 81:
		i -= 80
		return _MessageID_name_8[_MessageID_index_8[i]:_MessageID_index_8[i+1]]
	default:
		return "MessageID(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}