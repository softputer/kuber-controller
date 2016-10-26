# kuber-controller

#### Intrduction

This repo is inspired by Rancher's ingress-controller which uses k8s ingress to implement Layer 7 loadbalancer. But in the real world, maybe we don't have third-party loadbalancer support. So we have to implement our own loadbalancer.



# Usage

This repo only implements Layer 4 loadbalancer by haproxy or nginx. You can easily develop your own Layer 7 based on this repo.

