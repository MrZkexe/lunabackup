#Lunabackup
Lunabackup is a versatile backup tool designed to automate the backup process for both files and MariaDB databases. This README provides instructions on how to install and configure Lunabackup.


#Installation
To install Lunabackup, follow these steps:

Run the installer.sh script as the root user:
```bash
su -
./installer.sh
```
This script sets up a cron job to run Lunabackup every midnight and also executes upon completion.

#Configuration
If you wish to customize Lunabackup, locate the configuration file at the following path:

```bash
/etc/LunaBackup/lunaconf.json
```
Modify the settings in this JSON file according to your preferences.

#Backup
Lunabackup performs backups for both directories and files. The backup of MariaDB databases is a supported feature, and additional features and database support may be added in the future.

MariaDB Backup
Lunabackup backs up MariaDB databases automatically. If you need to back up other databases or files, consider contributing to the project and expanding its capabilities.

Feel free to reach out if you have any questions or encounter issues. Happy backing up!
