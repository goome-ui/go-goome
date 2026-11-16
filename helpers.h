#ifndef GOOME_GO_HELPERS_H
#define GOOME_GO_HELPERS_H
#include "native/goome.h"
void ggo_label(goome_ui *u, const char *s);
void ggo_set_ease(goome_ui_config *c, int on, uintptr_t user);
const char *ggo_version(void);
#endif
