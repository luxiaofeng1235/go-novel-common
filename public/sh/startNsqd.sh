killall -9 nsqd
nohup /www/server/nsq/nsq-1.2.1.linux-amd64.go1.16.6/bin/nsqd --lookupd-tcp-address=127.0.0.1:4160 >> nsqd.log 2>&1 &