package workerman_go

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
)

type GatewayIdInfo struct {
	//ClientGatewayPort Port   //2字节*8=16位 (整型)
	ClientGatewayNum uint64 //8字节*8=64位，必须是唯一的

	ClientGatewayAddr string //4字节 或者16字节
	//弃用
	//IPType            IpType //1个字节
	//ClientGatewayIpv4 net.IP //4字节*8=32位
	//ClientGatewayIpv6 net.IP //16字节*8=128位

}

// GenerateGatewayClientId 生成ClientToken Auth:GPT-3.5-turbo
func (c *GatewayIdInfo) GenerateGatewayClientId() string {

	var buf bytes.Buffer
	//binary.Write(&buf, binary.BigEndian, c.ClientGatewayPort)
	binary.Write(&buf, binary.BigEndian, c.ClientGatewayNum)

	buf.WriteString(c.ClientGatewayAddr)
	//binary.Write(&buf, binary.BigEndian, c.ClientGatewayAddr)

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// ParseGatewayClientId 解析code Auth:GPT-3.5-turbo
func ParseGatewayClientId(hexBuff string) (*GatewayIdInfo, error) {

	hexBuf, hexErr := base64.StdEncoding.DecodeString(hexBuff)

	if hexErr != nil {
		return nil, fmt.Errorf("解析失败:%s", hexErr)
	}

	//Data, err := hex.DecodeString(hexBuf)
	//if err != nil {
	//	return nil, err
	//}
	c := &GatewayIdInfo{}
	buf := bytes.NewReader(hexBuf)
	//binary.Read(buf, binary.BigEndian, &c.ClientGatewayPort)
	binary.Read(buf, binary.BigEndian, &c.ClientGatewayNum)

	remainingBytes := buf.Len()
	addrBytes := make([]byte, remainingBytes)
	_, err := buf.Read(addrBytes)
	if err != nil {
		return nil, fmt.Errorf("读取 ClientGatewayAddr 失败:%s", err)
	}

	c.ClientGatewayAddr = string(addrBytes)
	return c, nil
}

// genPrimaryKeyUint64 获取唯一的map[key]，注意自己设置读写锁
func genPrimaryKeyUint64(mapData map[uint64]InterfaceConnection) uint64 {

	for {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := mapData[num.Uint64()]; !exist {

			return num.Uint64()
		}
	}
	return 0
}
