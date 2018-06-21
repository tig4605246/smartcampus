#!/bin/sh

echo "*/5 6	* * *	root	test /home/ntust/tunnel/cronjob_5m.sh" >> /etc/crontab

echo "#" >> /etc/crontab

echo "/home/ntust/smartcampus/start.sh startmeter" >> /etc/rc.local

mkdir /home/ntust/smartcampus/

cp ./smartcampus /home/ntust/smartcampus/

cp ./macFile /home/ntust/smartcampus/

cp ./start.sh /home/ntust/smartcampus/

cd /home/ntust/smartcampus/