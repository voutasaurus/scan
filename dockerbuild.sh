set -ex
GOOS=linux go build .
docker build -t scan .
rm scan
docker run scan
