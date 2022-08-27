build:
	docker build --build-arg ENV=dev -t gokrakend .

run:
	docker run -p "8080:8080" -v $PWD:/etc/krakend/ gokrakend:latest run -c /etc/krakend/krakend.json

