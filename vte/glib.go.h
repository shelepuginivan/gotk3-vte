#include <gio/gio.h>
#include <glib.h>
#include <glib-object.h>

static GCancellable *toCancellable(void *p) {
    return (G_CANCELLABLE(p));
}

static GMenuModel *toMenuModel(void *p) {
    return (G_MENU_MODEL(p));
}

static uint gpointerToUint(gpointer i) {
    return (GPOINTER_TO_UINT(i));
}

static gpointer uintToGpointer(uint i) {
    return (GUINT_TO_POINTER(i));
}

static glong uintToGlong(uint i) {
    return ((glong)i);
}

static gsize uintToGsize(uint i) {
    return ((gsize)i);
}

static gssize intToGssize(int i) {
    return ((gssize)i);
}
