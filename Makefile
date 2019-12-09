.PHONY: build docker_build up clean seed

build:
	echo "Build ..."
	go build -o=./.build/txpost ./src/
	echo "Build done"

docker_build:
	docker build -t txpost/txpost .

up: docker_build
	docker-compose up

clean:
	docker-compose stop
	docker-compose rm -f

seed:
	./test/seed.sh
