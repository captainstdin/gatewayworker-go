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

// GenerateSignJsonTime 签名函数,参数 data 为 key-value 键值对，不包含 sign 字段
func GenerateSignJsonTime(data any, secretKey string, funcTime func() time.Duration) ([]byte, error) {

	jsonStr, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return nil, jsonErr
	}

	mapData := make(map[string]string)

	json.Unmarshal(jsonStr, &mapData)

	expireTime := time.Now().Unix() + int64(funcTime().Seconds())
	mapData[strconv.Itoa(ConstSignTokenTimeStamp)] = fmt.Sprintf("%d", expireTime)

	// 将 key-value 键值对排序
	keys := make([]string, 0, len(mapData))
	for k := range mapData {
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
		dataStr.WriteString(fmt.Sprintf("%s=%s", key, mapData[key]))
	}

	//凭借加密字符串
	dataStr.WriteString(secretKey)

	sign := md5.Sum(dataStr.Bytes())

	// 对签名MD5 [16]byte 进行 hex 编码
	hexSignStr := hex.EncodeToString(sign[:])

	mapData["sign"] = hexSignStr

	marshal, err := json.Marshal(mapData)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

// ParseAndVerifySignJsonTime 解析json字符串并验证签名和有效时间
func ParseAndVerifySignJsonTime(jsonStr string, secretKey string) (map[string]string, error) {
	// 解析json字符串
	mapData := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &mapData)
	if err != nil {
		return nil, err
	}

	// 验证签名
	sign, ok := mapData[ConstSignFieldName]
	if !ok {
		return nil, errors.New("sign field not found")
	}
	//delete(mapData, ConstSignFieldName)

	// 将 key-value 键值对排序
	keys := make([]string, 0, len(mapData))
	for k := range mapData {
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
		dataStr.WriteString(fmt.Sprintf("%s=%s", key, mapData[key]))
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
	expireTimeStr, ok := mapData[strconv.Itoa(ConstSignTokenTimeStamp)]
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

	return mapData, nil
}
