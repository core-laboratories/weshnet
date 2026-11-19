package weshbridge

import (
	"encoding/base64"
	"testing"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/accounttypes"
	bridgepb "github.com/core-laboratories/weshnet/mobile/v2/go/pkg/bertybridge"
	"google.golang.org/protobuf/proto"
)

func TestToneBridgeRoutesAccountRPC(t *testing.T) {
	root := t.TempDir()
	config := NewBridgeConfig()
	config.SetAppRootDir(root)
	config.SetSharedRootDir(root)
	bridge, err := NewBridge(config)
	if err != nil {
		t.Fatal(err)
	}
	defer bridge.Close()
	request, err := proto.Marshal(&accounttypes.NetworkConfigGetPreset_Request{})
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := proto.Marshal(&bridgepb.ClientInvokeUnary_Request{
		MethodDesc: &bridgepb.MethodDesc{
			Name: "/wesh.account.v1.AccountService/NetworkConfigGetPreset",
		},
		Payload: request,
	})
	if err != nil {
		t.Fatal(err)
	}
	response, err := bridge.InvokeBridgeMethod(
		"/wesh.bridge.v1.BridgeService/ClientInvokeUnary",
		base64.StdEncoding.EncodeToString(envelope),
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := base64.StdEncoding.DecodeString(response)
	if err != nil {
		t.Fatal(err)
	}
	var routed bridgepb.ClientInvokeUnary_Reply
	if err := proto.Unmarshal(data, &routed); err != nil {
		t.Fatal(err)
	}
	if routed.GetError().GetGrpcErrorCode() != bridgepb.GRPCErrCode_OK {
		t.Fatalf("account RPC failed: %s", routed.GetError().GetMessage())
	}
	var reply accounttypes.NetworkConfigGetPreset_Reply
	if err := proto.Unmarshal(routed.GetPayload(), &reply); err != nil {
		t.Fatal(err)
	}
	if reply.Config == nil {
		t.Fatal("account RPC returned no network configuration")
	}
	streamRequest, err := proto.Marshal(&bridgepb.ClientCreateStream_Request{
		MethodDesc: &bridgepb.MethodDesc{
			Name:           "/wesh.account.v1.AccountService/StreamLogfile",
			IsServerStream: true,
		},
		Payload: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	streamResponse, err := bridge.InvokeBridgeMethod(
		"/wesh.bridge.v1.BridgeService/CreateClientStream",
		base64.StdEncoding.EncodeToString(streamRequest),
	)
	if err != nil {
		t.Fatal(err)
	}
	streamData, err := base64.StdEncoding.DecodeString(streamResponse)
	if err != nil {
		t.Fatal(err)
	}
	var created bridgepb.ClientCreateStream_Reply
	if err := proto.Unmarshal(streamData, &created); err != nil {
		t.Fatal(err)
	}
	if created.GetStreamId() == "" || created.GetError().GetGrpcErrorCode() != bridgepb.GRPCErrCode_OK {
		t.Fatalf("stream was not created: %s", created.GetError().GetMessage())
	}
	closeRequest, err := proto.Marshal(&bridgepb.ClientStreamClose_Request{StreamId: created.GetStreamId()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bridge.InvokeBridgeMethod(
		"/wesh.bridge.v1.BridgeService/ClientStreamClose",
		base64.StdEncoding.EncodeToString(closeRequest),
	); err != nil {
		t.Fatal(err)
	}
}
