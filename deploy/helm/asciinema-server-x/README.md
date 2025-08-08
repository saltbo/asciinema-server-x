# asciinema-server-x Helm Chart

A Helm chart to deploy the asciinema-server-x Go service and optional SPA.

## Values overview

- image.repository/tag: container image for the server
- service.type/port/annotations: Kubernetes Service exposure and optional metadata annotations
- ingress: optional public routing
- persistence: PVC for data under STORAGE_ROOT (/data by default)
- admin.*: Basic Auth for admin endpoints
- config: environment variables consumed by the server

## Quickstart

- Update `values.yaml` image.repository and tag.
- If you have an Ingress controller, set `ingress.enabled=true` and configure hosts.

## Install

helm upgrade --install asciinema ./deploy/helm/asciinema-server-x \
  --set image.repository=ghcr.io/your-org/asciinema-server-x \
  --set image.tag=latest

## Admin API

Admin Basic is provided via Secret keys `ADMIN_BASIC_USER` and `ADMIN_BASIC_PASS`.

## Storage

- The server writes user data under `STORAGE_ROOT`.
- Enable `persistence.enabled` to retain data with a PVC.
