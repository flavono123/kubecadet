---
id: ntc2bx910r85cy4q3thlipb
title: Killersh
desc: ''
updated: 1737192729227
created: 1736989815631
---

## [Lorenzo Gironi](https://killercoda.com/lorenzo-g)

- [ ] Istio Playground
- [x] Traffic Management - Request Routing
- [x] Traffic Management - Fault Injection
- [x] Traffic Management - Traffic Shifting
- [x] Traffic Management - Circuit Breaking
- [x] Traffic Management - Traffic Mirroring
- [x] Ingress - Gateways
- [x] Ingress - Gateway Secure
- [x] Ingress - Gateway without TLS Termination
- [x] Ingress - Kubernetes
- [x] Egress - Accessing external services
- [x] Egress - External services TLS origination
- [x] Egress - Gateways
- [x] Egress - Gateway TLS Origination
- [x] Security - Authentication Policy - mTLS
- [x] Security - Authorization - HTTP traffic
- [x] Security - Authorization - JWT Token
- [x] Istio Installation using istioctl

---

## Settings and Tips

- istioctl

  ```sh
  # https://istio.io/latest/docs/setup/getting-started/#download
  source <(istioctl completion bash)
  alias i=istioctl
  ```

- vi setup

  ```sh
  alias v=vi
  v ~/.vimrc
  # file ~/.vimrc
  set et nu sw=2 sts=2 ts=2
  ```

  - master line delete, word wrapping
- tmux with vertical split

  ```sh
  tmux
  # C-b "
  ```

- validations are matters, should use `curl` well
  - default: `curl -s -X <METHOD> http://<uri>`
  - headers: `curl -s -X <METHOD> -H "<header>: <value>" http://<uri>`
  - query params: `curl -s -X <METHOD> http://<uri>?<param1>=<value1>&<param2>=<value2>`
  - test pod: `kubectl run test --image=nicolaka/netshoot -it --rm --restart=Never -- curl ...`
  - tls: `curl -s -X <METHOD> --cacert <ca-cert> https://<uri>`
  - mtls: `curl -s -X <METHOD> --cacert <ca-cert> --cert <client-cert> --key <client-key> https://<uri>`
- declarative strategy: make directory and file per question

```sh
mkdir 2 # for q2
v 2/dr.yaml
v 2/vs.yaml
k apply -f 2
# amend files and apply again
```

- check questions env(workload); not sure this is effective in the real exam too

```sh
k get po,svc
```

- beware on the array spec field after copy from the doc's example
  - e.g. virtualservice's spec.http is array of objects
- copy the spec field name, find in the doc, there are many manifest examples to copy
- many of fields are objects, such as Port; port.number, and Percentage; percentage.value, check the error message when apply is failed

---

## TODO on doc

- `VirtualService`
  - routing
  - fault injection; timeout
  - weighted routing
  - mirroring
- `DestinationRule`
  - circuit breaking
  - traffic policy, port level
- `Gateway`
  - ingress, to virtualservice
  - tls, mtls, termination/passthrough
  - egress, in virtualservice, to [`mesh` the default gateway](https://istio.io/latest/docs/reference/config/networking/virtual-service/#VirtualService-gateways) first, then to egress gateway
    - `ServiceEntry`
    - tls origination(service entry direct/gateway)
- `PeerAuthentication`
- `AuthorizationPolicy`
  - allow nothing = `spec: {}`
  - `from[*].source`, `to[*].operation`
- `RequestAuthentication` (jwt)
