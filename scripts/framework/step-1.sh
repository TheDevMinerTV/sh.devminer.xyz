---
description: Install the OEM D kernel for Ubuntu 22.04.
tags:
  - 'Framework 13" AMD'
  - 'Ubuntu 22.04'
---
#!/bin/bash

set -xe

sudo apt update
sudo apt upgrade -y
sudo snap refresh
sudo apt-get install linux-oem-22.04d -y
