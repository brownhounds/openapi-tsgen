#!/bin/sh

log_message_success() {
    local text="$1"
    echo "$text"
}

log_message_info() {
    local text="$1"
    echo "$text"
}

log_message_info "Fetching executable..."
wget https://github.com/brownhounds/openapi-tsgen/releases/download/v0.1.3/openapi-tsgen-linux-amd64

log_message_info "Installing..."
chmod +x ./openapi-tsgen-linux-amd64
bin_dir="$HOME/.local/bin"

# Check if the directory exists
if [ ! -d "$bin_dir" ]; then
    mkdir -p "$bin_dir"
fi

mv ./openapi-tsgen-linux-amd64 $HOME/.local/bin/openapi-tsgen

log_message_info 'Add following to .bashrc or .zshrc'
log_message_info 'export PATH="$HOME/.local/bin:$PATH"'
log_message_info 'Source the file or restart terminal session'

log_message_success 'All done.'
