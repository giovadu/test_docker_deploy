# Utilizamos una imagen base de Windows
FROM golang:1.23.1 AS builder

# Establecemos el directorio de trabajo
WORKDIR /build

# Copiamos el código fuente
COPY . .

# Descargamos los módulos de Go
RUN go mod download

# Compilamos la aplicación
RUN go build -o ./firebase_notifications_server

# Utilizamos una imagen base de debian12 para el contenedor final
FROM gcr.io/distroless/base-debian12

# Establecemos el directorio de trabajo en el contenedor
WORKDIR /app

# COPY .env .

COPY gd-notificacionesandroid-firebase-adminsdk-2v5rt-090a3f0a89.json .

# COPY notificaciones-push-1af7d-firebase-adminsdk-u0rd4-62ea0a3f0a.json .
# Exponemos el puerto 1909
EXPOSE 2122

# Copiamos el ejecutable de la aplicación desde la imagen del constructor
COPY --from=builder /build/firebase_notifications_server ./firebase_notifications_server

# Definimos el comando de inicio
CMD ["./firebase_notifications_server"]
