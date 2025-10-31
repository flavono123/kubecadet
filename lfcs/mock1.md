## 1

Write the total CPU cores count into /opt/mock/1/cpu_cores.
Write the total memory in GB (rounded down) into /opt/mock/1/memory_gb.

## 2

Create a cronjob for user backup-user that runs /opt/backup/daily.sh every day at 3:00am.

## 3

1. Create a new user developer with primary group devteam and home directory /home/dev/developer
2. User developer should be able to run sudo systemctl restart nginx without password

## 4

Configure iptables rules:

1. Block incoming connections on port 8080
2. Allow incoming connections on port 3306 only from localhost (127.0.0.1)

## 5

Find the disk (/dev/vdc or /dev/vde) with less than 20% usage, format it with xfs, and mount it to /mnt/storage.

## 6

There is a process consuming high CPU. Find any process using more than 50% CPU, kill it, and write the process name to /opt/mock/6/killed_process.

## 7

1. Change ownership of /var/www/html to user www-data and group www-data
2. Set permissions to 755 for all directories and 644 for all files under /var/www/html

## 8

Ensure the postgresql service is enabled to start on boot and is currently running. Write the service status (active/inactive) to /opt/mock/8/service_status.

## 9

In file /var/log/app/access.log, count how many requests returned HTTP status 404 and write the count to /opt/mock/9/error_count.

## 10

Create a tar.gz archive of /opt/important-data and save it to /backup/important-data-$(date +%Y%m%d).tar.gz. Verify the archive contains all files by listing contents to /backup/archive_contents.txt.
