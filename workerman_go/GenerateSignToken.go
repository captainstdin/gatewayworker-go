package workerman_go

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func toString(v any) (string, error) {

	switch v.(type) {
	case struct{}:
		marshal, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
		return "", errors.New(fmt.Sprintf("{ %+v } 暂不支持结构体", v))
	case string:
		return v.(string), nil
	case byte:
		return string(v.([]byte)), nil
	case int:
		return strconv.Itoa(v.(int)), nil
	case uint64:
		return strconv.FormatUint(v.(uint64), 64), nil
	default:
		return "", errors.New(fmt.Sprintf("{ %+v } 类型没有合适的转换器", v))
	}

}

type GenerateComponentSign struct {
	PackageLen uint32   //4字节 包体长度,不参与签名
	Sign       [16]byte //不参与签名
	TimeStamp  int64    //8字节，
	Cmd        int32    //4字节的指令
	Json       []byte   //不确定
}

func (g *GenerateComponentSign) ToString() string {
	return string(g.ToByte())
}

func (g *GenerateComponentSign) sumPackageLen(ExcludeJson bool) uint32 {
	//计算二进制长度
	b2 := bytes.Buffer{}
	binary.Write(&b2, binary.BigEndian, g.TimeStamp)
	binary.Write(&b2, binary.BigEndian, g.Sign)

	if !ExcludeJson {
		binary.Write(&b2, binary.BigEndian, g.Json)
	}

	binary.Write(&b2, binary.BigEndian, g.Cmd)

	return uint32(b2.Len())
}

// ToByte 生成通讯字节
func (g *GenerateComponentSign) ToByte() []byte {
	g.PackageLen = g.sumPackageLen(false)
	//4字节的包头 + 16字节的签名 + 8字节的unix时间戳(int64) + 2字节的指令 + n字节的json字符串
	b := bytes.Buffer{}
	binary.Write(&b, binary.BigEndian, g.PackageLen) //4字节的包头
	b.Write(g.Sign[:])                               //16字节的签名
	binary.Write(&b, binary.BigEndian, g.TimeStamp)  // 8字节的unix时间戳(int64)
	binary.Write(&b, binary.BigEndian, g.Cmd)        //4字节的指令
	b.Write(g.Json)                                  //n字节的json字符串
	return b.Bytes()
}

func (g *GenerateComponentSign) sumSign(secretKey string) [16]byte {
	//[]bye(sign签名) =  [8]byte(timeUnix)+[2]byte(Cmd)+[n]byte(json)+私钥
	ToBeSign := bytes.Buffer{}
	//时间戳
	binary.Write(&ToBeSign, binary.BigEndian, g.TimeStamp) //[8]byte(timeUnix)
	//指令
	binary.Write(&ToBeSign, binary.BigEndian, g.Cmd) //[4]byte(Cmd)
	//json内容
	ToBeSign.Write(g.Json) //[n]byte(json)
	ToBeSign.WriteString(secretKey)

	return md5.Sum(ToBeSign.Bytes())
}

// GenerateSignTimeByte 签名函数 防止篡改,参数 Data 为 key-value 键值对，不包含 sign 字段。
func GenerateSignTimeByte(Cmd int, data any, secretKey string, funcTime func() time.Duration) (*GenerateComponentSign, error) {

	//带格式的json字符串
	jsonDataStr, jsonErr := json.Marshal(data)

	if jsonErr != nil {
		return nil, jsonErr
	}

	expireTime := time.Now().Unix() + int64(funcTime().Seconds())

	gen := &GenerateComponentSign{}
	gen.TimeStamp = expireTime
	gen.Cmd = int32(Cmd)
	gen.Json = jsonDataStr

	// 对签名MD5 [16]byte 进行 hex 编码
	//hexSignStr := hex.EncodeToString(sign[:])
	gen.Sign = gen.sumSign(secretKey)
	return gen, nil
}

// ParseAndVerifySignJsonTime 解析json字符串并验证签名和有效时间
func ParseAndVerifySignJsonTime(dataByte []byte, secretKey string) (*GenerateComponentSign, error) {
	gen := GenerateComponentSign{}
	reader := bytes.NewReader(dataByte)

	minLen := int(gen.sumPackageLen(true))
	if reader.Len() <= minLen {

		return nil, errors.New("有效长度不足：" + strconv.Itoa(minLen))
	}

	//var length int32
	// 读取包头长度
	if err := binary.Read(reader, binary.BigEndian, &gen.PackageLen); err != nil {
		return nil, err
	}
	// 读取签名
	if _, err := reader.Read(gen.Sign[:]); err != nil {
		return nil, err
	}
	// 读取时间戳
	if err := binary.Read(reader, binary.BigEndian, &gen.TimeStamp); err != nil {
		return nil, err
	}
	// 读取指令
	if err := binary.Read(reader, binary.BigEndian, &gen.Cmd); err != nil {
		return nil, err
	}

	// 读取JSON数据 ,jsonLen 为json长度应该是  gen.PackageLen- (时间+签名+指令)
	jsonLen := gen.PackageLen - gen.sumPackageLen(false)
	gen.Json = make([]byte, jsonLen)
	if _, err := reader.Read(gen.Json); err != nil {
		return nil, err
	}

	if gen.Sign != gen.sumSign(secretKey) {
		return nil, errors.New("签名校验失败！sign error")
	}

	if gen.TimeStamp <= time.Now().Unix() {
		return nil, errors.New("时间已过期！Time has expired")
	}

	return &gen, nil
}
