# install
Create secret to store redis password
```bash
kubectl create secret generic -n database redis-secret \
--from-literal=redis-password=huangliandedianwei
```


```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install -n database redis bitnami/redis -f deployment/redis/redis.yaml
```

get redis password
```bash
kubectl get secret --namespace database redis-secret -o jsonpath="{.data.redis-password}" | base64 -d
```

```bash
kubectl run --namespace database redis-client --restart='Never'  --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:7.0.7-debian-11-r7 --command -- sleep infinity
```