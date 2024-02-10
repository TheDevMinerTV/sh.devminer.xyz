---
description: Update the GRUB default to the currently latest OEM D kernel and set up a prompt to notify the user when a new OEM kernel version is available.
tags:
  - 'Framework 13" AMD'
---
#!/bin/bash

set -xe

latest_oem_kernel=$(ls /boot/vmlinuz-* | grep '6.5.0-10..-oem' | sort -V | tail -n1 | awk -F'/' '{print $NF}' | sed 's/vmlinuz-//');
sudo sed -i.bak '/^GRUB_DEFAULT=/c\GRUB_DEFAULT="Advanced options for Ubuntu>Ubuntu, with Linux '"$latest_oem_kernel"'"' /etc/default/grub;
sudo update-grub
sudo apt install zenity

mkdir -p ~/.config/autostart
if [ ! -f ~/.config/autostart/kernel_check.desktop ]; then
  echo -e "[Desktop Entry]\nType=Application\nExec=bash -c \"latest_oem_kernel=\$(ls /boot/vmlinuz-* | grep '6.5.0-10..-oem' | sort -V | tail -n1 | awk -F'/' '{print \\\$NF}' | sed 's/vmlinuz-//') && current_grub_kernel=\$(grep '^GRUB_DEFAULT=' /etc/default/grub | sed -e 's/GRUB_DEFAULT=\\\"Advanced options for Ubuntu>Ubuntu, with Linux //g' -e 's/\\\"//g') && [ \\\"\\\${latest_oem_kernel}\\\" != \\\"\\\${current_grub_kernel}\\\" ] && zenity --text-info --html --width=300 --height=200 --title=\\\"Kernel Update Notification\\\" --filename=<(echo -e \\\"A newer OEM D kernel is available than what is set in GRUB. <a href='https://github.com/FrameworkComputer/linux-docs/blob/main/22.04-OEM-D.md'>Click here</a> to learn more.\\\")\"\nHidden=false\nNoDisplay=false\nX-GNOME-Autostart-enabled=true\nName[en_US]=Kernel check\nName=Kernel check\nComment[en_US]=\nComment=" \
    > ~/.config/autostart/kernel_check.desktop
fi
