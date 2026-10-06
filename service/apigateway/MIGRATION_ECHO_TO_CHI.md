# Migrasi Apigateway: Echo → Chi

> **Lingkup:** `monolith-ecommerce-grpc/service/apigateway` (+ harness `tests/` dan `pkg/upload_image`)
> **Router:** `github.com/labstack/echo/v4` → `github.com/go-chi/chi/v5`
> **Status:** ✅ Selesai — build ✅ · vet ✅ · unit test ✅
> **Pola:** sama dengan migrasi apigateway `monolith-payment-gateway-grpc` (lihat `MIGRATION_ECHO_TO_CHI.md` di sana).

---

## 1. Ringkasan

| Item | Jumlah |
|---|---|
| File yang import echo (sebelum) | 76 |
| Paket handler domain yang dimigrasi | 21 (`auth`, `banner`, `cart`, `category`, `merchant`, `merchant_award`, `merchant_business`, `merchant_detail`, `merchant_document`, `merchant_policy`, `merchant_social_link`, `order`, `order_item`, `product`, `review`, `review_detail`, `role`, `shipping_address`, `slider`, `transaction`, `user`) |
| Rute terdaftar (terkonversi) | 251 (120 GET · 111 POST · 20 DELETE) |
| Middleware global | 11 (lebih banyak dari payment-gateway) |
| Test harness `tests/` yang dikonversi | ±39 file |
| Rute terdampak duplikasi mount | 1 (`/api/order-item`) — dimerge |

## 2. Perbedaan dari Payment-Gateway (yang perlu disesuaikan)

Pola dasar identik, tapi ada 6 varian yang ditangani:

1. **Error handler global yang lebih pintar** — ecommerce punya `RegisterErrorHandler` yang memetakan `*AppError` shared → status code asli (bukan 500) dengan `trace_id` dari header `X-Request-ID`, dan fallback ke format default echo (`{"message": "..."}`). Di chi, logika ini pindah ke adapter `httpx.Handler` (untuk rute unwrapped) dan `apierror.HandleApiError` (untuk rute wrapped, plus dukungan `*echo.HTTPError` → `httpx.HTTPError`). `middlewares/errorhandler.go` dihapus.
2. **Middleware ekstra** — `trace.go` (OTel span per request) dan `metrics.go` (HTTP metrics) di-rewrite ke chi dengan `wrapResponseWriter` untuk menangkap status code (pengganti `c.Response().Status`); `rate limit` khusus prefix `/api/auth/` tetap sebagai middleware kondisional.
3. **Multipart / file upload** — handler `product`, `category`, `slider`, `merchant_document`, `review_detail`, `transaction` memakai `c.FormValue`/`c.FormFile`. Konversi: `r.FormValue`/`r.FormFile` (perhatikan `r.FormFile` mengembalikan **3 nilai** vs 2 di echo). `pkg/upload_image.ProcessImageUpload` diubah dari `echo.Context` → `http.ResponseWriter` (helper `writeJSONError` internal) — satu-satunya konsumennya adalah apigateway.
4. **JWT menyimpan user_id sebagai `int`** — `JWTAuth()` custom mengekstrak claim `sub` via `extractUserIDFromClaims` dan menyimpan `user_id` int ke context (echo-jwt lama menyimpan via `SuccessHandler`).
5. **Role validator berbasis gRPC** — `RoleValidatorGRPC` (bukan Kafka) di-rewrite ke chi middleware.
6. **Duplikasi mount `/api/order-item`** — echo mengizinkan dua `Group` dengan prefix sama; chi panik saat `Route()` duplikat. Rute query kini registrasi langsung di parent router dengan path penuh, command tetap pakai `Route()`.

## 3. Dependensi

**Ditambah:** `go-chi/chi/v5 v5.3.2`, `go-chi/cors v1.2.2`, `swaggo/http-swagger v1.3.4`.
**Tidak lagi di-import langsung:** `labstack/echo/v4`, `echo-jwt/v4`, `echo-swagger` (masih *indirect* via `shared/errors` yang dipakai modul `tests/`).

## 4. Struktur Baru

```
service/apigateway/
├── apierror/          # NEW — ApiHandler versi net/http (AppError + HTTPError aware)
├── httpx/             # NEW — JSON/Bind/RealIP, request-scoped values, error adapter
├── apps/client.go     # bootstrap chi (Trace → Metrics → Recover → ... → JWT)
├── handler/           # 43 file → chi.Router
├── middlewares/       # auth, role(+gRPC), requiredRole, ratelimit, metrics, trace, pyroscope
└── ...
pkg/upload_image/      # ProcessImageUpload(w http.ResponseWriter, ...)
```

## 5. Pemetaan Konversi (inti)

Sama dengan payment-gateway; ringkas:

| Echo | Chi / net/http |
|---|---|
| `func(c echo.Context) error` | `func(w http.ResponseWriter, r *http.Request) error` |
| `c.FormValue` / `c.FormFile` (2 nilai) | `r.FormValue` / `r.FormFile` (3 nilai: file, header, err) |
| `RegisterErrorHandler(e)` global | adapter `httpx.Handler` + `apierror.HandleApiError` |
| `echo.NewHTTPError` → handler default `{"message":...}` | `httpx.NewHTTPError` + `httpx.WriteHTTPError` (format sama) |
| `c.Response().Status` (metrics/trace) | `wrapResponseWriter.Status()` |
| `e.GET/POST/DELETE` , `:param`, `Group("")` | `r.Get/Post/Delete`, `{param}`, `Route("")` + `.Get("/")` |
| Chain `m1(m2(h.X))` | `r.With(m1, m2).Get("/", ...)` |

Urutan middleware global dipertahankan: `Trace → HTTPMetrics → Recoverer → RequestID → Logger(JSON) → Pyroscope → CORS → Compress → Secure → RateLimit(/api/auth/) → JWTAuth`.

## 6. Verifikasi

```bash
cd monolith-ecommerce-grpc/service/apigateway
go build ./...        # ✅
go vet ./...          # ✅
go test -count=1 ./...  # ✅ handler (gateway_test: inventory rute, swagger↔rute, smoke 401/404/400/503) + middlewares
```

- `handler/gateway_test.go` di-rewrite ke `chi.Walk` untuk inventory rute (chi `Routes()` tidak punya field method) + normalisasi trailing slash pada root subrouter.
- Harness `tests/` (39 file): `Deps{E: …}` → `{Router: …}`, bypass middleware `c.Set("user_id", …)` → `context.WithValue`, `RegisterErrorHandler(s.echo)` (no-op reinstall) dihapus. Catatan: klaim awal "`MockImageUpload` ikut signature baru" **tidak terbukti di tree** — mock masih membawa param `w http.ResponseWriter` (sisa echo) dan baru dibetulkan di §9.
- `shared` dan seluruh service gRPC tidak tersentuh; `pkg/upload_image` diubah karena satu-satunya konsumen adalah apigateway.

## 7. Catatan Perilaku

1. **Respons error tidak berubah:** `*AppError` → status asli + shape `ErrorResponse` (trace_id dari X-Request-ID untuk rute unwrapped, dari span OTel untuk rute wrapped); middleware rejection & 404/405 tetap format compact `{"message": "..."}`; 400 dari body/param tidak valid tetap 400.
2. **Smoke test 401/404/400/503 lulus** — matriks status pada middleware chain asli terjaga.
3. **Echo masih indirect dependency** via `shared/errors` (echo-typed, dipakai `tests/`).

## 8. Follow-up (Opsional)

- [ ] Smoke test end-to-end dengan skenario `hurl/` atau full stack docker.
- [ ] `swag init` regenerasi docs bila perlu (anotasi `@Router` tidak berubah).
- [x] Migrasi `tests/` penuh ke net/http + pb per-domain — selesai 2026-09-14 (§9). **Sisa:** `shared/errors` masih echo-typed (2 file: `api_errors.go`, `merchant_document_errors/errors.go`) → setelah dibersihkan, `labstack/echo` hilang total dari graph.
- [x] Titik yang sama untuk apigateway `monolith-pointofsale-grpc` — selesai 2026-09-14, lihat `monolith-pointofsale-grpc/SUMMARY_ECHO_TO_CHI.md`.

## 9. Delta 2026-09-14 — `tests/` ke pb per-domain + `pbcompat` keluar `internal/`

Scope sesi ini: HTTP layer sudah chi (bagian di atas); yang masih merah cuma modul `tests/`. Prod kode (`service/*`, `shared`, `pkg`) tidak tersentuh selain pemindahan paket.

| Item | Sebelum | Sesudah |
|------|---------|---------|
| Paket compat | `service/apigateway/internal/pbcompat` | `service/apigateway/pbcompat` (64 file import di-rewire, `rmdir internal`) |
| Referensi pb di tests | 67 file import flat `monolith-ecommerce-pb` (root) → `pb.X` untuk semua domain | `pb<domain>.X` (`pbusers`, `pbcategories`, `pborders`, …); root `pb` hanya untuk simbol auth yang memang tinggal di root (`pb.LoginRequest`, `pb.RegisterAuthServiceServer`, dst.) |
| `tests/mock_upload.go` | `ProcessImageUpload(w http.ResponseWriter, uploadDir, file, isDocument)` | `ProcessImageUpload(uploadDir, file, isDocument)` — sesuai `pkg/upload_image.ImageUploads` sekarang; import `net/http` dihapus |
| `tests/order/stats_handler_gapi_test.go` | 1 field `client pborders.OrderStatsServiceClient` dipakai untuk 8 method (dulu satu service flat berisi semuanya) | `client` (4 method core) + `clientByMerchant pborders.OrderStatsByMerchantServiceClient` (4 method `*ByMerchant`) |
| `tests/orchestration_helpers.go` `SetupOrderService` | register 3 service (`OrderQuery`/`OrderStats`/`OrderCommand`) → panggilan `*ByMerchant` balas `rpc error: Unimplemented desc = unknown service pb.order.stats.OrderStatsByMerchantService` (HTTP 500 di api suite, error di gapi suite) | + `pborders.RegisterOrderStatsByMerchantServiceServer(server, orderGapi.OrderStatsByMerchant)`, disamakan dengan `service/order/apps/server.go:114-117` |
| `tests/go.mod` / `go.sum` | stale untuk `GOWORK=off` (`go: updates to go.mod needed`) | `GOWORK=off go mod tidy` — `labstack/echo/v4 v4.15.0` keluar dari direct require, `go-chi/chi/v5 v5.3.2` masuk sebagai direct; `monolith-ecommerce-pb` v1.0.0 → v1.0.10 (tetap `replace => ../pb`) |

### Cara rewire (reusable)

Regex/prefix heuristic tidak dipakai. Dibangun **symbol table exact** dari generated code: `type|func|var|const [A-Z]…` per file `pb/**/*.go` → 732 simbol, **0 collision** antar paket; 596 referensi `pb.X` di 67 file resolve unik semua (0 unresolved → script abort kalau ada yang tidak resolve, jadi tidak ada simbol yang "ditebak"). Alias target = `pb<nama-direktori>`, konsisten dengan `monolith-pointofsale-grpc`.

Sisa yang benar-benar manual cuma 2 kasus semantik di tabel atas (interface yang ter-split & mock signature) — keduanya hasil `go vet`, bukan pembacaan manual.

### Verifikasi

| Perintah | Hasil |
|----------|-------|
| `go build ./...` + `go vet ./...` × 27 modul | exit 0 (noise pre-existing: `pkg/server/server.go:75` cancel leak, `unreachable code` di 3 `handler/mapping.go` service merchant_*) |
| `gofmt -l tests pb shared service pkg` | 0 file |
| `just build` | 23/23 service → `bin/` |
| `cd service/apigateway && go test -count=1 ./...` | ok (inventory rute `chi.Walk` + smoke 401/404/400/503) |
| `just test-unit` (`pkg`) | exit 0 |
| `just test-integration` (`tests`, `GOWORK=off`) | exit 0 — 19 paket `ok`, 88 `--- PASS`, 0 `FAIL`. Run pertama: paket `order` FAIL (2 suite) karena jebakan #4 di bawah, hijau setelah registrasi ditambahkan |

### Jebakan baru yang tercatat (di luar §3)

1. **`just test-integration` pakai `GOWORK=off`** → `go mod tidy` wajib dijalankan *tanpa workspace* juga; hasil tidy dalam mode workspace tidak selalu cukup.
2. **Satu service proto bisa ter-split jadi beberapa gRPC service** — symbol table memetakan *nama* ke *paket*, tapi tidak tahu *method* milik interface yang mana (`OrderStatsServiceClient` vs `OrderStatsByMerchantServiceClient`). Tipe field interface di struct test harus dicek ulang terhadap `pb/<domain>/*_grpc.pb.go`.
3. **`internal/` again** — sama seperti pointofsale: letakkan paket compat di path publik kalau modul `tests/` butuh memanggil ctornya.
4. **Harness gRPC bisa tertinggal dari prod bootstrap.** Setelah service di-split, `Setup*Service()` di `tests/orchestration_helpers.go` harus mendaftarkan *semua* sub-service. Deteksi tanpa membaca satu-satu:

   ```bash
   grep -rhoE "Register[A-Za-z]+ServiceServer" --include=*.go service | sort -u > /tmp/prod.txt
   grep -rhoE "Register[A-Za-z]+ServiceServer" --include=*.go tests  | sort -u > /tmp/harness.txt
   comm -23 /tmp/prod.txt /tmp/harness.txt   # ada di prod, tidak ada di harness
   ```

   Hasil sesi ini: 4 nama, yang benar-benar dipakai test cuma `RegisterOrderStatsByMerchantServiceServer` (merchant_document & merchant_social memang tidak punya suite).
