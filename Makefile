build-push:
	docker buildx build --platform linux/amd64,linux/arm64 --push -t zpaceway/ohce:latest -t zpaceway/ohce:latest .
