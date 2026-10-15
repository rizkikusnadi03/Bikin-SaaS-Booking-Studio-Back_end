# 📸 Booking Studio Foto — Backend API

Backend service untuk platform **SaaS Booking Studio Foto** berbasis **Go (Golang)** dengan arsitektur **Modular Monolith**. Dibangun menggunakan framework Gin, PostgreSQL, dan SQLC.

---

## 📋 Daftar Isi

- [Gambaran Umum](#gambaran-umum)
- [Tech Stack](#tech-stack)
- [Struktur Proyek](#struktur-proyek)
- [Prasyarat](#prasyarat)
- [Cara Menjalankan](#cara-menjalankan)
- [Konfigurasi Environment](#konfigurasi-environment)
- [API Endpoints](#api-endpoints)
- [Arsitektur Modul](#arsitektur-modul)
- [Shared Packages](#shared-packages)
- [Middleware](#middleware)
- [Database](#database)
- [Panduan Pengembangan](#panduan-pengembangan)

---

## Gambaran Umum

Sistem ini merupakan backend REST API untuk platform booking studio foto yang memungkinkan:
- **Studio Owner** mendaftarkan studio, mengelola paket foto, portofolio, dan menerima booking.
- **Visitor / Customer** melihat studio publik dan melakukan booking secara online.

---

## Tech Stack

| Kategori | Teknologi | Versi |
|---|---|---|
| Bahasa | Go (Golang) | 1.27.1 |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) | v1.12.0 |
| Database | PostgreSQL | 16 (Alpine) |
| DB Driver | [pgx](https://github.com/jackc/pgx) | v5.11.0 |
| Query Generator | [SQLC](https://sqlc.dev/) | v2 |
| Validasi | [go-playground/validator](https://github.com/go-playground/validator) | v10.30.3 |
| CORS | [gin-contrib/cors](https://github.com/gin-contrib/cors) | v1.7.9 |
| Env Config | [godotenv](https://github.com/joho/godotenv) | v1.5.1 |
| Containerisasi | Docker + Docker Compose | - |

---

## Struktur Proyek

```
backend/
├── cmd/
│   └── api/
│       └── main.go            # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go          # Konfigurasi dari environment variables
│   ├── database/
│   │   ├── postgres.go        # Koneksi & connection pool PostgreSQL
│   │   └── sqlc/              # Auto-generated SQLC queries (jangan edit manual)
│   ├── middleware/
│   │   ├── cors.go            # Middleware CORS
│   │   ├── logger.go          # Middleware HTTP request logger
│   │   └── recovery.go        # Middleware panic recovery
│   ├── modules/               # Business logic (Modular Monolith)
│   │   ├── auth/              # Autentikasi & session
│   │   ├── studio/            # Manajemen profil studio
│   │   ├── package/           # Manajemen paket foto
│   │   ├── portfolio/         # Galeri portofolio studio
│   │   ├── booking/           # Pengelolaan booking (owner)
│   │   └── public/            # Endpoint publik untuk visitor
│   └── shared/
│       ├── errors/
│       │   └── errors.go      # Definisi error & mapping ke HTTP status
│       ├── response/
│       │   └── response.go    # Standar format JSON response
│       └── validator/
│           └── validator.go   # Format pesan error validasi
├── migrations/                # File migrasi database SQL
├── queries/                   # File query SQL untuk SQLC
├── docker-compose.yml         # Konfigurasi Docker untuk PostgreSQL
├── sqlc.yaml                  # Konfigurasi SQLC code generator
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
└── .env.example               # Template konfigurasi environment
```

---

## Prasyarat

Pastikan tools berikut sudah terinstal di mesin Anda:

- [Go](https://go.dev/dl/) >= 1.21
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/)
- [SQLC](https://docs.sqlc.dev/en/latest/overview/install.html) (opsional, untuk generate query)
- [golang-migrate](https://github.com/golang-migrate/migrate) (opsional, untuk menjalankan migrasi)

---

## Cara Menjalankan

### 1. Clone repository dan masuk ke folder backend

```bash
git clone <repo-url>
cd booking_studio_foto/backend
```

### 2. Salin file environment

```bash
cp .env.example .env
```

Sesuaikan nilai di `.env` dengan kebutuhan lokal Anda (lihat bagian [Konfigurasi Environment](#konfigurasi-environment)).

### 3. Jalankan database PostgreSQL via Docker

```bash
docker compose up -d
```

Perintah ini akan menjalankan container PostgreSQL di port **5433** (host) → **5432** (container).

> **Catatan:** Port host defaultnya `5433` (bukan `5432`) untuk menghindari konflik dengan PostgreSQL lokal.

### 4. Install dependencies Go

```bash
go mod tidy
```

### 5. Jalankan server

```bash
go run ./cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`.

### 6. Verifikasi server berjalan

```bash
curl http://localhost:8080/api/v1/health
```

Response yang diharapkan:

```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "database": "connected",
    "environment": "development",
    "status": "UP",
    "version": "1.0.0"
  }
}
```

---

## Konfigurasi Environment

Salin `.env.example` menjadi `.env` dan sesuaikan nilainya:

```env
# === Aplikasi ===
PORT=8080                   # Port server berjalan
ENV=development             # Mode: development | production

# === Database (PostgreSQL) ===
DB_HOST=localhost           # Host database
DB_PORT=5433                # Port database (5433 jika pakai Docker Compose default)
DB_USER=postgres            # Username database
DB_PASSWORD=postgres        # Password database
DB_NAME=photo_studio_db     # Nama database
DB_SSLMODE=disable          # SSL mode: disable | require | verify-full

# === Connection Pool ===
DB_MAX_CONNS=25             # Maksimum koneksi terbuka
DB_MIN_CONNS=5              # Maksimum koneksi idle
DB_MAX_IDLE_MINS=15         # Waktu idle maksimum (menit)
DB_MAX_LIFE_MINS=60         # Umur koneksi maksimum (menit)

# === CORS ===
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
# Pisahkan beberapa origin dengan koma
```

### Mode Production

Saat `ENV=production`, server otomatis:
- Mengaktifkan **Gin Release Mode** (log minimal, performa optimal)
- Menonaktifkan debug output

---

## API Endpoints

### Base URL

```
http://localhost:8080/api/v1
```

### Endpoints yang Tersedia

| Method | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| `GET` | `/api/v1/health` | Health check server & database | ❌ |

> 🚧 **Modul bisnis sedang dalam pengembangan.** Endpoint lengkap akan ditambahkan sesuai progress modul.

### Format Response

Semua response menggunakan format JSON yang konsisten:

#### ✅ Success Response

```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": { }
}
```

#### ❌ Error Response

```json
{
  "success": false,
  "message": "Pesan error",
  "errors": { }
}
```

---

## Arsitektur Modul

Backend ini menggunakan pola **Modular Monolith** dimana setiap fitur bisnis dikemas dalam modul terpisah di bawah `internal/modules/`.

```
internal/modules/
├── auth/         → Autentikasi & manajemen sesi pengguna
├── studio/       → Profil & pengaturan studio owner
├── package/      → Paket foto yang ditawarkan studio
├── portfolio/    → Galeri foto portofolio studio
├── booking/      → Pengelolaan booking dari sisi owner
└── public/       → API publik untuk portal booking visitor
```

### Konvensi File per Modul

Setiap modul mengikuti struktur file yang konsisten:

| File | Peran |
|------|-------|
| `dto.go` | Data Transfer Objects (request & response struct) |
| `handler.go` | HTTP handler / controller layer |
| `service.go` | Business logic layer |
| `route.go` | Definisi route & binding middleware |

### Alur Request

```
HTTP Request
    │
    ▼
Middleware (Logger → Recovery → CORS)
    │
    ▼
Router (Gin)
    │
    ▼
Handler (dto.go parsing & validasi)
    │
    ▼
Service (business logic)
    │
    ▼
Database (SQLC queries via sqlc/)
    │
    ▼
HTTP Response (format standar via response.go)
```

---

## Shared Packages

Package-package di `internal/shared/` dapat digunakan oleh semua modul.

### `internal/shared/response`

Menyediakan helper untuk mengirim JSON response yang konsisten.

```go
// Success response
response.Success(c, http.StatusOK, "Data berhasil diambil", data)

// Error response
response.Error(c, http.StatusBadRequest, "Input tidak valid", errors)
```

### `internal/shared/errors`

Mendefinisikan error standar aplikasi beserta mapping ke HTTP status code.

| Error Var | HTTP Status | Deskripsi |
|-----------|-------------|-----------|
| `ErrNotFound` | 404 | Resource tidak ditemukan |
| `ErrUnauthorized` | 401 | Akses tanpa autentikasi |
| `ErrForbidden` | 403 | Akses tidak diizinkan |
| `ErrBadRequest` | 400 | Request tidak valid |
| `ErrConflict` | 409 | Konflik data (duplikasi) |
| `ErrValidationFailed` | 400 | Validasi input gagal |
| `ErrInternalServer` | 500 | Error server internal |

```go
// Contoh penggunaan
import apperrors "booking-studio-foto/backend/internal/shared/errors"

statusCode := apperrors.MapErrorToStatusCode(err)
```

### `internal/shared/validator`

Memformat pesan error validasi dari `go-playground/validator` menjadi map yang mudah dibaca.

```go
import "booking-studio-foto/backend/internal/shared/validator"

if err := c.ShouldBindJSON(&req); err != nil {
    errs := validator.FormatValidationErrors(err)
    response.Error(c, http.StatusBadRequest, "Validasi gagal", errs)
    return
}
```

**Contoh output:**

```json
{
  "email": "email must be a valid email address",
  "password": "password must be at least 8 characters"
}
```

---

## Middleware

Semua middleware global didaftarkan di `cmd/api/main.go` dan otomatis berlaku untuk semua route.

### `LoggerMiddleware`

Mencatat setiap HTTP request dalam format:
```
[HTTP] 200 |      1.234ms |       127.0.0.1 | GET     /api/v1/health
```

### `RecoveryMiddleware`

Menangkap `panic` yang tidak tertangani dan mengembalikan response `500 Internal Server Error` secara aman, mencegah server crash.

### `CORSMiddleware`

Mengonfigurasi CORS berdasarkan `CORS_ALLOWED_ORIGINS` di `.env`.

| Konfigurasi | Nilai |
|-------------|-------|
| Allowed Methods | GET, POST, PUT, PATCH, DELETE, OPTIONS |
| Allowed Headers | Origin, Content-Type, Accept, Authorization, X-Requested-With |
| Allow Credentials | `true` |
| Max Age | 12 jam |

---

## Database

### Koneksi & Connection Pool

Koneksi database dikelola di `internal/database/postgres.go` menggunakan `pgx` sebagai driver. Server tetap bisa booting meski database tidak tersedia (menampilkan warning), sehingga health check endpoint tetap bisa merespons dengan status `database: disconnected`.

### SQLC

Query database di-generate secara otomatis menggunakan [SQLC](https://sqlc.dev/).

**Workflow SQLC:**

1. Tulis file migrasi SQL di `migrations/`
2. Tulis file query SQL di `queries/`
3. Jalankan SQLC untuk generate Go code:

```bash
sqlc generate
```

4. Hasil generate akan berada di `internal/database/sqlc/`

> ⚠️ **Jangan edit file di `internal/database/sqlc/` secara manual.** File tersebut di-generate otomatis dan akan ditimpa saat `sqlc generate` dijalankan.

### Konfigurasi SQLC (`sqlc.yaml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "migrations"   # Lokasi file migrasi
    queries: "queries"      # Lokasi file query SQL
    gen:
      go:
        package: "sqlc"
        out: "internal/database/sqlc"
        emit_json_tags: true
        emit_interface: true
```

---

## Panduan Pengembangan

### Menambah Modul Baru

1. Buat folder baru di `internal/modules/<nama-modul>/`
2. Buat file berikut sesuai konvensi:
   - `dto.go` — Struct request & response
   - `service.go` — Business logic
   - `handler.go` — HTTP handler
   - `route.go` — Registrasi route
3. Daftarkan route modul di fungsi `setupRoutes()` di `cmd/api/main.go`

### Menambah Migrasi Database

1. Buat file SQL baru di `migrations/` dengan format `{nomor}_{deskripsi}.up.sql`
2. Tambahkan query ke `queries/` jika perlu
3. Jalankan `sqlc generate` untuk update generated code
4. Jalankan migrasi menggunakan tool migrasi pilihan (contoh: `golang-migrate`)

### Build untuk Production

```bash
go build -o bin/api ./cmd/api/main.go
```

### Menjalankan dengan environment production

```bash
ENV=production ./bin/api
```

---

## Graceful Shutdown

Server mendukung **graceful shutdown** dengan menunggu request yang sedang berjalan selesai (timeout 5 detik) sebelum mematikan server. Kirim sinyal `SIGINT` (Ctrl+C) atau `SIGTERM` untuk memicu shutdown.

```
[INFO] Shutting down server gracefully...
[INFO] Server exited successfully
```

---

> 📄 Dokumentasi ini dibuat untuk tim pengembang. Untuk pertanyaan lebih lanjut, hubungi maintainer proyek.
