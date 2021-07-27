package wire

import "time"

// algorithm in routing
const (
	AlgorithmHashSlots = "hashslots"
)

// Command defined data type between client and server
const (
	// login
	CommandLoginSignIn  = "login.signin"
	CommandLoginSignOut = "login.signout"

	// chat
	CommandChatUserTalk  = "chat.user.talk"
	CommandChatGroupTalk = "chat.group.talk"
	CommandChatTalkAck   = "chat.talk.ack"

	// 离线
	CommandOfflineIndex   = "chat.offline.index"
	CommandOfflineContent = "chat.offline.content"

	// 群管理
	CommandGroupCreate  = "chat.group.create"
	CommandGroupJoin    = "chat.group.join"
	CommandGroupQuit    = "chat.group.quit"
	CommandGroupMembers = "chat.group.members"
	CommandGroupDetail  = "chat.group.detail"
)

// Meta Key of a packet
const (
	// 消息将要送达的网关的ServiceName
	MetaDestServer = "dest.server"
	// 消息将要送达的channels
	MetaDestChannels = "dest.channels"
)

// Protocol Protocol
type Protocol string

// Protocol
const (
	ProtocolTCP       Protocol = "tcp"
	ProtocolWebsocket Protocol = "websocket"
)

// Service Name 定义统一的服务名
const (
	SNWGateway = "wgateway"
	SNTGateway = "tgateway"
	SNLogin    = "chat"    //login
	SNChat     = "chat"    //chat
	SNService  = "service" //rpc service
)

// ServiceID ServiceID
type ServiceID string

// SessionID SessionID
type SessionID string

type Magic [4]byte

var (
	MagicLogicPkt = Magic{0xc3, 0x11, 0xa3, 0x65}
	MagicBasicPkt = Magic{0xc3, 0x15, 0xa7, 0x65}
)

const (
	OfflineMessageExpiresIn = time.Hour * 24 * 30
	OfflineSyncIndexCount   = 3000
	OfflineMessageStoreDays = 30 //days
)

const (
	MessageTypeText  = 1 // 文本消息
	MessageTypeImage = 2
	MessageTypeVoice = 3
	MessageTypeVideo = 4
)

const ()
