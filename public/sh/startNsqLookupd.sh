killall -9 nsqlookupd
nohup /www/server/nsq/nsq-1.2.1.linux-amd64.go1.16.6/bin/nsqlookupd >> nsqLookupd.log 2>&1 &