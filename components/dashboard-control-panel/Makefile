build:
	@echo "Building custom Grafana"
	@docker build -t grafana-colmena .

run:
	@echo "Running custom Grafana on localhost:3000"
	@docker run -d -p 3000:3000 --name grafana grafana-colmena:latest

clean:
	@echo "Stopping and removing container and image"
	-@docker stop grafana
	-@docker rm grafana
	-@docker rmi grafana-colmena
