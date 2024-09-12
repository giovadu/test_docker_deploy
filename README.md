CREAR UN CONTENEDOR Y ENVIARLO
docker build -t pipe919/notifcations_api .
docker push pipe919/notifcations_api

REEMPLAZAR CONTENEDOR

docker pull pipe919/notifcations_api

EJECUTADOR UN CONTENEDOR POR CONSOLA

docker run -d -p 2122:2122 --restart always pipe919/notifcations_api
