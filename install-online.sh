#!/bin/sh
set -e

# Prevent execution if this script was only partially downloaded
{
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"
    [ "$ARCH" = "x86_64" ] && ARCH="amd64"
    [ "$ARCH" = "aarch64" ] && ARCH="arm64"

    if ! command -v curl >/dev/null 2>&1; then
        echo "Error: curl is required to install toofan."
        exit 1
    fi

    VERSION="1.0.1-online"
    echo "Fetching toofan-online v${VERSION}..."

    URL="https://github.com/vyrx-dev/toofan/releases/download/v${VERSION}/toofan_${VERSION}_${OS}_${ARCH}.tar.gz"

    if ! curl -sfL "$URL" -o /tmp/toofan.tar.gz; then
        echo "Error: Release not found at $URL"
        exit 1
    fi
    
    tar -xzf /tmp/toofan.tar.gz -C /tmp toofan
    chmod +x /tmp/toofan

    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    mv /tmp/toofan "$INSTALL_DIR/toofan"
    echo "Installed to $INSTALL_DIR/toofan"

    # Add ~/.local/bin to PATH if not already present
    case "$PATH" in
        *"$INSTALL_DIR"*) ;;
        *)
            SHELL_NAME="$(basename "$SHELL")"
            case "$SHELL_NAME" in
                bash)
                    SHELL_CONFIG="$HOME/.bashrc"
                    LINE='export PATH="$HOME/.local/bin:$PATH"'
                    ;;
                zsh)
                    SHELL_CONFIG="$HOME/.zshrc"
                    LINE='export PATH="$HOME/.local/bin:$PATH"'
                    ;;
                fish)
                    SHELL_CONFIG="$HOME/.config/fish/config.fish"
                    LINE='fish_add_path $HOME/.local/bin'
                    mkdir -p "$(dirname "$SHELL_CONFIG")"
                    ;;
                *)
                    SHELL_CONFIG=""
                    ;;
            esac

            if [ -n "$SHELL_CONFIG" ]; then
                if ! grep -qF '.local/bin' "$SHELL_CONFIG" 2>/dev/null; then
                    printf '\n# toofan\n%s\n' "$LINE" >> "$SHELL_CONFIG"
                    echo "Added ~/.local/bin to PATH in $SHELL_CONFIG"
                fi
            else
                echo "Unsupported shell: $SHELL_NAME"
                echo "Search: \"how to add to PATH on $SHELL_NAME\" and add ~/.local/bin"
            fi
            ;;
    esac

    export PATH="$INSTALL_DIR:$PATH"
    sleep 0.5
    toofan
}
