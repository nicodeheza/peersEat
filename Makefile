test:
	go test ./...
up_peer:
	docker-compose -f docker-compose/docker-compose.${PEER}.yml up -d
	docker-compose -f docker-compose/docker-compose.${PEER}.yml logs -f peer${PEER}
rebuild_peer:
	docker-compose -f docker-compose/docker-compose.${PEER}.yml up -d --build
	docker-compose -f docker-compose/docker-compose.${PEER}.yml logs -f peer${PEER}
down_peer:
	docker-compose -f docker-compose/docker-compose.${PEER}.yml down