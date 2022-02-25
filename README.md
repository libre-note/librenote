# LibreNote

Libre(Free as in freedom) note is a note taking applications. A alternative to google keep.

```bash
mkdir data
sudo chown -R 1000:1000 data

docker run --rm -it -p 8000:8000 \
 -v $(pwd)/config.yml:/app/config.yml \
 -v $(pwd)/infrastructure/db/migrations/sqlite:/app/migrations \
 -v $(pwd)/data:/persist \
 hrshadhin/librenote:latest

docker exec container_id /app/librenote migrate -p /app/migrations up

docker-compose -f _deploy/docker/docker-compose.yml up -d
docker-compose -f _deploy/docker/docker-compose.yml \
  exec librenote /app/librenote migrate -p /app/migrations up
```
