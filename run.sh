echo "Install and start docker"
sudo snap install docker &&
sleep 5s

cd app

echo "Run container with docker-compose"
if ! [ $(id -u) = 0 ]; then
    echo "Run as user"
    sudo docker-compose up
else
    echo "Run as root"
    docker-compose up
fi
