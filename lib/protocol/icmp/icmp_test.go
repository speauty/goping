package icmp

import (
	"fmt"
	"testing"
)

func BenchmarkIcmp_ToBytes(b *testing.B) {
	tmpIcmp := new(Icmp).DefaultForPing()
	for i := 0; i < b.N; i++ {
		tmpIcmp.ToBytes()
	}
}

func BenchmarkIcmp_DefaultForPing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		new(Icmp).DefaultForPing()
	}
}

func BenchmarkIcmp_GenerateChecksum(b *testing.B) {
	tmpIcmp := new(Icmp).DefaultForPing()
	for i := 0; i < b.N; i++ {
		tmpIcmp.SequenceNum = uint16(i + 1)
		tmpIcmp.GenerateChecksum()
	}
}

func TestIcmp_ToBytes(t *testing.T) {
	tmpIcmp := new(Icmp).DefaultForPing()
	tmpBytesStr := fmt.Sprintf("%x", tmpIcmp.ToBytes())
	targetStr := "0800000000010001"
	if tmpBytesStr != targetStr {
		t.Error("Expected", targetStr, "Got", tmpBytesStr)
	}
}

func TestIcmp_GenerateChecksum(t *testing.T) {
	sum := uint16(63485)
	tmpIcmp := new(Icmp).DefaultForPing()
	tmpIcmp.GenerateChecksum()
	if tmpIcmp.Checksum != sum {
		t.Error("Expected", sum, "Got", tmpIcmp.Checksum)
	}
}
