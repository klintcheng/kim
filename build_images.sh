
docker build --pull --rm -f "Dockerfile_royal" -t dockerklint/kim_royal:v1.3 "."

docker build --pull --rm -f "Dockerfile_gateway" -t dockerklint/kim_gateway:v1.3 "."

docker build --pull --rm -f "Dockerfile_server" -t dockerklint/kim_server:v1.3 "."

docker push dockerklint/kim_royal:v1.3
docker push dockerklint/kim_gateway:v1.3
docker push dockerklint/kim_server:v1.3