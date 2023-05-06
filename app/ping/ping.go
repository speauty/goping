package ping

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"goping/lib/protocol/icmp"
	"log"
	"net"
	"time"
)

var cmd = &cobra.Command{
	Use:     "ping addr",
	Short:   "goping，仿ping程序",
	Version: "0.0.0",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("必须指定 IP 地址。")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ipOrDomain := args[0]
		ipAddr := net.ParseIP(ipOrDomain)
		isDomain := false
		if ipAddr == nil {
			tmpIpAddr, err := net.ResolveIPAddr("ip", ipOrDomain)
			if err != nil {
				log.Fatalf("解析主机(%s)异常，错误: %s", ipOrDomain, err)
			}
			isDomain = true
			ipAddr = tmpIpAddr.IP
		}

		sendPack := new(icmp.Icmp).DefaultForPing()
		sendData := make([]byte, flagLSize)

		timeout := time.Millisecond * time.Duration(flagWTimeout)
		if isDomain {
			fmt.Printf("\n正在 Ping %s [%s] 具有 %d 字节的数据:\n", ipOrDomain, ipAddr, flagLSize)
		} else {
			fmt.Printf("\n正在 Ping %s 具有 %d 字节的数据:\n", ipAddr, flagLSize)
		}

		conn, err := net.DialTimeout("ip4:icmp", ipAddr.String(), timeout)
		if err != nil {
			log.Fatalf("链接主机(%s)异常，错误: %s", ipAddr, err)
		}
		defer func(conn net.Conn) {
			_ = conn.Close()
		}(conn)

		numSend := 0
		numReceived := 0
		minTS := 0
		maxTS := 0
		totalTS := 0

		for i := 0; i < flagNCount; i++ {
			sendPack.SequenceNum = uint16(i + 1)
			sendPack.GenerateChecksum(sendData)

			timeStart := time.Now()
			_ = conn.SetDeadline(timeStart.Add(timeout))
			_, err := conn.Write(append(sendPack.ToBytes(), sendData...))
			if err != nil {
				log.Fatalf("写入数据异常，错误: %s", err)
			}
			numSend++

			bufferResp := make([]byte, 65535)
			cntReply, err := conn.Read(bufferResp)
			if err != nil {
				fmt.Println("请求超时。")
				continue
			}
			numReceived++
			ts := int(time.Since(timeStart).Milliseconds())
			totalTS += ts
			if ts < minTS || minTS == 0 {
				minTS = ts
			}
			if ts > maxTS {
				maxTS = ts
			}
			fmt.Printf("来自 %s 的回复: 字节=%d 时间=%dms TTL=%d\n", ipAddr, len(bufferResp[28:cntReply]), ts, bufferResp[8])
			time.Sleep(time.Millisecond * 50)
		}
		fmt.Printf("\n%s 的 Ping 统计信息: \n", ipAddr)
		fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.0f%% 丢失)，\n", numSend, numReceived, numSend-numReceived, float32(numSend-numReceived)/float32(numSend))
		fmt.Println("往返行程的估计时间(以毫秒为单位):")
		fmt.Printf("    最短 = %dms，最长 = %dms，平均 = %dms\n", minTS, maxTS, totalTS/numReceived)
	},
}

var (
	flagNCount   int
	flagLSize    int
	flagWTimeout int
)

func init() {
	cmd.Flags().IntVarP(&flagNCount, "num", "n", 4, "发送请求数")
	cmd.Flags().IntVarP(&flagLSize, "large", "l", 6, "发送缓冲区大小")
	cmd.Flags().IntVarP(&flagWTimeout, "wait", "w", 1500, "等待每次回复的超时时间(ms)")
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("执行异常, 错误: %s", err)
	}
}
