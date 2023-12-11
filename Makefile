up:
	docker compose up -d --build
	cd frontend && yarn serve

lint:
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
	gofmt -s -w .
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --out-format colored-line-number -v

frontend-build:
	yarn add bootstrap timeago.js axios vue-native-websocket
	npm install sass-loader -D
	npm install node-sass -D