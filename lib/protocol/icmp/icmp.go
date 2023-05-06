package icmp

import (
	"bytes"
	"encoding/binary"
)

// Icmp ICMP协议包结构
type Icmp struct {
	Type        uint8  // 报文类型，1字节，主要分为信息类报文和差错报文两大类
	Code        uint8  // 报文代码，1字节
	Checksum    uint16 // 报文校验和，2字节
	Identifier  uint16 // 标识符，2字节
	SequenceNum uint16 // 序列号，2字节
}

func (o *Icmp) DefaultForPing() *Icmp {
	return &Icmp{
		Type:        8,
		Code:        0,
		Checksum:    0,
		Identifier:  1,
		SequenceNum: 1,
	}
}

func (o *Icmp) ToBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.BigEndian, o)
	return buffer.Bytes()
}

func (o *Icmp) GenerateChecksum(data []byte) {
	var sum uint32
	o.Checksum = 0
	// 拼接数据串
	dataset := append(o.ToBytes(), data...)
	datasetLen := len(dataset)
	idx := 0
	for datasetLen > 1 {
		sum += uint32(dataset[idx])<<8 + uint32(dataset[idx+1])
		idx += 2
		datasetLen -= 2
	}
	if datasetLen == 1 {
		sum += uint32(dataset[idx])
	}
	sum += sum >> 16
	o.Checksum = uint16(^sum)
	return
}
