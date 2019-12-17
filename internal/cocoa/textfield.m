#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GTextField : NSTextField <NSTextFieldDelegate>
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)controlTextDidChange:(NSNotification*)obj;
@end

@implementation GTextField

- (void)controlTextDidChange:(NSNotification*)obj {
	NSString* v = [self stringValue];
	// Drop const, not representable in Go type system
	textfieldOnChange( self,
	                   (char*)[v cStringUsingEncoding:NSUTF8StringEncoding] );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		textfieldOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		textfieldOnBlur( self );
	}
	return rc;
}

@end

@interface GPasswordField : NSSecureTextField <NSTextFieldDelegate>
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)controlTextDidChange:(NSNotification*)obj;
@end

@implementation GPasswordField

- (void)controlTextDidChange:(NSNotification*)obj {
	NSString* v = [self stringValue];
	// Drop const, not representable in Go type system
	textfieldOnChange( self,
	                   (char*)[v cStringUsingEncoding:NSUTF8StringEncoding] );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		textfieldOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		textfieldOnBlur( self );
	}
	return rc;
}

@end

void* textfieldNew( void* superview, char const* text, bool_t password ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( text );

	// Create the button
	NSTextField<NSTextFieldDelegate>* control = password 
		? [[GPasswordField alloc] init]
		: [[GTextField alloc] init];
	textfieldSetValue( control, text );
	[control setEditable:YES];
	//[control setUsesSingleLineMode:YES];
	[control setDelegate:control];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

bool_t textfieldIsEditable( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );

	return [(NSTextField*)handle isEditable];
}

bool_t textfieldIsPassword( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );

	return [(id)handle isKindOfClass:[GPasswordField class]];
}
char const* textfieldPlaceholder( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );

	NSString* text = [[(NSTextField*)handle cell] placeholderString];
	return [text cStringUsingEncoding:NSUTF8StringEncoding];
}

void textfieldSetEditable( void* handle, bool_t value ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );

	[(NSTextField*)handle setEditable:value];
}

void textfieldSetValue( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );
	assert( text );

	NSString* value = [[NSString alloc] initWithUTF8String:text];
	NSString* oldValue = [(NSTextField*)handle stringValue];
	if ( [value compare:oldValue] != NSOrderedSame ) {
		[(NSTextField*)handle setStringValue:value];
	}
	[value release];
}

void textfieldSetPlaceholder( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );
	assert( text );

	NSString* title = [[NSString alloc] initWithUTF8String:text];
	[[(NSTextField*)handle cell] setPlaceholderString:title];
	[title release];
}

char const* textfieldValue( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );

	NSString* text = [(NSTextField*)handle stringValue];
	return [text cStringUsingEncoding:NSUTF8StringEncoding];
}
