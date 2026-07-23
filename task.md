# Task Breakdown: Test Backend - Sharing Vision 2023

## 📌 Deskripsi Singkat
Implementasi backend microservice (Golang) untuk use case **Post Article** dengan database MySQL, fitur migration, Endpoint CRUD dengan validasi input, serta Postman Collection.

---

## 🗄️ 1. Database & Migration (Bobot 20%)
- [x] **Setup Database MySQL**
  - [x] Buat database `article` (diatur melalui `database/schema.sql` dan `config/database.go`).
- [x] **Skema Tabel `posts`**
  - [x] Kolom `Id`: `INT`, Auto Increment, Primary Key
  - [x] Kolom `Title`: `VARCHAR(200)`
  - [x] Kolom `Content`: `TEXT`
  - [x] Kolom `Category`: `VARCHAR(100)`
  - [x] Kolom `Created_date`: `TIMESTAMP`
  - [x] Kolom `Updated_date`: `TIMESTAMP`
  - [x] Kolom `Status`: `VARCHAR(100)` (Pilihan: `publish`, `draft`, `thrash`)
- [x] **Database Migration (Golang)**
  - [x] GORM AutoMigrate di `main.go` (dapat dijalankan via `go run main.go migrate` atau otomatis saat server start).

---

## 🚀 2. Microservice & REST API (Bobot 80%)

### 🛡️ Validasi Request Data (JSON)
- [x] **Rule Validasi Title**: `required`, minimal 20 karakter.
- [x] **Rule Validasi Content**: `required`, minimal 200 karakter.
- [x] **Rule Validasi Category**: `required`, minimal 3 karakter.
- [x] **Rule Validasi Status**: `required`, hanya boleh bernilai `"publish"`, `"draft"`, atau `"thrash"`.

### 🌐 Endpoints Implementation
- [x] **1. Create Article (`POST /article/`)**
  - [x] Validasi JSON request payload.
  - [x] Simpan article baru ke database.
  - [x] Response `{}` (200 OK).
- [x] **2. Get All Articles with Pagination (`GET /article/<limit>/<offset>`)**
  - [x] Ambil data daftar article dari database berdasarkan `limit` dan `offset`.
  - [x] Response format: Array of JSON objects `[{ "title": "...", "content": "...", "category": "...", "status": "..." }]`.
- [x] **3. Get Article by ID (`GET /article/<id>`)**
  - [x] Ambil detail article berdasarkan `id`.
  - [x] Response format: Single JSON object `{ "title": "...", "content": "...", "category": "...", "status": "..." }`.
- [x] **4. Update Article by ID (`POST` / `PUT` / `PATCH /article/<id>`)**
  - [x] Validasi JSON request payload.
  - [x] Update data article di database berdasarkan `id`.
  - [x] Response `{}`.
- [x] **5. Delete Article by ID (`POST` / `DELETE /article/<id>`)**
  - [x] Hapus data article dari database berdasarkan `id`.
  - [x] Response `{}`.

---

## 📄 3. Postman Collection (Bobot 15%)
- [x] File Postman Collection JSON (`postman/Sharing_Vision_Articles.postman_collection.json`).
- [x] Berisi 6 item request mencakup pengujian sukses dan error validasi.

---

## 🧪 4. Testing & Verifikasi
- [x] Skema database & migration terverifikasi.
- [x] Unit test validasi di `services/post_service_test.go` lulus 100%.
- [x] Build verifikasi `go build` tanpa error.
