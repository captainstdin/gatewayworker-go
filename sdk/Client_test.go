package sdk

import (
	"gatewaywork-go/workerman_go"
	"reflect"
	"testing"
	"time"
)

var Conf = workerman_go.ConfigGatewayWorker{
	RegisterListenAddr:             "0.0.0.0:1238",
	RegisterListenPort:             ":1238",
	TLS:                            false,
	TlsKeyPath:                     "",
	TlsPemPath:                     "",
	RegisterPublicHostForComponent: "127.0.0.1:1238",
	GatewayPublicHostForClient:     "",
	GatewayListenAddr:              ":2727",
	SkipVerify:                     false,
	SignKey:                        "da!!bskdhaskld#1238asjiocy89123",
}

func TestClient_BindUid(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
		uid       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.BindUid(tt.args.client_id, tt.args.uid)
		})
	}
}

func TestClient_CloseClient(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.CloseClient(tt.args.client_id)
		})
	}
}

func TestClient_GetAllClientIdCount(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "testCound", fields: fields{GatewayWorkerConfig: Conf}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetAllClientIdCount(); got != tt.want {
				t.Errorf("GetAllClientIdCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAllClientIdList(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetAllClientIdList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllClientIdList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAllClientSessions(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]workerman_go.SessionKv
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetAllClientSessions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllClientSessions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAllGroupCount(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetAllGroupCount(); got != tt.want {
				t.Errorf("GetAllGroupCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAllGroupIdList(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetAllGroupIdList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllGroupIdList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAllUidCount(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetAllUidCount(); got != tt.want {
				t.Errorf("GetAllUidCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAllUidList(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetAllUidList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllUidList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetClientIdByUid(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetClientIdByUid(tt.args.uid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientIdByUid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetClientIdCountByGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		group string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetClientIdCountByGroup(tt.args.group); got != tt.want {
				t.Errorf("GetClientIdCountByGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetClientIdListByGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		group string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetClientIdListByGroup(tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientIdListByGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetClientSessionsByGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		group string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]workerman_go.SessionKv
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetClientSessionsByGroup(tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientSessionsByGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetSession(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   workerman_go.SessionKv
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetSession(tt.args.client_id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetUidByClientId(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetUidByClientId(tt.args.client_id); got != tt.want {
				t.Errorf("GetUidByClientId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetUidCountByGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		group string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetUidCountByGroup(tt.args.group); got != tt.want {
				t.Errorf("GetUidCountByGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetUidListByGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		group string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.GetUidListByGroup(tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUidListByGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_IsOnline(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.IsOnline(tt.args.client_id); got != tt.want {
				t.Errorf("IsOnline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_IsUidOnline(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			if got := s.IsUidOnline(tt.args.uid); got != tt.want {
				t.Errorf("IsUidOnline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_JoinGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
		group     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.JoinGroup(tt.args.client_id, tt.args.group)
		})
	}
}

func TestClient_LeaveGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
		group     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.LeaveGroup(tt.args.client_id, tt.args.group)
		})
	}
}

func TestClient_SendToAll(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		data              []byte
		client_id_array   []string
		exclude_client_id []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.SendToAll(tt.args.data, tt.args.client_id_array, tt.args.exclude_client_id)
		})
	}
}

func TestClient_SendToClient(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
		send_data string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.SendToClient(tt.args.client_id, tt.args.send_data)
		})
	}
}

func TestClient_SendToGroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		group             string
		message           string
		exclude_client_id []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.SendToGroup(tt.args.group, tt.args.message, tt.args.exclude_client_id)
		})
	}
}

func TestClient_SendToUid(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		uid     string
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.SendToUid(tt.args.uid, tt.args.message)
		})
	}
}

func TestClient_SetSession(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
		data      workerman_go.SessionKv
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.SetSession(tt.args.client_id, tt.args.data)
		})
	}
}

func TestClient_UnbindUid(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
		uid       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.UnbindUid(tt.args.client_id, tt.args.uid)
		})
	}
}

func TestClient_Ungroup(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		group string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.Ungroup(tt.args.group)
		})
	}
}

func TestClient_UpdateSession(t *testing.T) {
	type fields struct {
		GatewayWorkerConfig workerman_go.ConfigGatewayWorker
	}
	type args struct {
		client_id string
		data      workerman_go.SessionKv
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Client{
				GatewayWorkerConfig: tt.fields.GatewayWorkerConfig,
			}
			s.UpdateSession(tt.args.client_id, tt.args.data)
		})
	}
}

func Test_curlPostBinaryData(t *testing.T) {
	type args struct {
		buffCmd []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := curlPostBinaryData(tt.args.buffCmd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("curlPostBinaryData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTime(t *testing.T) {
	tests := []struct {
		name string
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTime(); got != tt.want {
				t.Errorf("getTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
