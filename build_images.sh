git pull

docker build --pull --rm -f "Dockerfile_royal" -t kim_royal:latest "."

docker build --pull --rm -f "Dockerfile_gateway" -t kim_gateway:latest "."

docker build --pull --rm -f "Dockerfile_server" -t kim_server:latest "."