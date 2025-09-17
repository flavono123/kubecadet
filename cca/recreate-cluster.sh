#!/usr/bin/env sh

kind delete cluster

kind create cluster --config=kind-config.yaml

cilium install --version 1.18.1
cilium status --wait

cilium connectivity test

cilium hubble enable
cilium status
# Hubble CLI 설치 후:
hubble status -P
hubble observe -P --last 5
