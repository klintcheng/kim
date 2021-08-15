docker build --pull --rm -f "Dockerfile_royal" -t dockerklint/kim_royal:latest "."

docker build --pull --rm -f "Dockerfile_gateway" -t dockerklint/kim_gateway:latest "."

docker build --pull --rm -f "Dockerfile_server" -t dockerklint/kim_server:latest "."