# wa-app

`wa-app` là dịch vụ chuỗi ứng dụng WA, cung cấp quản lý tài khoản, dò số, đăng ký, kiểm tra trạng thái đăng nhập, phiên kết nối dài và xử lý tin nhắn, kèm dashboard quản trị tích hợp.

> [!CAUTION]
> Khi dùng dự án này, bạn đồng ý toàn bộ điều khoản trong [NOTICE](./NOTICE). Dự án chỉ dành cho mô hình hóa giao thức, demo phục vụ học tập, nghiên cứu bảo mật có ủy quyền và kiểm thử nội bộ phi thương mại. Cấm dùng cho mục đích thương mại, mục tiêu không được ủy quyền, hoặc tình huống vi phạm điều khoản của bên thứ ba.

## Tính năng

- Quản lý tài khoản: lưu WAAccount, client profile, bản ghi đăng ký và projection trạng thái đăng nhập.
- Số điện thoại và đăng ký: hỗ trợ dò số, dò SMS, gửi yêu cầu đăng ký, nộp OTP và kiểm tra trạng thái đăng nhập.
- Kết nối và tin nhắn: hỗ trợ phiên kết nối dài, nhận tin nhắn, ack tin nhắn, gửi tin nhắn văn bản 1:1 và xem hội thoại.
- Trích xuất dữ liệu: trích xuất ứng viên OTP/Flag từ tin nhắn, rồi lưu tham chiếu hoặc projection đã khử nhạy cảm theo rule dữ liệu nhạy cảm.
- Giao diện quản trị: dashboard cho tài khoản, liên hệ, tin nhắn, trạng thái kết nối và thao tác hồ sơ tài khoản.

## Triển khai

Cách khuyến nghị để chạy dịch vụ là dùng Docker Compose có sẵn trong repo:

```sh
cp .env.example .env
docker compose pull
docker compose up -d
```

Nếu chỉ cần chạy local nhanh, có thể gọi thẳng `docker compose up -d`. Dịch vụ vẫn khởi động được khi chưa tạo `.env`; giá trị thiếu sẽ dùng mặc định trong compose.

Cổng mặc định (cố định):

- Dashboard: `http://127.0.0.1:8080` (mapping trong `docker-compose.yml`)
- gRPC: `127.0.0.1:50091`

Nếu thực sự cần đổi cổng map trên host, sửa trực tiếp dòng `ports` trong `docker-compose.yml`. Chúng không phải config option riêng.

### Cấu hình

`.env` chỉ giữ một số cấu hình runtime tối thiểu:

- `WA_APP_IMAGE_TAG`: tag image. Production nên ghim version cố định.
- `WA_APP_AUTH_PASSWORD`: mật khẩu đơn cho dashboard, tùy chọn. Để trống để tắt auth.
- `WA_APP_DATA_DIR`: thư mục persistent trong container. Mặc định `/var/lib/wa-app`.
- `WA_APP_PG_DSN`: DSN PostgreSQL, tùy chọn. Để trống để dùng SQLite tích hợp cho persistence.
- `WA_APP_REDIS_URL`: Redis URL, tùy chọn. Để trống để dùng SQLite tích hợp cho runtime state.
- `WA_COMMON_PROXY`: proxy WA outbound mặc định, tùy chọn. Để trống thì kết nối trực tiếp.
- `WA_NUMBER_PROBE_PROXY`: proxy cho dò số/SMS, tùy chọn. Nếu trống sẽ fallback sang `WA_COMMON_PROXY`; nếu cả hai cùng trống thì đi trực tiếp.
- `WA_REGISTRATION_PROXY`: proxy cho đăng ký và gửi OTP, tùy chọn. Nếu trống sẽ fallback sang `WA_COMMON_PROXY`; nếu cả hai cùng trống thì đi trực tiếp.

PostgreSQL và Redis đều là thành phần tùy chọn. Muốn bật, bỏ comment service tương ứng trong `docker-compose.yml` và điền `WA_APP_PG_DSN` / `WA_APP_REDIS_URL` trong `.env`.

### Build image từ mã nguồn

`Dockerfile` hỗ trợ build trực tiếp trong repo `wa-app`, không cần `common-lib` làm build context:

```sh
docker build -t wa-app-service:local .
```

Khi build từ thư mục tổng `byte-v-forge` cũng dùng được:

```sh
docker build -f wa-app/Dockerfile -t wa-app-service:local .
```

## Liên kết

- [LINUX DO - Cộng đồng lý tưởng mới](https://linux.do/)
