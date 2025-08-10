#include <glib.h>

static GCancellable *toCancellable(void *p) { 
    return (G_CANCELLABLE(p));
}
