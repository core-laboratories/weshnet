package weshbridge

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/core-laboratories/weshnet/mobile/v2/go/internal/grpcutil"
	bridgepb "github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshbridge"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// The mobile API carries protobuf envelopes as base64 strings. Keeping this
// dispatcher here lets the app own routing, stream lifetime, and error handling
// without using the upstream bridge service implementation.
func (b *Bridge) invokeBridgeRPC(name, encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode bridge request: %w", err)
	}
	var reply proto.Message
	switch name {
	case "ClientInvokeUnary":
		var request bridgepb.ClientInvokeUnary_Request
		if err = proto.Unmarshal(data, &request); err == nil {
			reply, err = b.invokeUnary(&request)
		}
	case "CreateClientStream":
		var request bridgepb.ClientCreateStream_Request
		if err = proto.Unmarshal(data, &request); err == nil {
			reply, err = b.createStream(&request)
		}
	case "ClientStreamRecv":
		var request bridgepb.ClientStreamRecv_Request
		if err = proto.Unmarshal(data, &request); err == nil {
			reply, err = b.receiveStream(&request)
		}
	case "ClientStreamSend":
		var request bridgepb.ClientStreamSend_Request
		if err = proto.Unmarshal(data, &request); err == nil {
			reply, err = b.sendStream(&request)
		}
	case "ClientStreamClose":
		var request bridgepb.ClientStreamClose_Request
		if err = proto.Unmarshal(data, &request); err == nil {
			reply, err = b.closeStream(&request)
		}
	case "ClientStreamCloseAndRecv":
		var request bridgepb.ClientStreamCloseAndRecv_Request
		if err = proto.Unmarshal(data, &request); err == nil {
			reply, err = b.closeAndReceiveStream(&request)
		}
	default:
		return "", fmt.Errorf("unknown bridge method: %s", name)
	}
	if err != nil {
		return "", err
	}
	data, err = proto.Marshal(reply)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (b *Bridge) clientFor(desc *bridgepb.MethodDesc) (*grpcutil.LazyClient, error) {
	if desc == nil {
		return nil, errors.New("missing method description")
	}
	parts := strings.SplitN(strings.TrimPrefix(desc.Name, "/"), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, errors.New("invalid RPC method")
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, errors.New("bridge is closed")
	}
	client := b.clients[parts[0]]
	if client == nil {
		return nil, fmt.Errorf("service unavailable: %s", parts[0])
	}
	return client, nil
}

func outgoingContext(parent context.Context, headers []*bridgepb.Metadata) context.Context {
	values := metadata.MD{}
	for _, header := range headers {
		values.Append(header.GetKey(), header.GetValues()...)
	}
	return metadata.NewOutgoingContext(parent, values)
}

func bridgeMetadata(values metadata.MD) []*bridgepb.Metadata {
	result := make([]*bridgepb.Metadata, 0, len(values))
	for key, entries := range values {
		if key != "grpc-status-details-bin" {
			result = append(result, &bridgepb.Metadata{Key: key, Values: entries})
		}
	}
	return result
}

func bridgeError(err error) *bridgepb.Error {
	if err == nil {
		return &bridgepb.Error{}
	}
	return &bridgepb.Error{
		GrpcErrorCode: bridgepb.GRPCErrCode(status.Code(err)),
		Message:       err.Error(),
	}
}

func (b *Bridge) invokeUnary(request *bridgepb.ClientInvokeUnary_Request) (proto.Message, error) {
	client, err := b.clientFor(request.GetMethodDesc())
	if err != nil {
		return nil, err
	}
	desc := request.GetMethodDesc()
	if desc.GetIsClientStream() || desc.GetIsServerStream() {
		return nil, errors.New("stream method used as unary RPC")
	}
	trailer := metadata.MD{}
	reply, err := client.InvokeUnary(
		outgoingContext(b.ctx, request.GetHeader()),
		&grpcutil.LazyMethodDesc{Name: desc.GetName()},
		grpcutil.NewLazyMessage().FromBytes(request.GetPayload()),
		grpc.Trailer(&trailer),
	)
	result := &bridgepb.ClientInvokeUnary_Reply{Error: bridgeError(err), Trailer: bridgeMetadata(trailer)}
	if reply != nil {
		result.Payload = reply.Bytes()
	}
	return result, nil
}

func (b *Bridge) createStream(request *bridgepb.ClientCreateStream_Request) (proto.Message, error) {
	client, err := b.clientFor(request.GetMethodDesc())
	if err != nil {
		return nil, err
	}
	desc := request.GetMethodDesc()
	if !desc.GetIsClientStream() && !desc.GetIsServerStream() {
		return nil, errors.New("unary method used as stream RPC")
	}
	stream, err := client.InvokeStream(
		outgoingContext(b.ctx, request.GetHeader()),
		&grpcutil.LazyMethodDesc{
			Name:          desc.GetName(),
			ClientStreams: desc.GetIsClientStream(),
			ServerStreams: desc.GetIsServerStream(),
		},
		grpcutil.NewLazyMessage().FromBytes(request.GetPayload()),
	)
	result := &bridgepb.ClientCreateStream_Reply{Error: bridgeError(err)}
	if err == nil {
		result.StreamId = strconv.FormatUint(stream.ID(), 16)
		b.mu.Lock()
		if b.closed {
			b.mu.Unlock()
			stream.CancelFunc()
			return nil, errors.New("bridge is closed")
		}
		b.streams[result.StreamId] = stream
		b.mu.Unlock()
	}
	return result, nil
}

func (b *Bridge) streamFor(id string) (*grpcutil.LazyStream, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	stream := b.streams[id]
	if stream == nil {
		return nil, errors.New("invalid stream id")
	}
	return stream, nil
}

func (b *Bridge) forgetStream(id string) {
	b.mu.Lock()
	delete(b.streams, id)
	b.mu.Unlock()
}

func (b *Bridge) receive(id string, stream *grpcutil.LazyStream) *bridgepb.ClientStreamRecv_Reply {
	message := grpcutil.NewLazyMessage()
	err := stream.RecvMsg(message)
	if err != nil {
		b.forgetStream(id)
	}
	return &bridgepb.ClientStreamRecv_Reply{
		StreamId: id,
		Payload:  message.Bytes(),
		Error:    bridgeError(err),
		Trailer:  bridgeMetadata(stream.Trailer()),
	}
}

func (b *Bridge) receiveStream(request *bridgepb.ClientStreamRecv_Request) (proto.Message, error) {
	stream, err := b.streamFor(request.GetStreamId())
	if err != nil {
		return nil, err
	}
	return b.receive(request.GetStreamId(), stream), nil
}

func (b *Bridge) sendStream(request *bridgepb.ClientStreamSend_Request) (proto.Message, error) {
	id := request.GetStreamId()
	stream, err := b.streamFor(id)
	if err != nil {
		return nil, err
	}
	err = stream.SendMsg(grpcutil.NewLazyMessage().FromBytes(request.GetPayload()))
	if err != nil {
		b.forgetStream(id)
	}
	return &bridgepb.ClientStreamSend_Reply{
		StreamId: id,
		Error:    bridgeError(err),
		Trailer:  bridgeMetadata(stream.Trailer()),
	}, nil
}

func (b *Bridge) closeStream(request *bridgepb.ClientStreamClose_Request) (proto.Message, error) {
	id := request.GetStreamId()
	stream, err := b.streamFor(id)
	if err != nil {
		return nil, err
	}
	b.forgetStream(id)
	return &bridgepb.ClientStreamClose_Reply{Error: bridgeError(stream.Close())}, nil
}

func (b *Bridge) closeAndReceiveStream(request *bridgepb.ClientStreamCloseAndRecv_Request) (proto.Message, error) {
	id := request.GetStreamId()
	stream, err := b.streamFor(id)
	if err != nil {
		return nil, err
	}
	if err := stream.CloseSend(); err != nil {
		return &bridgepb.ClientStreamCloseAndRecv_Reply{StreamId: id, Error: bridgeError(err)}, nil
	}
	reply := b.receive(id, stream)
	return &bridgepb.ClientStreamCloseAndRecv_Reply{
		StreamId: id,
		Payload:  reply.Payload,
		Error:    reply.Error,
		Trailer:  reply.Trailer,
	}, nil
}
