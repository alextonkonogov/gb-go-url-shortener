APP=url_shortener
REMOTEHOST=c8555@h4.netangels.ru
DIR=:/gbt.alextonkonogov.ru

clear:
	rm -f ${APP} || true

build: clear
	env GOOS=linux GOARCH=386 go build -v -o ${APP} cmd/url-shortener/main.go

kill:
	-ssh -l root ${REMOTEHOST} pkill -9 ${APP}

up-app:
	scp ${APP} ${REMOTEHOST}${DIR}

up-pub:
	scp -r public ${REMOTEHOST}${DIR}