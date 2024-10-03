#!/bin/bash

start() {
	if nc -zv localhost 5432 2>/dev/null && nc -zv localhost 9092 2>/dev/null; then
		return
	fi

	SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 || exit 1 ; pwd -P )"
	cd "$SCRIPTPATH" || exit 1

	# delete kafka data to fix restart after killing container
	rm -rf data/kafka data/zookeeper
	mkdir -p data/{kafka/data,postgres,zookeeper/{data,datalog}} || exit 1

	docker-compose --project-name torque up -d || exit 1

	echo -n "waiting for readiness"
	while : ; do
		if docker exec -i postgres-torque pg_isready -U postgres 2>&1 >/dev/null \
		&& docker exec -i kafka zookeeper-shell zookeeper:2181 ls /brokers/ids 2>&1 | grep "\[1\]" >/dev/null; then
			break
		fi
		echo -n .
		sleep 0.1
	done

	echo "done"
}

stop() {
    SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 || exit 1 ; pwd -P )"
    cd "$SCRIPTPATH" || exit 1
    docker-compose --project-name torque down || exit 1
}

restart() {
    stop || exit 1
    start || exit 1
}


cmd=$1
env=$2

case $cmd in
    start|stop|restart)
        $cmd
        ;;
    *)
        echo "usage: $0 start|stop|restart + local|{empty}"
        echo "if local env is determined, the current state of images is saved."
        echo "e.g., you will see the same data after restarting."
        echo "otherwise, if env is empty, the current state is lost."
        ;;
esac
