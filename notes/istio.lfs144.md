---
id: 23h17up1bhg323nhuc4m7sx
title: Lfs144
desc: ''
updated: 1738379610859
created: 1735981566923
---

![LF144 Course](https://img.shields.io/badge/LF-Free_Course-003778.svg?logo=linux-foundation&labelColor=003778&link=https://training.linuxfoundation.org/courses/introduction-to-istio-lfs144)

## Milestone

- [x] 1/4(sat): 2
- [x] 1/4(sat): 3
- [x] 1/5(sun): 4
- [x] 1/5(sun): 5
- [x] 1/6(mon): 6
- [x] 1/7(tue): 7
- [ ] 1/12(sun): 8
- [ ] 1/14(tue): 9
- [x] 1/14(tue): 10

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

[killercoda installation scenario](https://killercoda.com/lorenzo-g/scenario/istio-installation-istioctl)

- meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY: only outbound traffic to ServiceEntry is allowed; [api](https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig-OutboundTrafficPolicy-Mode)

## 04. Observability

### Learning Objectives(4)

- metrics collection and observability of istio
- prometheus, promql, grafana, kiali
- distributed tracing
- (guess this is not exactly needed for the exam)

monitoring; profiling and stack tracing in a single application, monolith -> observability; distributed tracing, metrics and logs over multiple micro services
envoy sidecar collects metrics in a uniform way, no longer for each application/developer need to do that

### Hands-on(4)

no scenario, on [killercoda playground](https://killercoda.com/lorenzo-g/scenario/playground)

- prometheus scrape endpoint = metrics collection endpoint
- trace: end-to-end request-response flow; uid + spans
- span: a component of trace(e.g. call a service from another)
- tracing standards:
  - w3c [trace context](https://www.w3.org/TR/trace-context/)(otel)
  - x-b3(zipkin): trace id is propagated by this b3-header [b3-propagation](https://github.com/openzipkin/b3-propagation)
  - trace id gen counts to envoy(automatically); for new trace, envoy sidecar assigns a new trace id
  - propagating counts to app(configure): to propagate a trace id, application must set this context, headers, by including tracing client library

## 05. Traffic Management

### Learning Objectives(5)

- expose a service
- routing/(traffic get)routed
- lb: weighted, least conn, session affinity, ...
- resilience, failure injection, circuit breakering, ...
- ServiceEntry: an external service from the mesh

### Traffic Routing

- virtualservice: route on the request properties(e.g. weight, inject failures, delay, mirror, ...)
- destinationrule: after routing, configure how to reach the target service(e.g. outlier detection, load balancing, connection pool, tls, ...)
  - subsets by labels -> endpoints collection(=cluster); \<traffic-direction|port|subset|hostname\>
- serviceentry: register external service|api with istio's features

### Advanced Traffic Routing

- `match`: `uri`, `scheme`, `method`, `authority`(header), `headers`(kebab-case)
- matching method: exact, prefix, regex

### Rewriting and Redirecting

- `rewrite`, `redirect`(mutex with `destination`)

### Manipulating Headers

- request headers (`spec.http[*].headers`)
- response headers by destination(`spec.http[*].route[*].headers`)
- operators: `set`(overwrite), `remove`

### AND and OR Semantics

- AND: in same match rules
- OR: in splitted, a different element of match rules
- eval match rules in order of from top to bottom
  - if false, eval next
  - if no match(default), eval alway true
    - every traffic is routed
    - using as a fallback

### Service Resilience(virtualservice)

- timeout: envoy drop the reuqest(=response 408) over the set value
- retry
  - with conditions(`retryOn`) per timeout(`perTryTimeout`)
    - policy ref: [envoy `x-envoy-retry-on` header](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#x-envoy-retry-on)
  - a retrying endpoint is out of load balancing pool

### Circuit Breaking with Outlier Detection(destinationrule)

- passive health check: observe the health of endpoints, remove or eject unhealthy ones from load balancing pool
- health: consecutive failures, temporal success rate, latency, ...
- `connectionPool`
- `outlierDetection`
  - `maxEjectionPercent`: budget control
  - `baseEjectionTime`: ejected duration * how many times ejected
  - `interval`: checking interval for each pod

### Failure Injection

- delay: emulate a slow network or overloaded service
- abort: return a http error code to caller
- a vs' retry policy cannot retrigger a failure injection

### ServiceEntry

- to control all istio feature for a mesh-external service (`location: MESH_EXTERNAL`)
- or to furnish a cluster ip itself for a mesh-internal service (`location: MESH_INTERNAL, resolution: STATIC`)
- securing egress by `meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY`

### Hands-on(5)

[killercoda playground](https://killercoda.com/ica-scenarios/scenario/playground) with deleting bookinfo sample

```sh
alias i=istioctl
k delete -f https://raw.githubusercontent.com/istio/istio/release-${ISTIO_MINOR_VERSION}/samples/bookinfo/platform/kube/bookinfo.yaml
k get ns default -L istio-injection # check istio-injection label
```

#### Gateways

- (in|egress)gateway(envoy) <-> istio-system/istio-in|egressgateway((lb)svc)
- unfortunately, ingressgateway cannot get external ip in the playground
- ingressgateway -> gateway(hosts) -> virtualservice(hosts, gateways > route(destination host as a domain name of service), weight, match, redirect, mirror) > destinationrule(subsets) -> endpoints(cluster)
- a response header `server` is `istio-envoy`; the sidecar proxy of destination service's workload

#### Weight-Based Traffic Routing

- set weights on vs.spec.http[*].route[*].weight

## 06. Security

### Learning Objectives(6)

- authn, authz, access controls, identities in istio
- mTLS, TLS config for envoy(sidecar proxies, gateways)
- service, user authn
- issuing certificates

access control: a **principal** can do an **action** on an **object**.
authn: [SPIFFE](https://spiffe.io/)(Secure Production Identity Framework For Everyone) X.509 certificate to a kubernetes service account(e.g. jwt token)

- e.g. `spiffe://cluster.local/ns/default/sa/my-service-account`

mTLS: both client and server should provide certificates to each other; in traditional TLS, only server does

- after connection is established, a client(workload's) envoy proxy send to a server-side envoy proxy
  - proxies does, not the workload
- in handshaking, a caller check a secure naming, service account is authorized to the target service
- `PeerAuthentication`: for inbound traffic, `DestinationRule`: for outbound traffic

`PeerAuthentication` (`spec.mtls.mode`):

- `PERMISSIVE`(default): allow both
- `STRICT`: require mTLS
- `DISABLE`: disable mTLS
- `UNSET`: inherit from the parent(mesh, namespace), otherwise default(PERMISSIVE)

`DestinationRule` (`spec.trafficPolicy.tls.mode`):

- `ISTIO_MUTUAL`(default): mTLS by istio cert
- `DISABLE`: disable TLS
- `SIMPLE`: TLS(traditional) by server cert
- `MUTUAL`: mTLS, use key and cert

gatway tls: ingress(client outside, server mesh), egress(client mesh, server outside)

- `PASSTHROUGH`: do not terminate TLS, forward to a virtualservice of mathcing SNI
- `SIMPLE`: standard TLS
- `MUTUAL`: mTLS, `caCertificates` or `credentialName` required
- `AUTO_PASSTHROUGH`: no virtualservice(SNI service map) required, details are encoded in the SNI value
- `ISTIO_MUTUAL`: mTLS, using istio certs

authz: actions and objects of access control

`AuthorizationPolicy`:

- `from`(sources identities): principal(pa), request princial(ra),  namespace, ip block, remote ip block; each source is combined by AND
- `to`(request operation): hosts, ports, methods, paths; comb AND
- `when`(conditions): key as [istio attributes](https://istio.io/latest/docs/reference/config/security/conditions/)
- `action`: `ALLOW`, `DENY`, `CUSTOM`(custom configured in MeshConfig), `AUDIT`(for auditing requests, do not allow or deny)
  - default, allow all; but if ALLOW is configured, DENY all the others
  - eval order: CUSTOM > DENY > ALLOW

![flowchart](https://istio.io/latest/docs/concepts/security/authz-eval.svg)

provisioning identity(up-to-date)

![diagram](https://istio.io/latest/docs/concepts/security/id-prov.svg)

- proxy > (SDS req) > istio-agent > (CSR) > istiod > (key,cert) > istio-agent > (cached key,cert) > proxy
- iterate over key rotation
- a proxy only do SDS API call, istio-agent do actual things

### Hands-on(6)

[killercoda playground](https://killercoda.com/ica-scenarios/scenario/playground) with deleting bookinfo sample

```sh
alias i=istioctl
k delete -f https://raw.githubusercontent.com/istio/istio/release-${ISTIO_MINOR_VERSION}/samples/bookinfo/platform/kube/bookinfo.yaml
k delete gw bookinfo-gateway
k get ns default -L istio-injection # check istio-injection label
```

## 07. Extending the Mesh

### Learning Objectives(7)

- envoy, `EnvoyFilter`, Wasm Plugins and `WasmPlugin`
- json files: listners(what) -> routes(where) -> clusters(how) -> endpoints(can receive)
- listeners: ip + port(named network locations)
  - 0.0.0.0:15006 - inbound
  - 0.0.0.0:15001 - outbound
- envoy filters
  - filter chain: \<request\> (listener) ->  filter -> fileter -> filter -> (service)
  - listener filters: L4(e.g. TLS inspector fileter, if connection is TLS encrypted, extract TLS info)
  - network filters: deal with TCP packets(HTTP connection manager filter(HCM), rate limit filter, redis proxy filter, mongo proxy, connection limit filter, [all list](https://www.envoyproxy.io/docs/envoy/latest/configuration/listeners/network_filters/network_filters))
  - HTTP filters: L7(e.g. HTTP filter - optionally created by HCM)
- routes: the last of filter chain should be a route filter
- clusters: a group of upstream(server) hosts, similar to kube svc
- endpoints: a part of a cluster(ip address, hostname)
- inspecting: `istioctl proxy-config <listeners|routes|cluster|endpoints> ...`
- `EnvoyFilter`
  - `applyTo`: LISTENER, NETWORK_FILTER, HTTP_FILTER, VIRTUAL_HOST, ... [api](https://istio.io/latest/docs/reference/config/networking/envoy-filter/#EnvoyFilter-ApplyTo)
  - `match`: specify a more exact location
    - `context`: SIDECAR_INBOUND, SIDECAR_OUTBOUND, GATEWAY, ANY
    - `proxy`: proxy's properties match
    - `listener`: listener attributes match(port, name, filters in chain, ...)
    - `routeConfiguration`: route configuration attributes match(port name, number, gw name, virtual host name, ...)
    - `cluster`: cluster attributes match(name, ...)
  - `patch`:
    - `operation`: ADD, REMOVE, MERGE [api](https://istio.io/latest/docs/reference/config/networking/envoy-filter/#EnvoyFilter-Patch-Operation)
    - `value`: an actual patch, yaml
- custom filters: native c++ api, lua filters, wasm filters
- WASM:  a portable binary executed in a VM (.wasm) > `WasmPlugin` from istio 1.12(patched by `EnvoyFilter` before)
- SKIP `WasmPlugin`
  - not covered in the exam; actually this chapter, including filters, is not
  - tinygo + go sdk in the course is not recommended anymore

## 08. Advanced Topics

### Learning Objectives(8)

- [sidecar](https://istio.io/latest/docs/reference/config/networking/sidecar/) resources
- onboarding vm to mesh
- istio multicluster deployment scenarios

- sidecar: if not set, all sidecar proxy is conn

### Hands-on(8)

[killercoda playground](https://killercoda.com/ica-scenarios/scenario/playground) 'with' bookinfo sample

```sh
apt-get install -y bash-completion
eval $(istioctl completion bash)
alias i=istioctl
complete -F __start_istioctl i

k delete -f https://raw.githubusercontent.com/istio/istio/release-${ISTIO_MINOR_VERSION}/samples/bookinfo/platform/kube/bookinfo.yaml
k delete gw bookinfo-gateway
k get ns default -L istio-injection # check istio-injection label
```

## 10. Course Completion

- outlier detection = circuit breaker is passive health check; removing or ejecting unhealthy endpoints, which are OBSERVED, from load balancing pool
- endpoints: route to "clusters", and a cluster is a group of endpoints
