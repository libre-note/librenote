# LibreNote

Docker based deployment

```bash
# persistance data directory for sqlite
mkdir data
sudo chown -R 1000:1000 data

# run container in background
docker-compose up -d

# run migrations
docker-compose exec librenote /app/librenote migrate -p /app/migrations up
```
