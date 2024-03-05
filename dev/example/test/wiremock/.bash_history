ls
ls /
ls /source/
ce /source/
cd /source/
go build ./cmd/grpc2http/main.go 
ls
rm main 
cat Makefile 
which grpc2http
stat /go/bin/grpc2http
make install -C cmd/grpc2http/
stat /go/bin/grpc2http
ls /scripts/
CompileDaemon -color
CompileDaemon -color -h
#CompileDaemon -color -log-prefix -build=
pwd
ls
cat Makefile 
#CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command=''
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh' -h
ls 
ls pkg/
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
ls /etc/rsyslog.conf
cat /etc/rsyslog.conf
less /etc/rsyslog.conf
sudo service rsyslog start
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
systemctl stop rsyslog.service
supervisord -h
supervisord ctl
supervisord ctl -h
supervisord ctl stop rsyslog
supervisord ctl stop pid rsyslogd
supervisord ctl stop pid rsyslog
supervisord ctl pid rsyslog
supervisord ctl pid rsyslogd
supervisord ctl pid rsyslog.service
supervisord ctl stop rsyslog.service
supervisord ctl start rsyslog.service
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/entrypoint.sh'
whoami
user
uid
id
ls -al scripts/entrypoint.sh 
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/entrypoint.sh'
sudo rsyslogd 
sudo rsyslogd -h
man rsyslogd
cat /etc/supervisord/supervisord.conf 
ls /usr/sbin/rsyslogd
which rsyslogd
less /etc/rsyslog.conf
/imklog
less /etc/rsyslog.conf
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/entrypoint.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/entrypoint.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/entrypoint.sh'
docker compose -f build/docker-compose.yaml exec --help
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/entrypoint.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/entrypoint.sh'
echo $(./scripts/env.sh)
./scripts/env.sh
echo $SCRIPTS
echo $( dirname -- "${BASH_SOURCE[0]}" )
ls
source ./scripts/env.sh
echo $SCRIPTS/
echo $SCRIPTS
./scripts/init.sh 
source ./scripts/en
source ./scripts/env.sh 
ls ./scripts/proxy/
./scripts/proxy/install.sh 
source ./scripts/en
source ./scripts/env.sh 
ls ./scripts/proxy/
./scripts/proxy/install.sh 
source ./scripts/env.sh 
bash ./scripts/proxy/install.sh 
chown +x -r ./scripts/
lsof | grep 3010
lsof | head
lsof | grep 3010 | awk '$2'
lsof | grep 3010 | awk 'print $2'
lsof | grep 3010 | awk '{print $2}'
lsof | grep 3010 | awk '{print $2}' | sort | uniq | xargs kill
lsof | grep 3010
lsof | grep 3010 | awk '{print $2}' | sort | uniq | xargs kill
lsof | grep 3010 | awk '{print $2}' | sort | uniq | xargs -n1 kill
xargs -h
lsof | grep 3010 | awk '{print $2}' | sort | uniq | xargs kill -r
lsof | grep 3010 | awk '{print $2}' | sort | uniq | xargs -r kill 
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/entrypoint.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
lsof | { grep "${PROXY_INPUT_PORT}" || true } | awk '{print $2}' | sort | uniq | xargs -r kill
lsof | grep "${PROXY_INPUT_PORT}" || true | awk '{print $2}' | sort | uniq | xargs -r kill
lsof | grep "3010" || true | awk '{print $2}' | sort | uniq | xargs -r kill
lsof | {grep "3010" || true} | awk '{print $2}' | sort | uniq | xargs -r kill
lsof | { grep "3010" || true } | awk '{print $2}' | sort | uniq | xargs -r kill
lsof | { grep "3010" || true; } | awk '{print $2}' | sort | uniq | xargs -r kill
export PROXY_INPUT_PORT=3010
sudo lsof -t -i:${PROXY_INPUT_PORT}
lsof -t -i:${PROXY_INPUT_PORT}
lsof -t -i:${PROXY_INPUT_PORT}
lsof | grep 3010 | awk '{print $2}' 
lsof | grep 3010 
lsof | grep 3010 
lsof | grep 3010 
ps | grep init.sh
ps | grep init.sh
ps | grep '[.]/scripts/init.sh'
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | cut -f1
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | cut -d$'\t' -f1
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | cut -d$'\s*' -f1
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | cut -d$' ' -f1
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}'
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}' | xargs -r kill
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}'
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}'
ps | grep proxy
ps | grep proxy | sed 's/^.*0 //g'
ps | grep proxy | sed 's/^.*0 //g' | sort 
ps | grep proxy | sed 's/^.*0 //g' | sort | uniq -c
kill -h
man kill
#kill -P 
ps | grep init.sh
kill -P 65727
pkill -P 65727
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}'
ps | grep proxy | sed 's/^.*0 //g' | sort | uniq -c
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}'
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}' | xargs -r pkill -P
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}' | echo xargs -r pkill -P
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}' | xargs -r echo  pkill -P
ps | grep '[.]/scripts/init.sh' | grep -v CompileDaemon | awk '{print $1}' | xargs -r -n1 pkill -P
ps | grep proxy | sed 's/^.*0 //g' | sort | uniq -c
ps | grep script | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v CompileDaemon | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v CompileDaemon | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v CompileDaemon
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}'
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}' | xargs -r kill
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}'
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}' | xargs -r kill
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}' | xargs -r pkill -P
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}' | xargs -r -n1 pkill -P
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}' | xargs -r -n1 pkill -P
ps | grep [s]cript | grep -v CompileDaemon | awk '{print $1}' | xargs -r -n1 kill
ps | grep [s]cript | grep -v CompileDaemon 
ps | grep [s]cript 
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh'
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | xargs -r kill
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r kill
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r pkill -P
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 pkill -P
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 kill
ps | grep [s]cript
ps 
cat /var/log/supervisord/watcher-mocks.log
cat /var/log/supervisord/watcher-mocks.log.0 
cat /var/log/supervisord/mock-push-sender.log.0
ps 
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 pkill
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 -ISUB kill -SUB
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 -ISUB kill -SUB
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 -ISUB kill -- -SUB
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 -ISUB kill -- -SUB
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh'
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 -ISUB echo kill -- -SUB
kill -- -78508
kill -- 78508
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 -ISUB echo kill -- SUB
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r -n1 -ISUB kill -- SUB
ps | grep tail
ps | grep tail | awk '{print $1}'
ps | grep tail | awk '{print $1}' | xargs -r kill
ps | grep tail | awk '{print $1}' | xargs -r -n1 kill
killall tail
bash -c "killall tail; ./scripts/init.sh"
bash -c "killall tail; ./scripts/init.sh"
ps | grep '[s]cript'
ps | grep '[s]cript'
top
ps
ps | grep [s]cript | grep -v CompileDaemon | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}'
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}' | xargs -r kill
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep -e '[s]cript' -e 'tail' | grep -v 'CompileDaemon .*init.sh' | awk '{print $1}'
ps | grep -e '[s]cript' -e 'tail'
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
ps | grep [s]cript | grep -v 'CompileDaemon .*init.sh' | sed 's/^.*0 //g' | sort | uniq -c
chown +x +r ./scripts/
chown +x -R ./scripts/
chmod +x -R ./scripts/
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
make install -C cmd/grpc2http
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -h
CompileDaemon -h
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
sudo lsof | grep ${PROXY_INPUT_PORT} | awk '{print $2}' | sort | uniq | xargs -r kill
PROXY_INPUT_PORT="3010"; sudo lsof | grep "${PROXY_INPUT_PORT}"
echo $?
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash ./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash -c "killall tail; ./scripts/init.sh"'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='killall tail && ./scripts/init.sh'e
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='killall tail && ./scripts/init.sh'e
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash "killall tail && ./scripts/init.sh"'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash -c "killall tail && ./scripts/init.sh"'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='bash -c "killall tail && ./scripts/init.sh"'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
CompileDaemon -color -log-prefix -build='make install -C cmd/grpc2http' -command='./scripts/init.sh'
./dev/watch_wiremock.sh 
./dev/watch_wiremock.sh
./dev/watch_wiremock.sh
