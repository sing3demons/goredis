#### Redis Configuration
```
bind 0.0.0.0
appendonly yes
SAVE ""
```


#### k6 home brew
```
brew install k6
k6 run scripts/test.js
```

#### k6 on docker-compose
```
docker compose run --rm k6 run /scripts/test.js
docker compose run --rm k6 run /scripts/test.js -u5 -d5s
docker compose run --rm  k6 run --out json=output.json --summary-export=summary-export.json /scripts/test.js
docker compose run --rm k6
```

```
curl -k https://localhost:8000/healthz  
```