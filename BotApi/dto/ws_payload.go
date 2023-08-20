package dto

type EventType string

// websocket 消息载体结构
type OPCode int
type WSPayload struct {
	OPCode     OPCode      `json:"op"`
	Seq        uint32      `json:"s,omitempty"`
	Type       EventType   `json:"t,omitempty"`
	Data       interface{} `json:"d,omitempty"`
	RawMessage []byte      `json:"-"`
}

// WSHelloData
type WSHelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// WSReadyData
type WSReadyData struct {
	Version   int    `json:"version"`
	SessionID string `json:"session_id"`
	User      struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Bot      bool   `json:"bot"`
	} `json:"user"`
	Shard []uint32 `json:"shard"`
}

// WSATMessageData
type WSATMessageData Message

// WSIdentityData
type WSIdentityData struct {
	Token      string   `json:"token"`
	Intents    Intent   `json:"intents"`
	Shard      []uint32 `json:"shard"` // (shard_id, num_shards)
	Properties struct {
		Os      string `json:"$os,omitempty"`
		Browser string `json:"$browser,omitempty"`
		Device  string `json:"$device,omitempty"`
	} `json:"properties,omitempty"`
}

// OPCode
const (
	OPCODE_DISPATCH OPCode = iota
	OPCODE_HEARTBEAT
	OPCODE_IDENTIFY
	_
	_
	_
	OPCODE_RESUME
	OPCODE_RECONNECT
	_
	OPCODE_INVALID_SESSION
	OPCODE_HELLO
	OPCODE_HEARTBEAT_ACK
	OPCODE_HTTP_CALLBACK_ACK
)
