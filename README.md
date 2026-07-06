# SWIFT MX ISO 20022 Message Builder (Gin + Go)

Project builder untuk generate pesan **SWIFT MX (ISO 20022)** — `pacs.008`,
`pacs.009`, `pacs.002`, dan `pacs.004` — menggunakan **Gin** (Go web
framework). Setiap pesan yang berhasil dibangun akan di-*marshal* menjadi
XML format MX, lalu ditulis ke folder output secara **asynchronous**
menggunakan **worker pool (goroutine)**, sehingga endpoint HTTP tidak
perlu menunggu proses I/O selesai.

## Fitur

- 4 service pesan ISO 20022:
  - **pacs.008.001.08** — FIToFICustomerCreditTransfer
  - **pacs.009.001.08** — FinancialInstitutionCreditTransfer
  - **pacs.002.001.10** — FIToFIPaymentStatusReport
  - **pacs.004.001.09** — PaymentReturn
- Setiap service punya:
  - `POST /generate` → membangun dokumen MX dan mengirim job ke worker pool
  - `GET /inquiry/:messageId` → **inquiry status message** untuk mengecek
    progres penulisan file (`PENDING` → `PROCESSING` → `COMPLETED`/`FAILED`)
- Worker pool goroutine (jumlah worker & ukuran antrian bisa dikonfigurasi)
  menulis file `.xml` ke folder `output/{pacs008|pacs009|pacs002|pacs004}/`
- Graceful shutdown (menunggu worker menyelesaikan job yang sedang berjalan)

## Struktur Project

```
swift-mx-builder/
├── main.go                     # entrypoint, routing Gin
├── config/
│   └── config.go               # konfigurasi (port, output dir, worker count)
├── models/
│   ├── common.go                # tipe bersama (GroupHeader, Party, dst.)
│   ├── pacs008.go               # struct XML pacs.008.001.08
│   ├── pacs009.go               # struct XML pacs.009.001.08
│   ├── pacs002.go               # struct XML pacs.002.001.10
│   └── pacs004.go               # struct XML pacs.004.001.09
├── handlers/
│   ├── pacs008_handler.go       # POST /pacs008/generate
│   ├── pacs009_handler.go       # POST /pacs009/generate
│   ├── pacs002_handler.go       # POST /pacs002/generate
│   ├── pacs004_handler.go       # POST /pacs004/generate
│   └── inquiry_handler.go       # GET /{service}/inquiry/:messageId (generik, dipakai 4 service)
├── worker/
│   └── worker.go                 # worker pool goroutine + status tracking
├── utils/
│   ├── xmlbuilder.go             # marshal struct → MX XML
│   └── idgen.go                  # generator MsgId, TxId, UETR
└── output/                       # folder tujuan file MX hasil generate
```

## Menjalankan

```bash
go mod tidy      # unduh dependency gin (butuh akses internet penuh)
go run main.go
```

Server berjalan di `http://localhost:8080` (ubah lewat env `APP_PORT`).
Konfigurasi lain: `MX_OUTPUT_DIR` (default `./output`).

## Contoh Pemakaian API

### 1. Generate pacs.008 (Customer Credit Transfer)

```bash
curl -X POST http://localhost:8080/api/v1/pacs008/generate \
  -H "Content-Type: application/json" \
  -d '{
    "end_to_end_id": "E2E-INV-0001",
    "amount": 15000000,
    "currency": "IDR",
    "debtor_name": "PT Sumber Makmur",
    "debtor_agent_bic": "CENAIDJAXXX",
    "creditor_name": "PT Cahaya Abadi",
    "creditor_agent_bic": "BMRIIDJAXXX",
    "remittance_info": "Payment for invoice INV-0001",
    "settlement_method": "CLRG"
  }'
```

Response (202 Accepted):
```json
{
  "message_id": "PACS008-20260706153012-A1B2C3",
  "transaction_id": "TX-20260706153012-D4E5F6",
  "message_type": "pacs.008.001.08",
  "status": "PENDING",
  "file_name": "PACS008-20260706153012-A1B2C3.xml",
  "inquiry_url": "/api/v1/pacs008/inquiry/PACS008-20260706153012-A1B2C3"
}
```

### 2. Cek status (inquiry)

```bash
curl http://localhost:8080/api/v1/pacs008/inquiry/PACS008-20260706153012-A1B2C3
```

```json
{
  "message_id": "PACS008-20260706153012-A1B2C3",
  "message_type": "pacs008",
  "status": "COMPLETED",
  "file_path": "output/pacs008/PACS008-20260706153012-A1B2C3.xml",
  "submitted_at": "2026-07-06T15:30:12+07:00",
  "updated_at": "2026-07-06T15:30:12+07:00"
}
```

### 3. Generate pacs.009 (Financial Institution Credit Transfer)

```bash
curl -X POST http://localhost:8080/api/v1/pacs009/generate \
  -H "Content-Type: application/json" \
  -d '{
    "end_to_end_id": "E2E-COVER-0001",
    "amount": 500000,
    "currency": "USD",
    "debtor_bic": "CENAIDJAXXX",
    "debtor_agent_bic": "CHASUS33XXX",
    "creditor_agent_bic": "DEUTDEFFXXX",
    "creditor_bic": "BMRIIDJAXXX"
  }'
```

### 4. Generate pacs.002 (Payment Status Report)

```bash
curl -X POST http://localhost:8080/api/v1/pacs002/generate \
  -H "Content-Type: application/json" \
  -d '{
    "original_msg_id": "PACS008-20260706153012-A1B2C3",
    "original_msg_name_id": "pacs.008.001.08",
    "original_end_to_end_id": "E2E-INV-0001",
    "tx_status": "ACSC"
  }'
```

### 5. Generate pacs.004 (Payment Return)

```bash
curl -X POST http://localhost:8080/api/v1/pacs004/generate \
  -H "Content-Type: application/json" \
  -d '{
    "original_end_to_end_id": "E2E-INV-0001",
    "returned_amount": 15000000,
    "currency": "IDR",
    "reason_code": "AC04",
    "additional_info": "Account closed"
  }'
```

## Catatan Implementasi

- Struct XML pada folder `models/` adalah versi **disederhanakan** dari
  skema resmi ISO 20022 (elemen inti sudah mengikuti nama tag official:
  `GrpHdr`, `CdtTrfTxInf`, `PmtId`, `IntrBkSttlmAmt`, `TxSts`, `RtrRsnInf`,
  dsb). Untuk kebutuhan produksi/kepatuhan penuh terhadap skema SWIFT MX,
  tambahkan elemen wajib lain sesuai spesifikasi resmi
  (`Ccy`, `PstlAdr`, `RgltryRptg`, dll.) pada masing-masing struct.
- Worker pool berjalan dalam proses yang sama (in-memory status map).
  Untuk skala produksi/multi-instance, ganti status tracking dengan
  Redis/database agar status inquiry konsisten lintas instance/pod.
- Jalankan `go mod tidy` di lingkungan dengan akses internet penuh
  (sandbox pengujian ini memblokir sebagian domain dependency Go seperti
  `golang.org`, `gopkg.in`, `rsc.io` sehingga `go get` tidak bisa
  menyelesaikan seluruh graph modul secara otomatis). Semua source code
  telah divalidasi sintaksisnya dengan `gofmt`.
