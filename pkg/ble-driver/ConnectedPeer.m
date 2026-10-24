// +build darwin,!noproximitytransport
//
//  ConnectedPeer.m
//  WeshBridgeDemo
//
//  Created by Rémi BARBERO on 29/04/2021.
//

#import "ConnectedPeer.h"

@implementation ConnectedPeer

- (BOOL)isClientReady {
    return self.client != nil;
}

- (BOOL)isServerReady {
    return self.server != nil;
}

@end
