# test-github-repo

A throwaway repo for testing the shared Docker workflow in
[`leonidgrishenkov/github-actions`](../github-actions).

## What it builds

A tiny Go HTTP server (`app/main.go`) baked into a multi-stage, multi-arch
`scratch` image:

- `GET /` → `Hello from leonidgrishenkov/test-github-repo` (or `$GREETING`)
- `GET /healthz` → `ok`
- `--version` → `hello-docker <sha>` (used as the smoke test)

## CI

`.github/workflows/docker.yaml` calls the reusable workflow:

```yaml
jobs:
  docker:
    uses: leonidgrishenkov/github-actions/.github/workflows/docker.yaml@main
    with:
      image-name: hello-docker
      smoke-test-args: "--version"
    secrets: inherit
```

Pipeline: **lint → scan (amd64, Trivy, smoke test) → push (amd64+arm64, SBOM +
provenance, cosign sign)**.

## Prerequisites (in **this** repo's GitHub settings)

| Where | Name | Value |
|---|---|---|
| Settings → Secrets and variables → Actions → **Variables** | `YC_REGISTRY_ID` | your `crp…` registry id |
| Settings → Secrets and variables → Actions → **Secrets** | `YC_CR_SA_AUTH_JSON` | YC service account JSON key |

The reusable workflow itself lives in `leonidgrishenkov/github-actions` and
doesn't need these set on its own repo — they're read from the caller at runtime.

## Push it

```bash
git init -b main
git add -A
git commit -m "hello-docker test repo"
gh repo create leonidgrishenkov/test-github-repo --private --source=. --push
```

## Run locally

```bash
docker build -t hello-docker .
docker run --rm -p 8080:8080 hello-docker
curl localhost:8080/healthz   # ok
```
