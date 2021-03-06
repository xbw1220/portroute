package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	//DefaultCenterSvr = "127.0.0.1:3600"
	DefaultCenterSvr      = "122.152.250.74:3600"
	ForwareTunnelConn     = byte(1)
	ProxyTunnelConn       = byte(2)
	ForwardInstanceConn   = byte(3)
	ProxyInstanceConn     = byte(4)
	KickForwardTunnelConn = byte(20)
	KickProxyTunnelConn   = byte(21)
	AddForwardLink        = byte(31)
	ForwardLinkSuccess    = byte(32)
	SetTunnelKey          = byte(33)
	NotifyMessage         = byte(200)
	FwNotifyMessage       = byte(201)
	ConnectPing			  = byte(255)
)

func int32ToBytes(i int) []byte {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(i))
	return buf
}

func WriteString(writer io.Writer, str string) {
	bufBytes := []byte(str)
	bufLen := len(bufBytes)
	writer.Write(int32ToBytes(bufLen))
	writer.Write(bufBytes)
}

func ReadString(reader io.Reader) (string, error) {
	lenBytes := make([]byte, 4)
	n, err := io.ReadFull(reader, lenBytes)
	if n != 4 {
		return "", err
	}
	l := int(binary.LittleEndian.Uint32(lenBytes))
	bufBytes := make([]byte, l)
	n, err1 := io.ReadFull(reader, bufBytes)
	if n != l {
		return "", err1
	}
	return string(bufBytes), nil
}

func WriteByte(writer io.Writer, v byte) (int, error) {
	return writer.Write([]byte{v})
}

func ReadByte(reader io.Reader) (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func IoCopy(sconn net.Conn, dconn net.Conn, exitChan chan bool) {
	if sconn != nil && dconn != nil {
		io.Copy(dconn, sconn)
	}
	exitChan <- true
}

func Ping(conn net.Conn){
	for {
		_, err := WriteByte(conn, ConnectPing)
		if err != nil {
			conn.Close()
			break
		}
		time.Sleep(time.Second * 2)
	}
}

func PrintError() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}
