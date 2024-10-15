---
description: Copy all SSH keys from a folder onto multiple hosts.
tags: [ssh]
---
#!/bin/bash

# You can use this script to copy all public keys from your local machine to multiple remote machines.
#
# Usage: ./copy_ssh_keys.sh <non-root username> <own public key> <folder with pub keys> <host1> [host2]...
#
# Author: DevMiner <devminer@devminer.xyz>

COLOR_RESET='\033[0m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_BLUE='\033[1;34m'
COLOR_MAGENTA='\033[1;35m'
COLOR_BRIGHT_MAGENTA='\033[1;95m'

USAGE="Usage: ./copy_ssh_keys.sh <non-root username> <own public key> <folder with pub keys> <host1> [host2]..."

if [ $# -lt 4 ]; then
	echo "$USAGE"
	exit 1
fi

USERNAME="$1"
shift
OWN_PUBLIC_KEY="$1"
shift
SSH_KEYS_FOLDER="$1"
shift
HOSTS="$*"

fmt_prefix() {
	printf "%b%s%b:" "$COLOR_BRIGHT_MAGENTA" "$1" "$COLOR_RESET"
}

fmt_value() {
	printf "%b%s" "$COLOR_YELLOW" "$1"
}

ADD_PUBKEY_FN="$(
	cat <<EOF
add_pubkey() {
    PUBKEY="\$1"
    PUBKEY_FINGERPRINT="\$(echo \$PUBKEY | ssh-keygen -l -f -)"
    DEST="\$HOME/.ssh/authorized_keys"

    if [ ! -f "\$DEST" ]; then touch "\$DEST"; fi

    OLD_IFS="\$IFS"
    IFS=\$'\n'
    for PK in \$(cat \$DEST); do
        FINGERPRINT="\$(echo \$PK | ssh-keygen -l -f -)"

        if [ "\$FINGERPRINT" = "\$PUBKEY_FINGERPRINT" ]; then
		    printf "%bThe public key '%b%s%b' already exists on this machine\n" "$COLOR_GREEN" "$COLOR_YELLOW" "\$PUBKEY_FINGERPRINT" "$COLOR_GREEN"
		    return 0
        fi
    done
    IFS="\$OLD_IFS"
    
    echo "\$PUBKEY" >> "\$DEST"
	printf "%bCopied public key '%b%s%b'\n" "$COLOR_GREEN" "$COLOR_YELLOW" "\$PUBKEY_FINGERPRINT" "$COLOR_GREEN"
}
EOF
)"

for HOST in "${HOSTS[@]}"; do
	PREFIX="$(fmt_prefix "$HOST")"
	CONN="${USERNAME}@${HOST}"
	CONN_LOG="$(fmt_value "$CONN")"

	printf "%s%b Fetching public key...%b\n" "$PREFIX" "$COLOR_MAGENTA" "$COLOR_RESET"
	ssh-keyscan -H "${HOST}" >>~/.ssh/known_hosts

	printf "%s%b Copying own SSH key '%s%b' to %s${COLOR_RESET}\n" "$PREFIX" "$COLOR_MAGENTA" "$(fmt_value "$OWN_PUBLIC_KEY")" "$COLOR_MAGENTA" "$CONN_LOG"
	PUBKEY="$(cat "$OWN_PUBLIC_KEY")"
	ssh "$CONN" "$ADD_PUBKEY_FN; mkdir -p ~/.ssh; add_pubkey \"$PUBKEY\""

	printf "%s%b Creating temporary directory${COLOR_RESET}\n" "$PREFIX" "$COLOR_MAGENTA"
	ssh "$CONN" "mkdir -p /tmp/ssh_keys"

	printf "%s%b Copying SSH keys into temporary folder${COLOR_RESET}\n" "$PREFIX" "$COLOR_MAGENTA"
	scp "${SSH_KEYS_FOLDER}"/* "${CONN}:/tmp/ssh_keys"

	printf "%s%b Merging SSH keys into %s%b's authorized_keys%b\n" "$PREFIX" "$COLOR_MAGENTA" "$CONN_LOG" "$COLOR_MAGENTA" "$COLOR_RESET"
	ssh "$CONN" "$ADD_PUBKEY_FN; for KEY in /tmp/ssh_keys/*.pub; do PUBKEY=\"\$(cat \"\$KEY\")\"; add_pubkey \"\$PUBKEY\"; done"

	printf "%s%b Removing temporary directory%b\n" "$PREFIX" "$COLOR_MAGENTA" "$COLOR_RESET"
	ssh "$CONN" "rm -rf /tmp/ssh_keys"

	printf "%s%b Opening SSH session to %s%b\n" "$PREFIX" "$COLOR_MAGENTA" "$CONN_LOG" "$COLOR_RESET"

	printf "%s%b Please run '%s%b'%b\n" "$PREFIX" "$COLOR_BLUE" "$(fmt_value "su")" "$COLOR_BLUE" "$COLOR_RESET"
	printf "%s%b Please run '%s%b'%b\n" "$PREFIX" "$COLOR_BLUE" "$(fmt_value "mkdir -p ~/.ssh; cp .ssh/authorized_keys ~/.ssh")" "$COLOR_BLUE" "$COLOR_RESET"

	ssh "$CONN"
done

for HOST in "${HOSTS[@]}"; do
	PREFIX="$(fmt_prefix "$HOST")"

	printf "%s%b Checking that root login works for root@%s with key %b%s%b\n" "$PREFIX" "$COLOR_MAGENTA" "${HOST}" "$COLOR_YELLOW" "$OWN_PUBLIC_KEY" "$COLOR_RESET"

	printf "%s%b Running '%bwhoami%b': %b" "$PREFIX" "$COLOR_MAGENTA" "$COLOR_YELLOW" "$COLOR_MAGENTA" "$COLOR_GREEN"
	ssh "root@${HOST}" "whoami"
	printf "%b\n" "${COLOR_RESET}"
done