#!/bin/sh
#
# Executed before the installation of the new package
#
#   $1=install              : On installation
#   $1=upgrade              : On upgrade

set -e

HEM_USER=hem
HEM_GROUP=hem
HEM_HOME="/var/lib/$HEM_USER"
RESTART_FLAG_FILE="/tmp/.restartHemOnUpgrade"

if [ "$1" = "install" ] || [ "$1" = "upgrade" ]; then
	if [ -d /run/systemd/system ] && /bin/systemctl status hem.service > /dev/null 2>&1; then
	  deb-systemd-invoke stop hem.service >/dev/null || true
	  touch "$RESTART_FLAG_FILE"
	fi
  if ! getent group "$HEM_GROUP" > /dev/null 2>&1 ; then
    addgroup --system "$HEM_GROUP" --quiet
  fi
  if ! getent passwd "$HEM_USER" > /dev/null 2>&1 ; then
    adduser --quiet --system --ingroup "$HEM_GROUP" \
    --disabled-password --shell /bin/false \
    --gecos "hem runtime user" --home "$HEM_HOME" "$HEM_USER"
    chown -R "$HEM_USER:$HEM_GROUP" "$HEM_HOME"
    adduser --quiet "$HEM_USER" dialout
  else
    adduser --quiet "$HEM_USER" dialout
    homedir=$(getent passwd "$HEM_USER" | cut -d: -f6)
    if [ "$homedir" != "$HEM_HOME" ]; then
      mkdir -p "$HEM_HOME"
      chown "$HEM_USER:$HEM_GROUP" "$HEM_HOME"
      process=$(pgrep -u "$HEM_USER") || true
      if [ -z "$process" ]; then
        usermod -d "$HEM_HOME" "$HEM_USER"
      else
        echo "--------------------------------------------------------------------------------"
        echo "Warning: hem's home directory is incorrect ($homedir)"
        echo "but can't be changed because at least one other process is using it."
        echo "Stop offending process(es), then restart installation."
        echo "Hint: You can list the offending processes using 'pgrep -u $HEM_USER -a'"
        echo "Note that you should NOT use the hem user as login user, since that will"
        echo "inevitably lead to this error."
        echo "in that case, please create a different user as login user."
        echo "--------------------------------------------------------------------------------"
        exit 1
      fi
    fi
  fi
fi

exit 0
