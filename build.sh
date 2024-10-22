#!/bin/bash

IMAGE_NAME="tutitoos/storage-api"
TAG="latest"

echo "Construyendo la imagen Docker '$IMAGE_NAME:$TAG'..."

docker build -t $IMAGE_NAME:$TAG .

if [ $? -eq 0 ]; then
    echo "Imagen '$IMAGE_NAME:$TAG' construida correctamente."

    echo "Subiendo la imagen a Docker Hub..."
    docker push $IMAGE_NAME:$TAG

    if [ $? -eq 0 ]; then
        echo "Imagen '$IMAGE_NAME:$TAG' subida correctamente a Docker Hub."
    else
        echo "Error al subir la imagen a Docker Hub."
    fi
else
    echo "Error al construir la imagen Docker."
fi
