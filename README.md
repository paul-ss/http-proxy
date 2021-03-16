# http-proxy

## build & run
docker image build -t proxy . <br>
docker container run -i -t --rm -p 8080:8080 -p 8000:8000 proxy:latest <br>

## info
ProxyServer running at 0.0.0.0:8000 <br>
ApiServer running at 0.0.0.0:8080 <br>
<br>
You can change these default settings in configs/config.go
and then you must rebuild you docker image to apply them.
<br> <br>

## api
/requests <br>
/requests/id <br>
/repeat/id <br>
/scan/id <br>


