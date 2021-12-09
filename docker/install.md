# Docker-Compose

## Launch VM

```sh
multipass launch --name primary --cpus 4
```

## Install Docker-Compose

```sh
# Get a shell on the vm.
multipass shell primary

# Install docker and pip.
sudo apt-get -qy update && sudo apt-get -qy full-upgrade
sudo apt-get -qy install python3-pip

# Install docker-compose.
sudo pip3 install docker-compose
docker-compose -v

# Add your user to the docker group.
sudo usermod -aG docker $USER

# Activate the changes to groups
newgrp docker

# Verify that you can run docker commands without sudo.
docker ps
docker run hello-world
docker system prune --all
```
