package icmp

import (
	"bytes"
	"encoding/binary"
)

type icmpHead struct {
	Type        uint8  // 报文类型，1字节，主要分为信息类报文和差错报文两大类
	Code        uint8  // 报文代码，1字节
	Checksum    uint16 // 报文校验和，2字节
	Identifier  uint16 // 标识符，2字节
	SequenceNum uint16 // 序列号，2字节
}

// Icmp ICMP协议包结构
type Icmp struct {
	icmpHead
	Data []byte // 数据包
}

func (o *Icmp) DefaultForPing() *Icmp {
	tmpIcmp := new(Icmp)
	tmpIcmp.Type = 8
	tmpIcmp.Code = 0
	tmpIcmp.Checksum = 0
	tmpIcmp.Identifier = 1
	tmpIcmp.SequenceNum = 1
	return tmpIcmp
}

func (o *Icmp) ToBytes() []byte {
	var buffer bytes.Buffer
	buffer.Reset()
	// binary.Write参数中的data必须为定长数据, 也就是说不能出现切片, 因此潜入了icmpHead
	_ = binary.Write(&buffer, binary.BigEndian, o.icmpHead)
	if len(o.Data) > 0 {
		buffer.Write(o.Data)
	}
	return buffer.Bytes()
}

// GenerateChecksum 生成校验和
// 参与部分 ICMP头部+数据(这里就是Icmp.Data)
// 具体操作:
// 1. 把校验和字段(Icmp.Checksum)置为0;
// 2. 对参与部分每16bit进行二进制求和;
// 3. 如果校验和的高16bit不为0，则将校验和的高16bit和低16bit反复相加，直到校验和的高16bit为0，从而获得一个16bit的值;
// 4. 将该16位的值取反，存入校验和字段;
// 验证时, 也是先计算校验和, 然后直接比较
func (o *Icmp) GenerateChecksum() {
	var sum uint32 // 16bit相加可能溢出, 因此, 这里采用32bit存储计算值(+)

	o.Checksum = 0

	dataset := o.ToBytes() // 转字节数组(注意大小端模式, 网络传输是大端)(一个字节8bit)

	for i := 0; i < len(dataset); i++ {
		if i%2 == 0 { // 高位8bit数据
			sum += uint32(dataset[i]) << 8
		} else { // 低位8bit数据
			sum += uint32(dataset[i])
		}
	}

	for sum>>16 != 0 { // 将进位到高位的16bit与低16bit再相加
		sum = sum>>16 + sum&0xffff
	}

	o.Checksum = ^uint16(sum) // 取反 强制截断 16bit
	return
}
