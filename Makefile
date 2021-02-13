NAME := weirdprogram-ams129
CONTAINER_CMD := docker

build:
	${CONTAINER_CMD} build -t ${NAME} .
run:
	${CONTAINER_CMD} run --rm -it \
	-v ${PWD}:/go/src/app \
	--dns 1.1.1.1 \
	-w /go/src/app \
	--name ${NAME} -d ${NAME} tail -f /dev/null

