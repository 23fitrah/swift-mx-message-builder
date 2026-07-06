# SWIFT MX ISO 20022 Message Builder (Gin + Go)

A Go-based REST API for generating SWIFT MX (ISO 20022) messages.

This service accepts JSON requests and builds ISO 20022 compliant XML messages. It is designed to simplify the process of generating SWIFT MX messages and to provide a reusable component for banking and payment integrations.

The project follows a modular structure, making it easier to support additional message types as business requirements grow.

## Feature

- 4 services message ISO 20022:
  - **pacs.008.001.08** — FI To FI Customer Credit Transfer
  - **pacs.009.001.08** — Financial Institution Credit Transfer
  - **pacs.002.001.10** — FI To FI Payment Status Reject Or Accept
  - **pacs.004.001.09** — Payment Return
- Each service has:
  - `POST /generate` → builds an MX document and sends the job to the worker pool
  - `GET /inquiry/:messageId` → **inquiry status message** to check
file writing progress (`PENDING` → `PROCESSING` → `COMPLETED`/`FAILED`)
- Worker pool go routine (number of workers and queue size are configurable)
  - writes an `.xml` file to the `output/{pacs008|pacs009|pacs002|pacs004}/` folder
  - Graceful shutdown (waits for the worker to complete the current job)

## Tech Stack

* Go (Golang)
* Gin Framework
* Docker Container
* Unit Test
* CI/CD Jenkins
* Rabbit MQ

## API List

| Method | Endpoint | Description |
|---|---|---|
| POST | /api/v1/pacs008/generate | FI To FI Customer Credit Transfer |
| GET | /api/v1/pacs008/inquiry/:messageId | Inquiry Status pacs008 |
| POST | /api/v1/pacs002/generate | FI To FI Payment Status Reject Or Accept |
| GET | /api/v1/pacs002/inquiry/:messageId | Inquiry Status pacs002 |
| POST | /api/v1/pacs004/generate | Payment Refund |
| GET | /api/v1/pacs004/inquiry/:messageId  | Inquiry Status pacs004 |
| POST | /api/v1/pacs009/generate | Financial Institution Credit Transfer |
| GET | /api/v1/pacs009/inquiry/:messageId | Inquiry Status pacs009|

##  Project Structure

```
swift-mx-builder/
├── main.go                      # entrypoint, routing Gin
├── config/
│   └── config.go                # configuration (port, output dir, worker count)
├── dto/                         # request input
├── models/
│   ├── common.go                # shared type (GroupHeader, Party, dst.)
│   ├── pacs008.go               # struct XML pacs.008.001.08
│   ├── pacs009.go               # struct XML pacs.009.001.08
│   ├── pacs002.go               # struct XML pacs.002.001.10
│   └── pacs004.go               # struct XML pacs.004.001.09
├── handlers/
│   ├── pacs008_handler.go       # POST /pacs008/generate
│   ├── pacs009_handler.go       # POST /pacs009/generate
│   ├── pacs002_handler.go       # POST /pacs002/generate
│   ├── pacs004_handler.go       # POST /pacs004/generate
│   └── inquiry_handler.go       # GET /{service}/inquiry/:messageId (generic, used 4 services)
├── worker/
│   └── worker.go                 # worker pool goroutine + status tracking
├── utils/
│   ├── xmlbuilder.go             # marshal struct → MX XML
│   ├── auth.go                   # Authorization middleware
│   ├── validation.go             # validation input
│   └── idgen.go                  # generator MsgId, TxId, UETR
└── output/                       # folder for result generate file mx
```

## Running

```bash
go mod tidy      # unduh dependency gin (butuh akses internet penuh)
go run main.go
```

Server running at `http://localhost:9090` (change at .env `APP_PORT`).
other config: `MX_OUTPUT_DIR` (default `./output`).

## Example using API

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

Response (200 Accepted):
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

### 2. Check status (inquiry)

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

## Notes

- The XML models in the models/ directory are based on a simplified version of the ISO 20022 schema. Core elements such as GrpHdr, CdtTrfTxInf, PmtId, IntrBkSttlmAmt, TxSts, and RtrRsnInf follow the official SWIFT MX element names. If full ISO 20022 compliance is required, additional mandatory elements (such as Ccy, PstlAdr, RgltryRptg, and others) should be added according to the relevant message specification.
- The worker pool currently runs within a single process and keeps job status in memory. This works well for development and single-instance deployments. For production environments with multiple instances or pods, use a shared storage such as Redis or a database so job status can be accessed consistently across all instances.
- Run go mod tidy in an environment with internet access before building the project. Some Go module hosts may not be reachable in restricted or sandboxed environments, which can prevent dependencies from being downloaded successfully.

## License

This project is intended for portfolio demonstration.

