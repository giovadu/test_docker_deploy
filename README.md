docker build -t pipe919/dockerhub:notifcations_server .
docker push pipe919/dockerhub:notifcations_server
docker run --memory=128g --cpus=8 -p 1909:1909 pipe919/dockerhub:notifcations_server

REEMPLAZAR CONTENEDOR

docker ps
docker stop CONTAINER ID
docker rm CONTAINER ID
docker pull pipe919/dockerhub:notifcations_server
docker run -d -p 2122:2122 --restart always pipe919/dockerhub:notifcations_server
