# Test Backend Sharing Vision 2023

Aplikasi REST API Backend berbasis Golang dan MySQL untuk use case Post Article.

## Persyaratan System
- Go (v1.26+)
- MySQL Database Server (port 3306)

## Langkah Penggunaan

### 1. Import Database SQL
Import berkas `article.sql` ke MySQL Database Anda untuk membuat database `article` dan tabel `posts`.

### 2. Konfigurasi Environment (`.env`)
Salin berkas `.env.example` menjadi `.env` dan sesuaikan kredensial MySQL lokal Anda:
```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=article
SERVER_PORT=8080
```

### 3. Menjalankan Server
Jalankan perintah berikut pada direktori backend:
```bash
go run main.go
```
Server REST API akan aktif di `http://localhost:8080`.

## Daftar Endpoint REST API
- `POST /article/` - Membuat artikel baru (Request Body: `title`, `content`, `category`, `status`)
- `GET /article/<limit>/<offset>` - Menampilkan daftar artikel dengan pagination
- `GET /article/<id>` - Menampilkan detail artikel berdasarkan ID
- `PUT /article/<id>` - Memperbarui data artikel berdasarkan ID
- `DELETE /article/<id>` - Menghapus artikel berdasarkan ID

## Pengujian via Postman
Berkas `postman_collection.json` dapat di-import ke aplikasi Postman untuk menguji seluruh endpoint di atas.
