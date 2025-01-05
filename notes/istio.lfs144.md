---
id: 23h17up1bhg323nhuc4m7sx
title: Lfs144
desc: ''
updated: 1736061923586
created: 1735981566923
---

![Linux Foundation](https://img.shields.io/badge/LF-FreeCouse-003778.svg?logo=linux-foundation&labelColor=003778&link=https://training.linuxfoundation.org/courses/introduction-to-istio-lfs144)

## Milestone

- [x] 1/4(sat): 2
- [x] 1/4(sat): 3
- [x] 1/5(sun): 4
- [ ] 1/6(mon): 5
- [ ] 1/7(tue): 6
- [ ] 1/8(wed): 7
- [ ] 1/9(thu): 8
- [ ] 1/10(fri): 9
- [ ] 1/11(sat): 10

## 02. Overview

### Learning Objectives(2)

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

### Learning Objectives(3)

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

### Hands-on(3)

https://killercoda.com/lorenzo-g/scenario/istio-installation-istioctl

- meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY: only outbound traffic to ServiceEntry is allowed; https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig-OutboundTrafficPolicy-Mode

## 04. Observability

### Learning Objectives(4)

- metrics collection and observability of istio
- prometheus, promql, grafana, kiali
- distributed tracing
- (guess this is not exactly needed for the exam)

monitoring; profiling and stack tracing in a single application, monolith -> observability; distributed tracing, metrics and logs over multiple micro services
envoy sidecar collects metrics in a uniform way, no longer for each application/developer need to do that

### Hands-on(4)

no scenario, on playground: https://killercoda.com/lorenzo-g/scenario/playground

- prometheus scrape endpoint = metrics collection endpoint
- trace: end-to-end request-response flow; uid + spans
- span: a component of trace(e.g. call a service from another)
- tracing standards:
  - w3c [trace context](https://www.w3.org/TR/trace-context/)(otel)
  - x-b3(zipkin): trace id is propagated by this b3-header [b3-propagation](https://github.com/openzipkin/b3-propagation)
  - trace id gen counts to envoy(automatically); for new trace, envoy sidecar assigns a new trace id
  - propagating counts to app(configure): to propagate a trace id, application must set this context, headers, by including tracing client library
