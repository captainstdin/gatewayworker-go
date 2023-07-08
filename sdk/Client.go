package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gatewaywork-go/workerman_go"
	"io"
	"net/http"
	"time"
)

type Client struct {
	GatewayWorkerConfig workerman_go.ConfigGatewayWorker
}

func (s Client) SendToAll(data []byte, client_id_array []string, exclude_client_id []string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandSendToAll, workerman_go.GcmdSendToAll{
		Data:            data,
		ClientIdArray:   client_id_array,
		ExcludeClientId: exclude_client_id,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) SendToClient(client_id string, send_data string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandSendToClient, workerman_go.GcmdSendToClient{
		ClientId: client_id,
		SendData: send_data,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) CloseClient(client_id string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandCloseClient, workerman_go.GcmdCloseClient{
		ClientId: client_id,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) IsOnline(client_id string) int {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandIsOnline, workerman_go.GcmdIsOnline{
		ClientId: client_id,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return 0
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultIsOnline
	json.Unmarshal(buff, &parse)
	return parse.IsOnline
}

func (s Client) BindUid(client_id string, uid string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandBindUid, workerman_go.GcmdBindUid{
		ClientId: client_id,
		Uid:      uid,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())

}

func (s Client) UnbindUid(client_id string, uid string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandUnbindUid, workerman_go.GcmdUnbindUid{
		ClientId: client_id,
		Uid:      uid,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) IsUidOnline(uid string) int {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandIsUidOnline, workerman_go.GcmdIsUidOnline{
		Uid: uid,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return 0
	}
	buff := curlPostBinaryData(timeByte.ToByte())
	var parse workerman_go.GResultIsUidOnline
	json.Unmarshal(buff, &parse)
	return parse.IsUidOnline
}

func (s Client) GetClientIdByUid(uid string) []string {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetClientIdByUid, workerman_go.GcmdGetClientIdByUid{
		Uid: uid,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())
	var parse workerman_go.GResultGetClientIdByUid
	json.Unmarshal(buff, &parse)
	return parse.ClientIDList
}

func (s Client) GetUidByClientId(client_id string) string {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetUidByClientId, workerman_go.GcmdGetUidByClientId{
		ClientId: client_id,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return ""
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetUidByClientId
	json.Unmarshal(buff, &parse)
	return parse.UID
}

func (s Client) SendToUid(uid string, message string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandSendToUid, workerman_go.GcmdSendToUid{
		Uid:     uid,
		Message: message,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) JoinGroup(client_id string, group string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandJoinGroup, workerman_go.GcmdJoinGroup{
		ClientId: client_id,
		Group:    group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) LeaveGroup(client_id string, group string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandLeaveGroup, workerman_go.GcmdLeaveGroup{
		ClientId: client_id,
		Group:    group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) Ungroup(group string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandUngroup, workerman_go.GcmdUngroup{
		Group: group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) SendToGroup(group string, message string, exclude_client_id []string) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandSendToGroup, workerman_go.GcmdSendToGroup{
		Group:           group,
		Message:         message,
		ExcludeClientId: exclude_client_id,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) GetClientIdCountByGroup(group string) int {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetClientIdCountByGroup, workerman_go.GcmdGetClientIdCountByGroup{
		Group: group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return 0
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetClientIdCountByGroup
	json.Unmarshal(buff, &parse)
	return parse.ClientCount
}

func (s Client) GetClientSessionsByGroup(group string) map[string]workerman_go.SessionKv {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetClientSessionsByGroup, workerman_go.GcmdGetClientSessionsByGroup{
		Group: group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetClientSessionsByGroup
	json.Unmarshal(buff, &parse)
	return parse.ClientSessions
}

func (s Client) GetAllClientIdCount() int {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetAllClientIdCount, workerman_go.GcmdGetAllClientIdCount{}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return 0
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetAllClientIdCount
	json.Unmarshal(buff, &parse)
	return parse.ClientCount
}

func (s Client) GetAllClientSessions() map[string]workerman_go.SessionKv {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetAllClientSessions, workerman_go.GcmdGetAllClientSessions{}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetAllClientSessions
	json.Unmarshal(buff, &parse)
	return parse.ClientSessions
}

func (s Client) SetSession(client_id string, data workerman_go.SessionKv) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandSetSession, workerman_go.GcmdSetSession{
		ClientId: client_id,
		Data:     data,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())
}

func (s Client) UpdateSession(client_id string, data workerman_go.SessionKv) {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandUpdateSession, workerman_go.GcmdUpdateSession{
		ClientId: client_id,
		Data:     data,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return
	}
	curlPostBinaryData(timeByte.ToByte())

}

func (s Client) GetSession(client_id string) workerman_go.SessionKv {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetSession, workerman_go.GcmdGetSession{
		ClientId: client_id,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetSession
	json.Unmarshal(buff, &parse)
	return parse.Session
}

func (s Client) GetClientIdListByGroup(group string) []string {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetClientIdListByGroup, workerman_go.GcmdGetClientIdListByGroup{
		Group: group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetClientIdListByGroup
	json.Unmarshal(buff, &parse)
	return parse.ClientIDList
}

func (s Client) GetAllClientIdList() []string {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetAllClientIdList, workerman_go.GcmdGetAllClientIdList{}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetAllClientIdList
	json.Unmarshal(buff, &parse)
	return parse.ClientIDList
}

func (s Client) GetUidListByGroup(group string) []string {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetUidListByGroup, workerman_go.GcmdGetUidListByGroup{
		Group: group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetUidListByGroup
	json.Unmarshal(buff, &parse)
	return parse.UIDList
}

func (s Client) GetUidCountByGroup(group string) int {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetUidCountByGroup, workerman_go.GcmdGetUidCountByGroup{
		Group: group,
	}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return 0
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetUidCountByGroup
	json.Unmarshal(buff, &parse)
	return parse.UIDCount
}

func (s Client) GetAllUidList() []string {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetAllUidList, workerman_go.GcmdGetAllUidList{}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetAllUidList
	json.Unmarshal(buff, &parse)
	return parse.UIDList
}

func (s Client) GetAllUidCount() int {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetAllUidCount, workerman_go.GcmdGetAllUidCount{}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return 0
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetAllUidCount
	json.Unmarshal(buff, &parse)
	return parse.UIDCount
}

func (s Client) GetAllGroupIdList() []string {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetAllGroupIdList, workerman_go.GcmdGetAllGroupIdList{}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return nil
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetAllGroupIdList
	json.Unmarshal(buff, &parse)
	return parse.GroupIDList
}

func (s Client) GetAllGroupCount() int {
	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.GatewayCommandGetAllGroupCount, workerman_go.GcmdGetAllGroupCount{}, s.GatewayWorkerConfig.SignKey, getTime)
	if err != nil {
		return 0
	}
	buff := curlPostBinaryData(timeByte.ToByte())

	var parse workerman_go.GResultGetAllGroupCount
	json.Unmarshal(buff, &parse)
	return parse.GroupCount
}

func curlPostBinaryData(buffCmd []byte) []byte {
	// 创建一个 HTTP 请求的客户端
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
		Timeout: time.Second * 10, // 设置超时时间
	}

	// 创建一个 POST 请求

	req, err := http.NewRequest("POST", "http://example.com/api/endpoint", bytes.NewBuffer(buffCmd))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// 设置请求头等其他参数
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	all, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil
	}

	return all

}

func getTime() time.Duration {

	return time.Second * 60
}
