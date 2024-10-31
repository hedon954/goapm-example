genpb:
	protoc --go_out=./protos --go-grpc_out=./protos ./protos/*.proto

docker-up:
	make build-ubuntu
	docker compose -f zscripts/setup/docker-compose.yml up -d
	bash ./zscripts/ci/wait_mysql_start.sh
	mysql -h 127.0.0.1 -P 3306 -u root -p'root' < zscripts/setup/init.sql
	open http://127.0.0.1:16686
	open http://127.0.0.1:3000
	open http://127.0.0.1:9090
	open http://127.0.0.1:5601

docker-down:
	docker compose -f zscripts/setup/docker-compose.yml down

docker-restart:
	make docker-down
	make docker-up

setup:
	# docker compose  -f zscripts/setup/docker-compose.yml up -d
	# mysql -h 127.0.0.1 -P 3306 -u root -p'root' < zscripts/setup/init.sql

	make usr
	make sku
	make order

	ps -ef | grep usrsvc
	ps -ef | grep skusvc
	ps -ef | grep ordersvc

usr:
	APP_NAME=usrsvc go run usrsvc/main.go > logs/usrsvc.log 2>&1 &

sku:
	APP_NAME=skusvc go run skusvc/main.go > logs/skusvc.log 2>&1 &

order:
	APP_NAME=ordersvc go run ordersvc/main.go > logs/ordersvc.log 2>&1 &

stop:
	lsof -i :30001 | grep "main" | awk '{print $$2}' | xargs kill
	lsof -i :30002 | grep "main" | awk '{print $$2}' | xargs kill
	lsof -i :30003 | grep "main" | awk '{print $$2}' | xargs kill

	ps -ef | grep usrsvc
	ps -ef | grep skusvc
	ps -ef | grep ordersvc

restart:
	make stop
	make setup

status:
	ps -ef | grep usrsvc
	ps -ef | grep skusvc
	ps -ef | grep ordersvc
	lsof -i :30001
	lsof -i :30002
	lsof -i :30003
ab:
	bash zscripts/setup/ab.sh

build-ubuntu:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o zscripts/setup/build/usrsvc usrsvc/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o zscripts/setup/build/skusvc skusvc/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o zscripts/setup/build/ordersvc ordersvc/main.go
