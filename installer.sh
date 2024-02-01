#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "your user not is root."
  exit 1
fi

file_url="https://github.com/MrZkexe/lunabackup/releases/download/backup/lunabackup"
destination="/sbin/lunabackup"

wget "$file_url" -O "$destination"
if [ $? -eq 0 ]; then
  echo "Download sucesses."
else
  echo "Error Download"
  exit 1
fi

chmod +x "$destination"
(crontab -l 2>/dev/null; echo "0 0 * * * $destination") | crontab -

echo "Configuration completed. The script will run every day at midnight."
/sbin/lunabackup
