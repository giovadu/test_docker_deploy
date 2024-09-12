docker build -t pipe919/firebase_notifications_server .

docker run -e DB_USER=_ -e DB_PASSWORD=_ -e DB_HOST=_ -e DB_PORT=_ -e DB_DATABASE=\* -d pipe919/firebase_notifications_server
