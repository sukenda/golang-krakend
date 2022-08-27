### Build docker image
```shell
docker build --build-arg ENV=dev -t gokrakend .
```

### Running 
```shell
docker run -p "8080:8080" -v $PWD:/etc/krakend/ gokrakend:latest run -c /etc/krakend/krakend.json
```