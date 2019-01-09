package stun

import (
	"NatCheck/utils"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var HeaderType = map[string]uint16{
	"BIND_REQUEST":          0x0001,
	"BIND_RESPONSE":         0x0101,
	"BIND_ERROR_RESPONSE":   0x0111,
	"SHARE_SECRET_REQUEST":  0x0002,
	"SHARE_SECRET_RESPONSE": 0x0102,
	"SHARE_SECRET_ERROR":    0x0112,
}

var AttributeType = map[string]uint16{
	"MAPPED_ADDRESS":     0x0001,
	"RESPONSE_ADDRESS":   0x0002,
	"CHANGE_REQUEST":     0x0003,
	"SOURCE_ADDRESS":     0x0004,
	"CHANGED_ADDRESS":    0x0005,
	"USERNAME":           0x0006,
	"PASSWORD":           0x0007,
	"MESSAGE_INTEGRITY":  0x0008,
	"ERROR_CODE":         0x0009,
	"UNKNOWN_ATTRIBUTES": 0x000a,
	"REFLECTED_FROM":     0x000b,
	"XOR_MAPPED_ADDRESS": 0x8020,
	"SERVER":             0x8022,
	"SECONDARY_ADDRESS":  0x8050,
}

type Message struct {
	Header     Header
	Attributes []Attribute
}

type Header struct {
	Type          []byte
	Length        []byte
	TransactionId []byte
}

type Attribute struct {
	Type   []byte
	Length []byte
	Value  []byte
}

type Address struct {
	Flag   byte
	Family byte
	Port   []byte
	Ip     []byte
}

func (m *Message) ToBytes() []byte {
	var ret []byte
	ret = append(ret, m.Header.ToBytes()...)
	for _, attr := range m.Attributes {
		ret = append(ret, attr.ToBytes()...)
	}
	return ret
}

func (m *Message) FromBytes(b []byte) {
	m.Header.Type = b[0:2]
	m.Header.Length = b[2:4]
	m.Header.TransactionId = b[4:20]

	length := utils.BytesToUint16(m.Header.Length)
	if length > 0 {
		var attrs []Attribute
		start := 20
		for {
			if start < 20+int(length) {
				var attr Attribute
				attr.Type = b[start : start+2]
				attr.Length = b[start+2 : start+4]
				attr.Value = b[start+4 : start+4+int(utils.BytesToUint16(attr.Length))]
				attrs = append(attrs, attr)
				start = start + int(attr.Len())
			} else {
				break
			}
		}
		m.Attributes = attrs
	}
}

func (m *Message) Len() uint16 {
	var ret uint16
	for _, attr := range m.Attributes {
		ret = ret + uint16(len(attr.ToBytes()))
	}
	return ret
}

func (m *Message) GetMappedAddress() *Address {
	var address Address
	for _, attr := range m.Attributes {
		if utils.BytesToUint16(attr.Type) == AttributeType["MAPPED_ADDRESS"] {
			address.FromBytes(attr.Value)
			return &address
		}
	}
	return &address
}

func (m *Message) GetChangedAddress() *Address {
	var address Address
	for _, attr := range m.Attributes {
		if utils.BytesToUint16(attr.Type) == AttributeType["CHANGED_ADDRESS"] {
			address.FromBytes(attr.Value)
		}
	}
	return &address
}

func (h *Header) ToBytes() []byte {
	var ret []byte
	ret = append(ret, h.Type...)
	ret = append(ret, h.Length...)
	ret = append(ret, h.TransactionId...)
	return ret
}

func (a *Attribute) ToBytes() []byte {
	var ret []byte
	ret = append(ret, a.Type...)
	ret = append(ret, a.Length...)
	ret = append(ret, a.Value...)
	return ret
}

func (a *Attribute) Len() uint16 {
	return uint16(len(a.ToBytes()))
}

func (a *Address) FromBytes(b []byte) {
	a.Flag = b[0]
	a.Family = b[1]
	a.Port = b[2:4]
	a.Ip = b[4:8]
}

func (a *Address) String() string {
	ip, port := "nil", "nil"
	if len(a.Ip) == 4 {
		ip = (net.IP)(a.Ip).String()
	}
	if len(a.Port) == 2 {
		port = strconv.Itoa(int(utils.BytesToUint16(a.Port)))
	}
	return ip + ":" + port
}

func RandTransactionId() utils.Uint128 {
	var ret utils.Uint128
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	ret[0] = rnd.Uint64()
	ret[1] = rnd.Uint64()
	return ret
}

func NewBindRequest() *Message {
	var message Message
	message.Header.Type = utils.Uint16ToBytes(HeaderType["BIND_REQUEST"])
	message.Header.TransactionId = utils.Uint128ToBytes(RandTransactionId())
	message.Header.Length = utils.Uint16ToBytes(message.Len())
	return &message
}

func NewChangeRequest(isChangeIp bool, isChangePort bool) *Message {
	var value uint32

	if isChangeIp {
		value = 1 << 2
	}

	if isChangePort {
		value = value + 1<<1
	}

	var attr Attribute
	attr.Type = utils.Uint16ToBytes(AttributeType["CHANGE_REQUEST"])
	attr.Value = utils.Uint32ToBytes(value)
	attr.Length = utils.Uint16ToBytes(uint16(len(attr.Value)))

	var message Message
	message.Attributes = append(message.Attributes, attr)
	message.Header.Type = utils.Uint16ToBytes(HeaderType["BIND_REQUEST"])
	message.Header.TransactionId = utils.Uint128ToBytes(RandTransactionId())
	message.Header.Length = utils.Uint16ToBytes(attr.Len())

	return &message
}
