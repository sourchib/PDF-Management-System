# Dokumentasi API - PDF Management System

Dokumentasi ini menjelaskan seluruh endpoint yang tersedia pada sistem manajemen PDF, termasuk autentikasi dan operasi file.

## Informasi Umum
- **Base URL**: `http://localhost:8080`
- **Format Response**: JSON
- **Otentikasi**: JWT (JSON Web Token) Bearer-Token

---

## 1. Authentication

### Register User
Digunakan untuk membuat akun pengguna baru.
- **Endpoint**: `POST /api/auth/register`
- **Method**: `POST`
- **Body Request**:
```json
{
  "name": "Nama Lengkap",
  "email": "user@example.com",
  "password": "password_aman",
  "address": "Alamat Tinggal",
  "phone_number": "08123456789",
  "post_code": "12345",
  "role_id": 1
}
```
> **Note**: `role_id` tersedia: 1 (Project Manager), 2 (Financial), 3 (HRD).

- **Response Success (200 OK)**:
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": { "id": 1, "name": "Nama Lengkap", ... }
}
```

### Login
Digunakan untuk mendapatkan Token akses.
- **Endpoint**: `POST /api/auth/login`
- **Method**: `POST`
- **Body Request**:
```json
{
  "email": "user@example.com",
  "password": "password_aman"
}
```
- **Response Success (200 OK)**:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1Ni...",
    "user": { "id": 1, "email": "user@example.com", ... }
  }
}
```

---

## 2. PDF Management
*Seluruh endpoint di bawah ini membutuhkan Header:*
`Authorization: Bearer <JWT_TOKEN>`

### Generate Report PDF
Membuat file PDF secara otomatis berdasarkan parameter.
- **Endpoint**: `/api/pdf/generate`
- **Method**: `POST`
- **Body Request**:
```json
{
    "title": "Laporan Keuangan Kuartal I 2026",
    "institution_name": "Arema FC Finance Department",
    "address": "Jl. Kertanegara No. 7, Malang, Jawa Timur",
    "phone": "(0341) 333-1987",
    "logo_url": "https://i.ibb.co.com/3Yf2yg0t/Arema-FC-2017-logo.png",
        "content": "Laporan ini merangkum kinerja keuangan klub Arema FC pada kuartal pertama tahun 2026. Pendapatan utama berasal dari penjualan tiket pertandingan kandang dan kerjasama sponsor baru. Total pendapatan tercatat meningkat sebesar 20% dibandingkan periode yang sama tahun lalu. Beban operasional terkendali dengan fokus pada optimalisasi biaya akademi pemain muda."
    }
```
- **Response Success (200 OK)**:
```json
{
  "success": true,
  "message": "PDF generated successfully",
  "data": {
    "id": 1,
    "filename": "report_20260128_abc123.pdf",
    "filepath": "/uploads/pdf/report_20260128_abc123.pdf",
    "status": "CREATED",
    "created_at": "2026-01-28T12:00:00Z"
  }
}
```

### Upload PDF
Mengupload file PDF yang sudah ada.
- **Endpoint**: `/api/pdf/upload`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
- **Body**:
  - `file`: (Binary File) Hanya file `.pdf` maks 10MB.
- **Response Success (200 OK)**:
```json
{
  "success": true,
  "message": "PDF uploaded successfully",
  "data": {
    "id": 2,
    "filename": "upload_20260128_xyz789.pdf",
    "original_name": "dokumen.pdf",
    "status": "UPLOADED",
    "size": 1024567
  }
}
```

### List PDF Files
Menampilkan daftar semua file PDF.
- **Endpoint**: `/api/pdf/list`
- **Method**: `GET`
- **Query Parameters**:
  - `status`: Filter status (CREATED, UPLOADED, DELETED)
  - `page`: Nomor halaman (default: 1)
  - `limit`: Data per halaman (default: 10)
- **Response Success (200 OK)**:
```json
{
  "success": true,
  "data": [
    { "id": 1, "filename": "...", "status": "CREATED", ... }
  ],
  "pagination": { "page": 1, "limit": 10, "total": 1 }
}
```

### Delete PDF (Soft Delete)
Menghapus file dari daftar tanpa menghapus file fisiknya.
- **Endpoint**: `/api/pdf/{id}`
- **Method**: `DELETE`
- **Response Success (200 OK)**:
```json
{
  "success": true,
  "message": "PDF deleted successfully",
  "data": {
    "id": 1,
    "status": "DELETED",
    "deleted_at": "2026-01-28T12:30:00Z"
  }
}
```

---

## 3. Error Codes & Messages

| Status Code | Message | Deskripsi |
|---|---|---|
| 400 | Invalid request body | Payload JSON tidak sesuai format |
| 401 | Missing/Invalid Token | Tidak ada atau token JWT salah |
| 403 | Forbidden | Akses ditolak |
| 404 | File not found | ID PDF yang dicari tidak ditemukan |
| 500 | Internal Server Error | Kesalahan pada server |

---
**Author**: Muchammad Muchib Zainul Fikry
**Project**: PDF Management System Technical Test
