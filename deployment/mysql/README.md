# install 
Create secret to store mysql password
```bash
kubectl create secret generic -n database database-secret \
--from-literal=mysql-root-password=${MYSQL_ROOT_PASSWORD} \
--from-literal=mysql-password=${MSYQL_CUSTOM_PASSWORD}
```

install mysql using helm
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install mysql bitnami/mysql -n database --create-namespace -f deployment/mysql/mysql.yaml
```

get root password
```bash
kubectl get secret --namespace database database-secret -o jsonpath="{.data.mysql-root-password}" | base64 -d
```

start a mysql client connect to mysql using root account
```bash
kubectl run mysql-client --rm --tty -i --restart='Never' --image  docker.io/bitnami/mysql:8.0.31-debian-11-r30 --namespace database --env MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD --command -- bash
```

## some trouble 
```text
mysqladmin: connect to server at 'localhost' failed
error: 'Access denied for user 'root'@'localhost' (using password: YES)'
```
for pvc not delete
[handled  address](https://github.com/bitnami/charts/issues/12132)