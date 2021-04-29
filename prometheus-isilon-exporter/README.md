# Install chart 
```
helm upgrade --install prometheus-isilon-exporter ./prometheus-isilon-exporter/ --values prometheus-isilon-exporter/values.yaml --namespace prom-isilon --create-namespace
```