docker swarm init
powershell Start-Sleep -m 2000
cmd /k docker stack deploy -c swarm-compose.yml kadem
