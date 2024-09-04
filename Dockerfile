# Utilizamos una imagen base de Windows
FROM golang:1.22.1 AS builder

# Establecemos el directorio de trabajo
WORKDIR /build

# Copiamos el código fuente
COPY . .

# Descargamos los módulos de Go
RUN go mod download

# Compilamos la aplicación
RUN go build -o ./notifcations_server

# Utilizamos una imagen base de debian12 para el contenedor final
FROM gcr.io/distroless/base-debian12

# Establecemos el directorio de trabajo en el contenedor
WORKDIR /app

# Exponemos el puerto 1909
EXPOSE 1909

# Copiamos el ejecutable de la aplicación desde la imagen del constructor
COPY --from=builder /build/notifcations_server ./notifcations_server

# Definimos el comando de inicio
CMD ["./notifcations_server"]