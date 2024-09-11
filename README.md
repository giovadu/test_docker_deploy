docker build -t pipe919/dockerhub:firebase_notifications_server .
docker push pipe919/dockerhub:firebase_notifications_server
docker run --memory=48g --cpus=8 pipe919/dockerhub:firebase_notifications_server

REEMPLAZAR CONTENEDOR

docker ps
docker stop CONTAINER ID
docker rm CONTAINER ID
docker pull pipe919/firebase_notifications_server
docker run -d pipe919/firebase_notifications_server
