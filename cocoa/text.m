#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* textNew( void* superview, char const* text ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( text );

	// Create the text view
	NSText* control = [[NSText alloc] init];
	[control setDrawsBackground:NO];
	textSetText( control, text );
	[control setEditable:NO];

	// Add the control as the view for the window
	[(NSView*)superview addSubview:control];

	return control;
}

int textAlignment( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );
	switch ( [(NSText*)handle alignment] ) {
	default:
	case NSTextAlignmentLeft:
		return 0;

	case NSTextAlignmentCenter:
		return 1;

	case NSTextAlignmentRight:
		return 2;

	case NSTextAlignmentJustified:
		return 3;
	}
}

void textSetText( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );
	assert( text );

	NSString* nsText = [[NSString alloc] initWithUTF8String:text];
	[(NSText*)handle setText:nsText];
	[nsText release];
}

void textSetAlignment( void* handle, int align ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	switch ( align ) {
	default:
	case 0:
		[(NSText*)handle setAlignment:NSTextAlignmentLeft];
		break;
	case 1:
		[(NSText*)handle setAlignment:NSTextAlignmentCenter];
		break;
	case 2:
		[(NSText*)handle setAlignment:NSTextAlignmentRight];
		break;
	case 3:
		[(NSText*)handle setAlignment:NSTextAlignmentJustified];
		break;
	}
}

char const* textText( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	NSString* text = [(NSText*)handle text];
	return [text cStringUsingEncoding:NSUTF8StringEncoding];
}
