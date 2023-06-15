package workerman_go

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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
	Sign      [lenSign]byte
	TimeStamp [lenTimeStamp]byte //10位字符
	Json      []byte
}

const (
	lenSign      = 16
	lenTimeStamp = 10
)

func (g *GenerateComponentSign) ToByte() []byte {
	b := bytes.Buffer{}
	b.Write(g.Sign[:])
	b.Write(g.TimeStamp[:])
	b.Write(g.Json)
	return b.Bytes()
}

// GenerateSignJsonTime 签名函数 防止篡改,参数 data 为 key-value 键值对，不包含 sign 字段。最后效果  sign[16] + json{}
func (g *GenerateComponentSign) GenerateSignJsonTime(data any, secretKey string, funcTime func() time.Duration) ([]byte, error) {
	expireTime := time.Now().Unix() + int64(funcTime().Seconds())
	expireTimeString := fmt.Sprintf("%d", expireTime)
	g.TimeStamp = [10]byte([]byte(expireTimeString))

	//带格式的json字符串
	jsonDataStr, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return nil, jsonErr
	}

	// 按照 key-value 的顺序组装字符串
	dataStr := bytes.Buffer{}

	dataStr.Write(jsonDataStr)
	//凭借加密字符串
	dataStr.WriteString(secretKey)

	sign := md5.Sum(dataStr.Bytes())

	// 对签名MD5 [16]byte 进行 hex 编码
	hexSignStr := hex.EncodeToString(sign[:])

	g.Sign = [16]byte([]byte(hexSignStr))

	return []byte(hexSignStr), nil
}

// ParseAndVerifySignJsonTime 解析json字符串并验证签名和有效时间
func ParseAndVerifySignJsonTime(jsonStr string, secretKey string) (map[any]any, error) {

	// 解析json字符串
	mapSourceData := make(map[any]any)
	err := json.Unmarshal([]byte(jsonStr), &mapSourceData)
	if err != nil {
		return nil, err
	}
	mapSignStringData := make(map[string]string)
	for sourceK, sourceV := range mapSourceData {
		sk, errSk := toString(sourceK)
		if errSk != nil {
			return nil, errSk
		}
		sv, errSv := toString(sourceV)
		if errSv != nil {
			return nil, errSv
		}
		mapSignStringData[sk] = sv
	}

	// 验证签名
	sign, ok := mapSignStringData[ConstSignFieldName]
	if !ok {
		return nil, errors.New("sign field not found")
	}
	//delete(mapData, ConstSignFieldName)

	// 将 key-value 键值对排序
	keys := make([]string, 0, len(mapSignStringData))
	for k := range mapSignStringData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 按照 key-value 的顺序组装字符串
	dataStr := bytes.Buffer{}
	for i, key := range keys {
		//key是sign的时候不加入签名计算
		if key == ConstSignFieldName {
			continue
		}
		//去掉最后一个&
		if i != 0 && i != (len(key)-1) {
			dataStr.WriteString("&")
		}
		dataStr.WriteString(fmt.Sprintf("%s=%s", key, mapSignStringData[key]))
	}

	//凭借加密字符串
	dataStr.WriteString(secretKey)

	signBytes := md5.Sum(dataStr.Bytes())

	// 对签名MD5 [16]byte 进行 hex 编码
	hexSignStr := hex.EncodeToString(signBytes[:])

	if sign != hexSignStr {
		return nil, errors.New("sign verification failed")
	}

	// 验证有效时间
	expireTimeStr, ok := mapSignStringData[strconv.Itoa(ConstSignTokenTimeStamp)]
	if !ok {
		return nil, errors.New("expire time field not found")
	}
	expireTime, err := strconv.ParseInt(expireTimeStr, 10, 64)
	if err != nil {
		return nil, err
	}
	if time.Now().Unix() > expireTime {
		return nil, errors.New("token expired")
	}

	return mapSourceData, nil
}
