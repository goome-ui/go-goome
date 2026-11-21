package goome

/*
#cgo CFLAGS: -std=c99 -O2
#cgo linux CFLAGS: -D_DEFAULT_SOURCE
#cgo windows CFLAGS: -D_CRT_SECURE_NO_WARNINGS
#define GOOME_IMPLEMENTATION
#include "native/goome.h"
#include "helpers.h"

extern double goomeGoEase(double t, uintptr_t user);

static double ggo_ease_call(double t, void *user)
{
    return goomeGoEase(t, (uintptr_t)user);
}

void ggo_label(goome_ui *u, const char *s)
{
    goome_ui_label(u, "%s", s);
}

void ggo_set_ease(goome_ui_config *c, int on, uintptr_t user)
{
    c->ease = on ? ggo_ease_call : 0;
    c->ease_user = on ? (void *)user : 0;
}

const char *ggo_version(void)
{
    return GOOME_VERSION;
}
*/
import "C"

// Version is the version of the goome C library this binding was built against.
var Version = C.GoString(C.ggo_version())
