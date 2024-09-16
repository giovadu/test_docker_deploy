docker build -t pipe919/firebase_notifications_server .
docker pull pipe919/firebase_notifications_server

docker run -e DB*USER=service_user -e DB*PASSWORD=Service@#!@2019 -e DB*HOST=useast.datanetcenter.com -e DB*PORT=3306 -e DB_DATABASE=tracker -d pipe919/firebase_notifications_server

DB_USER=service_user
DB_PASSWORD=Service@#.!@2019
DB_HOST=useast.datanetcenter.com
DB_PORT=3306
DB_DATABASE=tracker

# PORT=3000

docker run -e DB_USER=service_user -e DB_PASSWORD=Service@#.!@2019 -e DB_HOST=singapur.datanetcenter.com -e DB_PORT=3306 -e DB_DATABASE=telocalizo -d pipe919/firebase_notifications_server
