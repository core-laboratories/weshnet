//
//  ConnectedPeer.h
//  WeshBridgeDemo
//
//  Created by Rémi BARBERO on 29/04/2021.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@class WeshDevice;

@interface ConnectedPeer : NSObject

@property (nonatomic, assign, nullable) WeshDevice *client;
@property (nonatomic, assign, nullable) WeshDevice *server;
@property (readwrite, getter=isConnected) BOOL connected;

- (BOOL)isClientReady;
- (BOOL)isServerReady;

@end

NS_ASSUME_NONNULL_END
