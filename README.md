# Backend ILITS - REST API Pencatatan Akta Kelahiran

REST API sederhana untuk mencatat data akta kelahiran (CRUD) menggunakan **Go**, **Gin**, dan **GORM**.
Database memakai SQLite supaya project bisa langsung dijalankan tanpa perlu setup server database.

## Tech Stack

- Go 1.25+
- [Gin](https://github.com/gin-gonic/gin) sebagai HTTP framework
- [GORM](https://gorm.io) sebagai ORM
- SQLite (driver `github.com/glebarez/sqlite`, tidak butuh CGO)

## Struktur Folder

```
backend-ilits/
├── main.go                 # entry point, menjalankan server
├── config/
│   └── database.go         # koneksi database + auto migrate
├── models/
│   └── akta.go             # struct AktaKelahiran dan AktaInput
├── handlers/
│   └── akta_handler.go     # logic tiap endpoint
└── routes/
    └── routes.go           # daftar route
```

Alurnya: request masuk ke `routes`, diteruskan ke `handlers`, lalu handler membaca/menulis
data lewat GORM ke tabel yang strukturnya didefinisikan di `models`.

## Cara Menjalankan

```bash
git clone <link-repo-ini>
cd backend-ilits
go mod download
go run main.go
```

Server berjalan di `http://localhost:8080`. Port bisa diganti lewat environment variable `PORT`.
File database `akta.db` otomatis dibuat saat pertama kali dijalankan.

## Data Akta Kelahiran

| Field           | Keterangan                                           |
| --------------- | ---------------------------------------------------- |
| `nomor_akta`    | dibuat otomatis, format `AKL-<ddmmyyyy>-<id>`        |
| `nama_anak`     | wajib                                                |
| `jenis_kelamin` | wajib, `L` atau `P`                                  |
| `tempat_lahir`  | wajib                                                |
| `tanggal_lahir` | wajib, format `YYYY-MM-DD`, tidak boleh di masa depan |
| `nama_ayah`     | opsional                                             |
| `nama_ibu`      | wajib                                                |

## Endpoint

| Method | Endpoint         | Keterangan                                            |
| ------ | ---------------- | ----------------------------------------------------- |
| GET    | `/ping`          | cek server                                            |
| GET    | `/api/akta`      | ambil semua akta (filter `?nama=` dan `?tahun=2024`)  |
| GET    | `/api/akta/:id`  | ambil satu akta                                       |
| POST   | `/api/akta`      | catat akta baru                                       |
| PUT    | `/api/akta/:id`  | ubah data akta                                        |
| DELETE | `/api/akta/:id`  | hapus akta                                            |

Contoh body untuk POST dan PUT:

```json
{
  "nama_anak": "Dhaniel Dhaneswara",
  "jenis_kelamin": "L",
  "tempat_lahir": "Surabaya",
  "tanggal_lahir": "2007-03-15",
  "nama_ayah": "Pedrosa",
  "nama_ibu": "Aminah"
}
```

Kalau input tidak valid, API mengembalikan status `400`.
Kalau id akta tidak ada, API mengembalikan status `404`.

### Contoh request

```bash
# catat akta baru
curl -X POST http://localhost:8080/api/akta \
  -H "Content-Type: application/json" \
  -d '{"nama_anak":"Dhaniel Dhaneswar","jenis_kelamin":"L","tempat_lahir":"Surabaya","tanggal_lahir":"2007-03-15","nama_ayah":"Pedrosa","nama_ibu":"Aminah"}'

# lihat semua akta
curl http://localhost:8080/api/akta

# cari akta yang lahir tahun 2024
curl "http://localhost:8080/api/akta?tahun=2024"

# hapus akta id 1
curl -X DELETE http://localhost:8080/api/akta/1
```

---

## Dasar Version Control System (Git)

### Apa itu VCS dan Git?

Version Control System (VCS) adalah sistem yang mencatat setiap perubahan pada file dalam sebuah
project. Dengan VCS kita bisa melihat riwayat perubahan, kembali ke versi sebelumnya kalau ada
yang rusak, dan bekerja bareng orang lain tanpa saling menimpa pekerjaan.

Git adalah VCS yang paling banyak dipakai. Git bersifat *distributed*, artinya setiap orang
menyimpan salinan lengkap repository beserta riwayatnya di komputer masing-masing. GitHub adalah
layanan hosting untuk repository Git, jadi tempat menyimpan repository secara online.

### Istilah penting

- **Repository (repo)**: folder project yang perubahannya dilacak oleh Git.
- **Working directory**: file yang sedang kita edit.
- **Staging area**: tempat menampung perubahan yang akan dimasukkan ke commit berikutnya.
- **Commit**: snapshot perubahan yang disimpan beserta pesan, penulis, dan waktunya.
- **Branch**: jalur pengembangan terpisah, dipakai untuk mengerjakan fitur tanpa mengganggu `main`.
- **Merge**: menggabungkan perubahan dari satu branch ke branch lain.
- **Remote**: repository di server (misalnya GitHub), biasanya bernama `origin`.
- **Pull Request**: permintaan untuk menggabungkan branch ke `main`, biasanya direview dulu oleh tim.
- **Merge conflict**: terjadi saat dua branch mengubah baris yang sama, harus diselesaikan manual.

### Perintah yang sering dipakai

```bash
git init                      # membuat repository baru
git clone <url>               # menyalin repository dari remote
git status                    # melihat file yang berubah
git add <file>                # memasukkan perubahan ke staging (git add . untuk semua)
git commit -m "pesan"         # menyimpan perubahan yang sudah di-stage
git log --oneline             # melihat riwayat commit
git diff                      # melihat perubahan yang belum di-stage

git branch                    # melihat daftar branch
git checkout -b fitur-baru    # membuat branch baru lalu pindah ke sana
git switch main               # pindah ke branch main
git merge fitur-baru          # menggabungkan branch fitur-baru ke branch aktif

git remote add origin <url>   # menghubungkan repo lokal dengan GitHub
git push -u origin main       # mengirim commit ke remote
git pull                      # mengambil perubahan terbaru dari remote
```

### Alur kerja yang dipakai di project ini

1. Buat repository dengan `git init`, lalu hubungkan ke GitHub dengan `git remote add origin`.
2. Kerjakan satu bagian (misalnya model, handler, atau route), cek dengan `git status` dan `git diff`.
3. `git add` lalu `git commit` dengan pesan yang menjelaskan apa yang diubah.
4. `git push` ke GitHub.

Kalau bekerja dalam tim, fitur baru sebaiknya dikerjakan di branch sendiri
(`git checkout -b feat/nama-fitur`), lalu digabung ke `main` lewat Pull Request setelah direview.

### Tips pesan commit

Pesan commit sebaiknya singkat dan menjelaskan perubahan, misalnya:

- `feat: tambah endpoint update akta`
- `fix: validasi tanggal lahir di masa depan`
- `docs: tambah cara menjalankan di README`

`.gitignore` dipakai untuk mengecualikan file yang tidak perlu masuk repository, seperti file
database lokal (`*.db`) dan file konfigurasi rahasia.
