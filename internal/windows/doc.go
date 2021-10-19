// Package windows provides platform-dependent routines required to support the
// package goey.
// In particular, on WIN32, the goal is to fill in some missing APIs that are not
// provided by lxn's WIN32 binding.
// Anything found herein should be a candidate for upstreaming.
// Since the WIN32 naming convention  is also camel case, most of the functions
// in this package are named exactly as their C API counterpart.
//
// This package is intended for internal use.
//
// This package contains platform-specific details.
package windows
