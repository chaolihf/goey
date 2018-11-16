#ifndef GOEY_COCOA_H
#define GOEY_COCOA_H

#include <stdint.h>

// Cannot use std bool.  The builtin type _Bool does not play well with CGO.
// Need an alternate for the binding.

typedef unsigned bool_t;

typedef struct nssize_tag {
	int32_t width;
	int32_t height;
} nssize_t;

typedef struct nscolor_tag {
	uint8_t r, g, b, a;
} nscolor_t;

/* Event loop */
extern void init( void );
extern void run( void );
extern void do_thunk( void );
extern void stop( void );
extern bool_t isMainThread( void );

extern void trace( char const* func );
#ifdef NTRACE
#define TRACE() ( (void)0 )
#else
#define TRACEx) trace(__func__)
#endif

/* Window */
extern void* windowNew( char const* title, unsigned width, unsigned height );
extern void windowClose( void* handle );
extern nssize_t windowContentSize( void* handle );
extern void* windowContentView( void* handle );
extern void windowMakeFirstResponder( void* handle, void* control );
extern void windowSetContentSize( void* handle, int width, int height );
extern void windowSetMinSize( void* handle, int width, int height );
extern void windowSetIconImage( void* handle, void* nsimage );
extern void windowSetScrollVisible( void* handle, bool_t horz, bool_t vert );
extern void windowSetTitle( void* handle, char const* title );
extern char const* windowTitle( void* handle );

/* View */
extern void viewSetFrame( void* handle, int x, int y, int dx, int dy );
extern void viewClose( void* handle );

/* Control */
extern bool_t controlIsEnabled( void* handle );
extern void controlSetEnabled( void* handle, bool_t value );
extern nssize_t controlIntrinsicContentSize( void* handle );
extern bool_t controlMakeFirstResponder( void* handle );

/* Button */
extern void* buttonNew( void* superview, char const* title );
extern void* buttonNewCheck( void* window, char const* title, bool_t value );
extern void buttonPerformClick( void* handle );
extern bool_t buttonState( void* handle );
extern void buttonSetState( void* handle, bool_t checked );
extern char const* buttonTitle( void* handle );
extern void buttonSetTitle( void* handle, char const* title );

/* Decoration */
extern void* decorationNew( void* superview, nscolor_t fill, nscolor_t stroke,
                            nssize_t radius );
extern nssize_t decorationBorderRadius( void* control );
extern nscolor_t decorationFillColor( void* control );
extern nscolor_t decorationStrokeColor( void* control );
extern void decorationSetBorderRadius( void* control, nssize_t radius );
extern void decorationSetFillColor( void* control, nscolor_t fill );
extern void decorationSetStrokeColor( void* control, nscolor_t stroke );

/* HR */
extern void* hrNew( void* superview );

/* PopUpButton */
extern void* popupbuttonNew( void* superview );
extern void popupbuttonAddItem( void* control, char const* text );
extern void popupbuttonSetValue( void* control, int index );

/* ProgressIndicator */
extern void* progressNew( void* superview, double min, double value,
                          double max );
extern double progressMax( void* handle );
extern double progressMin( void* handle );
extern double progressValue( void* handle );
extern void progressUpdate( void* handle, double min, double value,
                            double max );

/* Slider */
extern void* sliderNew( void* superview, double min, double value, double max );
extern double sliderMax( void* handle );
extern double sliderMin( void* handle );
extern double sliderValue( void* handle );
extern void sliderUpdate( void* handle, double min, double value, double max );

/* TabView */
extern void* tabviewNew( void* superview );
extern void tabviewAddItem( void* control, char const* text );
extern void tabviewSelectItem( void* control, int index );
extern void* tabviewContentView( void* control, int index );
extern nssize_t tabviewContentInsets( void* control );

/* Text */
extern void* textNew( void* superview, char const* text );
extern int textAlignment( void* handle );
extern int textEightyEms( void* handle );
extern int textMinHeight( void* handle, int width );
extern int textMinWidth( void* handle );
extern void textSetText( void* handle, char const* text );
extern void textSetAlignment( void* handle, int align );
extern char const* textText( void* handle );

/* TextField */
extern void* textfieldNew( void* superview, char const* text );
extern bool_t textfieldIsEditable( void* handle );
extern char const* textfieldPlaceholder( void* handle );
extern void textfieldSetEditable( void* handle, bool_t value );
extern void textfieldSetValue( void* handle, char const* text );
extern void textfieldSetPlaceholder( void* handle, char const* text );
extern char const* textfieldValue( void* handle );

/* TextView */
extern void* textviewNew( void* superview, char const* text );
extern void textviewSetValue( void* handle, char const* text );

/* Image */
extern void* imageNewFromRGBA( uint8_t* imageData, int width, int height,
                               int stride );
extern void* imageNewFromGray( uint8_t* imageData, int width, int height,
                               int stride );
extern void imageClose( void* handle );

/* ImageView */
extern void* imageviewNew( void* superview, void* image );
extern void imageviewSetImage( void* control, void* image );

#endif
