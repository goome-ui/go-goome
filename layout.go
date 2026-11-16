package goome

/*
#include "native/goome.h"
*/
import "C"

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func layoutTypes() map[string]reflect.Type {
	return map[string]reflect.Type{
		"goome_opts":      reflect.TypeOf(C.goome_opts{}),
		"goome_rect":      reflect.TypeOf(C.goome_rect{}),
		"goome_monitor":   reflect.TypeOf(C.goome_monitor{}),
		"goome_image":     reflect.TypeOf(C.goome_image{}),
		"goome_canvas":    reflect.TypeOf(C.goome_canvas{}),
		"goome_event":     reflect.TypeOf(C.goome_event{}),
		"goome_pacer":     reflect.TypeOf(C.goome_pacer{}),
		"goome_radii":     reflect.TypeOf(Radii{}),
		"goome_bitmap":    reflect.TypeOf(C.goome_bitmap{}),
		"goome_theme":     reflect.TypeOf(Theme{}),
		"goome_skin":      reflect.TypeOf(Skin{}),
		"goome_ui_item":   reflect.TypeOf(C.goome_ui_item{}),
		"goome_combo":     reflect.TypeOf(Combo{}),
		"goome_ui_config": reflect.TypeOf(UIConfig{}),
	}
}

func fieldName(f reflect.StructField) string {
	if c := f.Tag.Get("c"); c != "" {
		return c
	}
	return strings.TrimPrefix(f.Name, "_")
}

func layoutReport() string {
	var lines []string
	for name, t := range layoutTypes() {
		lines = append(lines, fmt.Sprintf("%s %d", name, t.Size()))
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Name == "_" {
				continue
			}
			lines = append(lines, fmt.Sprintf("%s.%s %d %d", name, fieldName(f), f.Offset, f.Type.Size()))
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}
