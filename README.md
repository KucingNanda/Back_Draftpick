# Draft Pick Backend

Aplikasi *backend* super cepat berbasis **Golang** (Framework Fiber) untuk menyuplai *API endpoint* ke aplikasi simulasi *Draft Pick*. Dibekali dengan sistem manajemen database ORM otomatis via **GORM** dan **PostgreSQL**.

## 🚀 Fitur Utama
- **RESTful API Cepat**: Mendukung penuh operasi CRUD untuk entitas Karakter, Peserta, Riwayat, BGM, dan Settings.
- **Keamanan Lapis Ganda**: 
  - Validasi *JSON Web Token* (JWT)
  - Pembatasan Hak Akses (*Role*: `wasit` & `peserta`)
  - Hashing Sandi via Bcrypt
- **Dokumentasi API Otomatis**: Dilengkapi dengan Swagger UI bawaan yang interaktif.
- **Cloud Storage**: Terhubung ke sistem *Cloudinary* untuk menampung gambar profil dan file mp3 BGM secara eksternal.

## 🛠️ Prasyarat
- Golang (v1.20 atau versi lebih baru)
- Database PostgreSQL (Contoh: Alwaysdata)
- Kredensial akun Cloudinary

## ⚙️ Instalasi & Menjalankan di Lokal

1. Buat file bernama `.env` berdasarkan salinan file `.env.example` dan lengkapi rinciannya:
   ```env
   DB_HOST=postgresql-kucing27.alwaysdata.net
   DB_PORT=5432
   DB_USER=yourdbuser
   DB_PASSWORD=yoursecret
   DB_NAME=yourdbname
   
   CLOUDINARY_URL=cloudinary://<API_KEY>:<API_SECRET>@<CLOUD_NAME>
   
   JWT_SECRET=super_rahasia_123
   SWAGGER_HOST=127.0.0.1:3000
   ```
2. Unduh dan perbarui semua modul:
   ```bash
   go mod tidy
   ```
3. *(Opsional)* Jika Anda baru saja memodifikasi komentar Swagger di *handler*, jangan lupa *generate* ulang dokumen:
   ```bash
   swag init
   ```
4. Jalankan *server* lokal:
   ```bash
   go run main.go
   ```
5. Selamat! API kini berjalan. Buka dokumentasi Swagger di: `http://localhost:3000/docs/index.html`.

## 📦 Panduan Deployment (Render)
Ketika melakukan *deploy* ke platform cloud (Render / Koyeb):
1. **Penting:** Anda tidak mengunggah file `.env`.
2. Anda harus menyalin isi variabel di atas satu per satu ke menu **Environment Variables** di *dashboard* platform cloud tujuan Anda.
3. Gunakan *Build Command*: `go build main.go` dan *Start Command*: `./main`.
