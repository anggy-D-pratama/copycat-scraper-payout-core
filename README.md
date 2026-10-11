# Copycat Scraper Payout Core

Mock core engine untuk aplikasi scraper payout berbasis Go (Fiber v2).

---

## CI/CD Deployment ke Google Cloud Run (via Workload Identity Federation)

Pipeline deployment otomatis menggunakan GitHub Actions (`.github/workflows/deploy.yml`) saat commit di-push ke branch `main` atau melalui trigger manual (`workflow_dispatch`).

### 1. Prasyarat GCP Resource & IAM

Service Account yang digunakan oleh GitHub Actions (`WIF_SERVICE_ACCOUNT`) memerlukan peran minimal di GCP:

- `roles/run.admin` (Cloud Run Admin)
- `roles/artifactregistry.writer` (Artifact Registry Writer)
- `roles/iam.serviceAccountUser` (Service Account User)

### 2. GitHub Actions Secrets / Variables

Konfigurasikan Secrets berikut di repositori GitHub (**Settings > Secrets and variables > Actions > Secrets**):

| Secret Name | Deskripsi | Contoh Nilai |
|-------------|-----------|--------------|
| `GCP_PROJECT_ID` | Google Cloud Project ID | `my-gcp-project-123456` |
| `WIF_PROVIDER` | Full resource name Workload Identity Provider | `projects/123456789/locations/global/workloadIdentityPools/github-pool/providers/github-provider` |
| `WIF_SERVICE_ACCOUNT` | Email Service Account GCP dengan IAM role terkait | `github-actions-sa@my-gcp-project-123456.iam.gserviceaccount.com` |
| `GAR_LOCATION` *(Opsional)* | Lokasi region Artifact Registry (default: `asia-southeast2`) | `asia-southeast2` |
| `GAR_REPOSITORY` *(Opsional)* | Nama repository Artifact Registry (default: `copycat-payout-core`) | `copycat-payout-core` |
| `CLOUD_RUN_SERVICE` *(Opsional)* | Nama service Google Cloud Run (default: `copycat-payout-core`) | `copycat-payout-core` |
| `CLOUD_RUN_REGION` *(Opsional)* | Region Google Cloud Run (default: `asia-southeast2`) | `asia-southeast2` |

---

## Menjalankan Secara Lokal

### Menggunakan Go langsung

```bash
export PORT=8080
go run main.go
```

### Menggunakan Docker

```bash
docker build -t copycat-payout-core:latest .
docker run -p 8080:8080 -e PORT=8080 copycat-payout-core:latest
```
