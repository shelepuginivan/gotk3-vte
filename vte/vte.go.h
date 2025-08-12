#include <vte/vte.h>

static VteTerminal *toTerminal(void *p) {
    return (VTE_TERMINAL(p));
}

static gboolean isTerminal(void *p) {
    return (VTE_IS_TERMINAL(p));
}
