#include <gio/gio.h>
#include <glib.h>
#include <glib-object.h>

static GCancellable *toCancellable(void *p) { 
    return (G_CANCELLABLE(p));
}
