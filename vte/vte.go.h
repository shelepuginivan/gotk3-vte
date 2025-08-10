#include <vte/vte.h>

static VteTerminal *toTerminal(void *p) {
    return (VTE_TERMINAL(p));
}
