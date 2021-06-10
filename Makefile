all:
	$(MAKE) -C client
	$(MAKE) -C ws_lambda
	$(MAKE) -C http_lambda

deploy: all
	./deploy-ws.sh
	./deploy-http.sh
