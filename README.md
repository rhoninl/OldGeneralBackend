# OldGeneralBackend

<div align="right">

[![Go Report Card](https://goreportcard.com/badge/github.com/leepala/OldGeneralBackend)](https://goreportcard.com/report/github.com/leepala/OldGeneralBackend)
[![codecov](https://codecov.io/github/leepala/OldGeneralBackend/branch/main/graph/badge.svg?token=0ZWSUO9ZST)](https://codecov.io/github/leepala/OldGeneralBackend)

</div>

Graduation design, A Goal management and scheduled punching system Backend

[Frontend Portal](https://github.com/leepala/OldGeneralFrontend)

## infra
### loki stack (loki/promtail) 
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install -f deployment/infra/prometheus-stack/values.yaml prometheus-stack prometheus-community/kube-prometheus-stack --version 45.7.1 -n monitoring
```
### prometheus stack (prometheus/grafana)
```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm upgrade --install loki grafana/loki-stack -n monitoring -f deployment/infra/loki-stack/values.yaml
```