killall -9 nsqadmin
nohup /www/server/nsq/nsq-1.2.1.linux-amd64.go1.16.6/bin/nsqadmin --lookupd-http-address=127.0.0.1:4161 >> nsqAdmin.log 2>&1 &