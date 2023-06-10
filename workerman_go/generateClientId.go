package workerman_go

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net"
)

const IpTypeV4 = IpType(4)
const IpTypeV6 = IpType(6)

type Port uint16
type GatewayNum uint64
type IpType uint8

type ClientToken struct {
	IPType            IpType     //1个字节
	ClientGatewayIpv4 net.IP     //4字节*8=32位
	ClientGatewayIpv6 net.IP     //16字节*8=128位
	ClientGatewayPort Port       //2字节*8=16位 (整型)
	ClientGatewayNum  GatewayNum //8字节*8=64位，必须是唯一的
}

// GenerateGatewayClientId 生成ClientToken Auth:GPT-3.5-turbo
func (c *ClientToken) GenerateGatewayClientId() string {

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint8(c.IPType))
	if c.IPType == 4 {
		binary.Write(&buf, binary.BigEndian, c.ClientGatewayIpv4.To4())
	} else if c.IPType == 6 {
		binary.Write(&buf, binary.BigEndian, c.ClientGatewayIpv6.To16())
	}
	binary.Write(&buf, binary.BigEndian, c.ClientGatewayPort)
	binary.Write(&buf, binary.BigEndian, c.ClientGatewayNum)
	return hex.EncodeToString(buf.Bytes())
}

// ParseGatewayClientId 解析code Auth:GPT-3.5-turbo
func ParseGatewayClientId(hexBuff string) (*ClientToken, error) {
	data, err := hex.DecodeString(hexBuff)

	if err != nil {

		return nil, err
	}
	c := &ClientToken{}
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.BigEndian, &c.IPType)
	if c.IPType == IpType(4) {
		ipv4 := make([]byte, 4)
		binary.Read(buf, binary.BigEndian, &ipv4)
		c.ClientGatewayIpv4 = net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3])
	} else if c.IPType == IpType(6) {
		ipv6 := make([]byte, 16)
		binary.Read(buf, binary.BigEndian, &ipv6)
		c.ClientGatewayIpv6 = net.IP(ipv6)
	}
	binary.Read(buf, binary.BigEndian, &c.ClientGatewayPort)
	binary.Read(buf, binary.BigEndian, &c.ClientGatewayNum)
	return c, nil
}
