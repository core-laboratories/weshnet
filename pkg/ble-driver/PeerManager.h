//
//  PeerManager.h
//  WeshBridgeDemo
//
//  Created by Rémi BARBERO on 29/04/2021.
//

#import <Foundation/Foundation.h>

#import "Logger.h"
#import "ConnectedPeer.h"

NS_ASSUME_NONNULL_BEGIN

@class WeshDevice;

@interface PeerManager : NSObject

@property (nonatomic, strong, nonnull) NSMutableDictionary *connectedPeers;
@property (nonatomic, strong, nonnull) Logger *logger;

- (instancetype __nonnull)initWithLogger:(Logger *__nonnull)logger;
- (ConnectedPeer *__nonnull)getPeer:(NSString *__nonnull) peerID;
- (ConnectedPeer *__nullable)registerDevice:(WeshDevice *__nonnull)device withPeerID:(NSString *__nonnull)peerID isClient:(BOOL)isClient;
- (void)unregisterDevice:(WeshDevice *)device;
- (void)removePeer:(NSString *__nonnull) peerID;
- (void)removeAllPeers;

@end

NS_ASSUME_NONNULL_END
