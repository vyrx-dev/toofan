# Topic: Script Boilerplate
#!/bin/bash
set -euo pipefail

log() {
    echo "[$(date +%H:%M:%S)] $*"
}

main() {
    log "starting"
    setup
    log "done"
}

main "$@"
