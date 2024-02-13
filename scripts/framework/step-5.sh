---
description: Fix bogus keypresses while in suspend.
tags:
  - 'Framework 13" AMD'
  - 'Ubuntu 22.04'
  - Optional
---
#!/bin/bash

set -xe

if [ ! -f /etc/udev/rules.d/20-suspend-fixes.rules ]; then
  echo "ACTION==\"add\", SUBSYSTEM==\"serio\", DRIVERS==\"atkbd\", ATTR{power/wakeup}=\"disabled\"" | sudo tee /etc/udev/rules.d/20-suspend-fixes.rules
fi