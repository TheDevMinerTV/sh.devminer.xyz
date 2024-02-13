---
description: Allow both CPU and platform drivers to be simultaneously active.
tags:
  - 'Framework 13" AMD'
  - 'Ubuntu 22.04'
---

#!/bin/bash

set -xe

sudo add-apt-repository ppa:superm1/ppd
sudo apt update
sudo apt upgrade -y