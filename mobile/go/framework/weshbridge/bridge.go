// Package weshbridge is the mobile binding package used by the existing
// Android and iOS clients. Tone owns this bridge implementation; the account
// and messenger services it hosts remain separate upstream dependencies.
package weshbridge

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"

	manet "github.com/multiformats/go-multiaddr/net"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/text/language"
	"google.golang.org/grpc"

	"github.com/core-laboratories/weshnet/mobile/v2/go/internal/grpcutil"
	"github.com/core-laboratories/weshnet/mobile/v2/go/internal/notification"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/accounttypes"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/bertyaccount"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/mdns"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/osversion"
	weshgrpcutil "github.com/core-laboratories/weshnet/v2/pkg/grpcutil"
	"github.com/core-laboratories/weshnet/v2/pkg/lifecycle"
	"github.com/core-laboratories/weshnet/v2/pkg/logutil"
	"github.com/core-laboratories/weshnet/v2/pkg/netmanager"
)

// Bridge owns the in-process account server and routes protobuf requests from
// the Flutter method channel to the account, messenger, and protocol clients.
type Bridge struct {
	// Keep the generated iOS ServiceClient protocol shape while routing calls
	// through Tone's implementation below.
	ServiceClient
	mu               sync.RWMutex
	closed           bool
	cancel           context.CancelFunc
	ctx              context.Context
	clients          map[string]*grpcutil.LazyClient
	streams          map[string]*grpcutil.LazyStream
	server           *grpc.Server
	listener         *weshgrpcutil.BufListener
	accountConn      *grpc.ClientConn
	serviceAccount   bertyaccount.Service
	lifecycleManager *lifecycle.Manager
	netManager       *netmanager.NetManager
	logger           *zap.Logger
}

var _ LifeCycleHandler = (*Bridge)(nil)

type serviceRegistry struct{ bridge *Bridge }

func (r *serviceRegistry) RegisterService(name string, connection *grpc.ClientConn) {
	r.bridge.registerService(name, connection)
}

func NewBridge(config *BridgeConfig) (*Bridge, error) {
	if config == nil {
		return nil, errors.New("missing bridge configuration")
	}
	ctx, cancel := context.WithCancel(context.Background())
	logger := logutil.NewNativeLogger("tonebridge")
	bridge := &Bridge{
		ctx:              ctx,
		cancel:           cancel,
		clients:          make(map[string]*grpcutil.LazyClient),
		streams:          make(map[string]*grpcutil.LazyStream),
		lifecycleManager: lifecycle.NewManager(lifecycle.StateActive),
		logger:           logger,
	}
	bridge.ServiceClient = &localServiceClient{bridge: bridge}

	if config.netDriver != nil {
		interfaces := &inet{net: config.netDriver, logger: logger}
		mdns.SetNetDriver(interfaces)
		manet.SetNetInterface(interfaces)
	}
	var locker sync.Locker = &noopNativeMDNSLockerDriver{}
	if runtime.GOOS == "android" && osversion.GetVersion().Major() >= 30 && config.mdnsLockerDriver != nil {
		locker = config.mdnsLockerDriver
	}
	if config.connectivityDriver != nil {
		if current := config.connectivityDriver.GetCurrentState(); current != nil && current.info != nil {
			bridge.netManager = netmanager.NewNetManager(*current.info)
		}
		config.connectivityDriver.RegisterHandler(bridge)
	}
	if bridge.netManager == nil {
		bridge.netManager = netmanager.NewNoopNetManager()
	}

	languages := make([]language.Tag, 0, len(config.languages))
	for _, value := range config.languages {
		if tag, err := language.Parse(value); err == nil {
			languages = append(languages, tag)
		}
	}
	var notifications notification.Manager = notification.NewLoggerManager(logger)
	if config.notifdriver != nil {
		notifications = newNotificationManagerAdapter(logger, config.notifdriver)
	}
	if config.lc != nil {
		config.lc.RegisterHandler(bridge)
	}

	account, err := bertyaccount.NewService(&bertyaccount.Options{
		AppRootDirectory:      config.AppRootDirPath,
		SharedRootDirectory:   config.SharedRootDirPath,
		MDNSLocker:            locker,
		NetManager:            bridge.netManager,
		Languages:             languages,
		ServiceClientRegister: &serviceRegistry{bridge: bridge},
		NotificationManager:   notifications,
		Logger:                logger,
		LifecycleManager:      bridge.lifecycleManager,
		BleDriver:             config.bleDriver,
		NBDriver:              config.nbDriver,
		Keystore:              config.keystoreDriver,
	})
	if err != nil {
		cancel()
		return nil, err
	}
	bridge.serviceAccount = account
	server := grpc.NewServer()
	accounttypes.RegisterAccountServiceServer(server, account)
	listener := weshgrpcutil.NewBufListener(bufListenerSize)
	connection, err := listener.NewClientConn(ctx)
	if err != nil {
		server.Stop()
		listener.Close()
		account.Close()
		cancel()
		return nil, fmt.Errorf("connect account service: %w", err)
	}
	bridge.server = server
	bridge.listener = listener
	bridge.accountConn = connection
	for service := range server.GetServiceInfo() {
		bridge.registerService(service, connection)
	}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Warn("account server stopped", zap.Error(err))
		}
	}()
	return bridge, nil
}

type localServiceClient struct{ bridge *Bridge }

func (c *localServiceClient) InvokeBridgeMethod(method, payload string) (string, error) {
	b := c.bridge
	parts := strings.SplitN(strings.TrimPrefix(method, "/"), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", errors.New("invalid RPC method")
	}
	if parts[0] == "wesh.bridge.v1.BridgeService" {
		return b.invokeBridgeRPC(parts[1], payload)
	}
	b.mu.RLock()
	client := b.clients[parts[0]]
	closed := b.closed
	b.mu.RUnlock()
	if closed {
		return "", errors.New("bridge is closed")
	}
	if client == nil {
		return "", fmt.Errorf("service unavailable: %s", parts[0])
	}
	request, err := grpcutil.NewLazyMessage().FromBase64(payload)
	if err != nil {
		return "", fmt.Errorf("decode RPC request: %w", err)
	}
	reply, err := client.InvokeUnary(b.ctx, &grpcutil.LazyMethodDesc{Name: method}, request)
	if err != nil {
		return "", err
	}
	return reply.Base64(), nil
}

func (c *localServiceClient) InvokeBridgeMethodWithPromiseBlock(promise PromiseBlock, method, payload string) {
	go func() {
		reply, err := c.InvokeBridgeMethod(method, payload)
		if err != nil {
			promise.CallReject(err)
			return
		}
		promise.CallResolve(reply)
	}()
}

// registerService is called by the account service when an account opens or
// switches. Replacing the client under a lock makes the next call use the new
// account without requiring a new native bridge.
func (b *Bridge) registerService(name string, connection *grpc.ClientConn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.closed {
		b.clients[name] = grpcutil.NewLazyClient(connection)
	}
}

func (b *Bridge) Close() error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil
	}
	b.closed = true
	b.clients = nil
	for _, stream := range b.streams {
		stream.CancelFunc()
	}
	b.streams = nil
	b.mu.Unlock()
	b.cancel()
	if b.server != nil {
		b.server.Stop()
	}
	if b.listener != nil {
		b.listener.Close()
	}
	if b.accountConn != nil {
		b.accountConn.Close()
	}
	if b.serviceAccount != nil {
		return b.serviceAccount.Close()
	}
	return nil
}

func (b *Bridge) HandleState(state int) {
	if state == AppStateActive {
		b.lifecycleManager.UpdateState(lifecycle.StateActive)
	} else if state == AppStateInactive || state == AppStateBackground {
		b.lifecycleManager.UpdateState(lifecycle.StateInactive)
	}
	b.lifecycleManager.WaitForTasks()
}

func (b *Bridge) HandleTask() LifeCycleBackgroundTask {
	return newBackgroundTask(b.logger, func(ctx context.Context) error {
		b.lifecycleManager.UpdateState(lifecycle.StateActive)
		defer b.lifecycleManager.UpdateState(lifecycle.StateInactive)
		return b.serviceAccount.WakeUp(ctx)
	})
}

func (b *Bridge) WillTerminate() { _ = b.Close() }

func (b *Bridge) HandleConnectivityUpdate(info *ConnectivityInfo) {
	if info != nil && info.info != nil {
		b.netManager.UpdateState(*info.info)
	}
}

const (
	LevelDebug = int(zapcore.DebugLevel)
	LevelInfo  = int(zapcore.InfoLevel)
	LevelWarn  = int(zapcore.WarnLevel)
	LevelError = int(zapcore.ErrorLevel)
)

func (b *Bridge) Log(level int, subsystem, message string) {
	b.logger.Named(subsystem).Log(zapcore.Level(level), message) //nolint:gosec
}
