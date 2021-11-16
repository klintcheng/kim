
docker build --pull --rm -f "Dockerfile_royal" -t dockerklint/kim_royal:v1.4 "."

docker build --pull --rm -f "Dockerfile_gateway" -t dockerklint/kim_gateway:v1.4 "."

docker build --pull --rm -f "Dockerfile_server" -t dockerklint/kim_server:v1.4 "."

docker build --pull --rm -f "Dockerfile_router" -t dockerklint/kim_router:v1.1 "."

docker push dockerklint/kim_royal:v1.4
docker push dockerklint/kim_gateway:v1.4
docker push dockerklint/kim_server:v1.4
docker push dockerklint/kim_router:v1.1