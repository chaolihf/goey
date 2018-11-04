#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <assert.h>

@interface GNOPThread : NSThread
- (void)main;
@end

@implementation GNOPThread

- (void)main {
	// Do nothing.  This is a NOP thread.
	return;
}

@end

@interface GApplicationDelegate : NSObject <NSApplicationDelegate>
- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication*)sender;
@end

@implementation GApplicationDelegate

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication*)sender {
	// Will try to close all of the windows.
	NSArray* windows = [sender windows];
	assert( windows );

	int i;
	for ( i = 0; i < [windows count]; ++i ) {
		NSWindow* w = [windows objectAtIndex:i];
		assert( w );
		[w performClose:self];
	}

	// Only close if all of the windows are closed.
	return [[sender windows] count] == 0 ? NSTerminateNow : NSTerminateCancel;
}

@end

static void detachAThread() {
	TRACE();

	// We need to make sure that Cocoa is running multithreaded.  Otherwise,
	// use of autopool from other threads will not work propertly.  The notes
	// for NSAutoreleasePool indicate that we need to detach a thread to
	// cause this transition.
	NSThread* thread = [[GNOPThread alloc] init];
	[thread start];
	[thread release];
}

static void initApplication() {
	TRACE();

	static GApplicationDelegate* delegate = 0;
	if ( !delegate ) {
		delegate = [[GApplicationDelegate alloc] init];
	}

	NSString* quitString = [[NSString alloc] initWithUTF8String:"Quit "];
	NSString* qString = [[NSString alloc] initWithUTF8String:"q"];

	[NSApplication sharedApplication];
	//[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
	[NSApp setDelegate:delegate];

	id menubar = [[NSMenu new] autorelease];
	id appMenuItem = [[NSMenuItem new] autorelease];
	[menubar addItem:appMenuItem];
	[NSApp setMainMenu:menubar];
	id appMenu = [[NSMenu new] autorelease];
	id appName = [[NSProcessInfo processInfo] processName];
	id quitTitle = [quitString stringByAppendingString:appName];
	id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
	                                              action:@selector( terminate: )
	                                       keyEquivalent:qString] autorelease];
	[appMenu addItem:quitMenuItem];
	[appMenuItem setSubmenu:appMenu];
}

static NSAutoreleasePool* pool = 0;

void init() {
	TRACE();

	assert( !pool );
	assert( !NSApp );
	assert( [NSThread isMainThread] );

	detachAThread();
	assert( [NSThread isMultiThreaded] );

	// This is a global release pool that we will keep around.  This will still
	// cause leaks, but until we identify where the autoreleasepool is required,
	// this will get us running.
	pool = [[NSAutoreleasePool alloc] init];
	assert( pool );
	initApplication();
	assert( NSApp && ![NSApp isRunning] );
}

void run() {
	TRACE();

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && ![NSApp isRunning] );
	assert( pool );

	[NSApp activateIgnoringOtherApps:YES];
	[NSApp run];
}

@interface DoStop : NSObject
- (void)main;
@end

@implementation DoStop

- (void)main {
	TRACE();

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && [NSApp isRunning] );

	[NSApp stop:nil];
}

@end

void stop() {
	TRACE();

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && [NSApp isRunning] );

	if ( [NSApp.windows count] > 0 ) {
		// Want to post an action to the main thread.
		// This will allow any windows that are probably
		// still shutting down to complete that action.
		id thunk = [[DoStop alloc] init];
		[thunk performSelectorOnMainThread:@selector( main )
		                        withObject:nil
		                     waitUntilDone:NO];
		[thunk release];
	} else {
		[NSApp stop:nil];
	}
}

@interface DoThunk : NSObject
- (void)main;
@end

@implementation DoThunk

- (void)main {
	TRACE();

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && [NSApp isRunning] );

	callbackDo();
}

@end

void do_thunk() {
	TRACE();

	assert( [NSThread isMultiThreaded] );
	assert( ![NSThread isMainThread] );
	assert( NSApp );

	while ( ![NSApp isRunning] ) {
		[NSThread sleepForTimeInterval:0.001];
	}

	id thunk = [[DoThunk alloc] init];
	[thunk performSelectorOnMainThread:@selector( main )
	                        withObject:nil
	                     waitUntilDone:YES];
	[thunk release];
}

bool_t isMainThread( void ) {
	TRACE();

	return [NSThread isMainThread];
}
