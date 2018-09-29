powershell "docker service rm $(docker service ls -q)"
powershell Start-Sleep -m 2000
docker swarm leave --force
