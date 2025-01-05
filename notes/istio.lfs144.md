---
id: 23h17up1bhg323nhuc4m7sx
title: Lfs144
desc: ''
updated: 1736057884730
created: 1735981566923
---

![Linux Foundation](https://img.shields.io/badge/LF-FreeCouse-003778.svg?logo=linux-foundation&labelColor=003778&link=https://training.linuxfoundation.org/courses/introduction-to-istio-lfs144)

## Milestone

- [x] 1/4(sat): 2
- [x] 1/4(sat): 3
- [ ] 1/5(sun): 4
- [ ] 1/6(mon): 5
- [ ] 1/7(tue): 6
- [ ] 1/8(wed): 7
- [ ] 1/9(thu): 8
- [ ] 1/10(fri): 9
- [ ] 1/11(sat): 10

## 02. Overview

### Learning Objectives

- what is the problem service mesh solve and how address
- design and architecture of istio

### New problems for microservices

- Service discovery
- Load balancing
  - rr, session affinity, weighted
- Service call handling
  - circuit breaker, timeout, retry
- Resilience
  - circuit breaker?, fallback
- Programming models
- Diagnosis and troubleshooting
- Resource utilization
- Automated testing
- Traffic management

sidecar injection, automated by mutating admission webhook

iptables: the proxy captures all traffic(packets) of the main container(on "init-container")
\+ more option: as cni plugin

identity(security): [SPIFFE](https://spiffe.io/)(Secure Production Identity Framework for Everyone) framework.(url form: `spiffe://<trust domain>/<workload identifier>`)
CSR > istio (xDS API) > Envoy proxy

configuring: on istio control plane; envoy just do(not for each workload sidecared), dynamic, not static(envoy can wihtout restarting by xDS API)

edge gateway: in/egress to/from mesh; contour, emissary-ingress, in/egress gateway

## 03. Installing Istio

### Learning Objectives

- helm chart
- istio operator api
- istio installation profiles

configuration profiles: [doc](https://istio.io/latest/docs/setup/additional-setup/config-profiles/) is up to date, no `istioctl profile list` command


| Component | default | demo | minimal | remote | empty | preview | ambient |
|-----------|---------|------|---------|--------|-------|---------|----------|
| **Core components** | | | | | | | |
| istio-egressgateway | | ✔ | | | | | |
| istio-ingressgateway | ✔ | ✔ | | | | ✔ | |
| istiod | ✔ | ✔ | ✔ | | | ✔ | ✔ |
| CNI | | | | | | | ✔ |
| Ztunnel | | | | | | | ✔ |

a profile is a custom resource of istio operator api itself

- global: profile name, root docker image path, image tags, namespace, revision, ...
- mesh configuration(meshConfig): controlplane components' things; access log format, log encoding, default proxy config, discovery selector, trust domain, ...
- component configuration(components): individual/additional(e.g. multiple in/egress gateway) components'(e.g. pilot, in/egress gateway, ...), kubernetes resources(e.g. cpu, memory, labels, annotations, replicas, ...)
-

helm: base(validating webhook, sa), istiod(controlplane, mutating webhook), gateway(in/egress gateway)

handson: https://killercoda.com/lorenzo-g/scenario/istio-installation-istioctl

- meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY: only outbound traffic to ServiceEntry is allowed; https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig-OutboundTrafficPolicy-Mode


