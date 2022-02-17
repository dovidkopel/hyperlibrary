rm -rf /tmp/state-store /tmp/msp client/libclient
cd client && go build -o libclient && ./libclient $1