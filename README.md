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
```