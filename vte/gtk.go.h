#include <gtk/gtk.h>

static GtkWidget *toGtkWidget(void *p) {
    return (GTK_WIDGET(p));
}
