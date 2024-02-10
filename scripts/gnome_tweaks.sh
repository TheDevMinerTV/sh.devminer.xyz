---
description: Some tweaks to make GNOME more usable
tags:
  - GNOME
---
# Make Win+Right Click resize windows
gsettings set org.gnome.desktop.wm.preferences resize-with-right-button true

# Enable workspaces on all monitors
gsettings set org.gnome.mutter workspaces-only-on-primary false

# Disable emoji picker (Ctrl+.)
gsettings set org.freedesktop.ibus.panel.emoji hotkey "[]"