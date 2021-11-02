// Package windows provides platform-dependent routines required to support the
// package goey when targeting WIN32.
//
// In particular, the goal of this package is to fill in some missing APIs that
// are not provided by lxn's WIN32 binding.  Most functions found herein should
// be a candidate for upstreaming. Since the WIN32 naming convention is also
// camel case, most of the functions in this package are named exactly as their
// C API counterpart.
//
// Some functions provided by lxn's WIN32 binding are recreated in this package.
// Unfortunately, the functions provided by lxn eat error codes, which
// complicates proper error handling.  Certain critical functions have been
// rewritten to support error handling. See golang.org/x/sys/windows/mkwinsyscall
// for guidance on proper wrapping of WIN32 API calls.
//
// This package is intended for internal use.
//
// This package contains platform-specific details.
package windows
