## Panop Kubernetes service resource import 

cp kservices in directory

```
export KUBECONFIG=/home/.kube/config
PANOP_ACCESS_KEY=xxx PANOP_HOST=tower.panop.io  terraform apply
```

## test binary 

```
echo '{"ns":"panop"}' | ./kservices
```
output
```json
{
  "service1.panop.svc": "ok",
  "service2.panop.svc": "ok"
}
```