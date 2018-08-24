#!/bin/sh -xe

if [ -f "/etc/init.d/loom" ]; then
  update-rc.d loom defaults
  invoke-rc.d loom start
fi
if [ -f "/lib/systemd/system/loom.service" ]; then
  systemctl daemon-reload
  systemctl enable loom.service
  systemctl start loom.service
fi
