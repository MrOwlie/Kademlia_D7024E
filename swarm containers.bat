docker swarm init
powershell Start-Sleep -m 2000
docker network create --driver=bridge --subnet=11.11.0.0/16 --gateway=11.11.0.254 --ip-range=11.11.0.0/24 --scope=swarm --internal --attachable kademlia_network
powershell Start-Sleep -m 1000
cmd /k docker stack deploy -c swarm-compose.yml kadem
