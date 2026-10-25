//
//  WriteDataCache.h
//  WeshBridgeDemo
//
//  Created by Rémi BARBERO on 03/08/2021.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@class WeshDevice;

@interface WriteDataCache : NSObject

@property (nonatomic, strong, nonnull) WeshDevice *device;
@property (nonatomic, strong, nonnull) NSData *data;

- (instancetype __nonnull) initWithDevice:(WeshDevice *__nonnull)device withData:(NSData *__nonnull)data;

@end

NS_ASSUME_NONNULL_END
