# gotodo

> Công cụ quản lý công việc ngay trên terminal — nhanh, đẹp, không cần rời khỏi dòng lệnh.

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)

---

## Tính năng

- ✅ **Quản lý đầy đủ** — thêm, xem danh sách, hoàn thành, chỉnh sửa và xóa công việc
- 🎨 **Giao diện terminal đẹp** — màu sắc theo độ ưu tiên, emoji trực quan, bảng hiển thị linh hoạt theo `go-pretty`
- 🏷️ **Tags & độ ưu tiên** — phân loại công việc bằng tag tùy chỉnh và 3 mức độ ưu tiên `high / medium / low`
- 📅 **Ngày tháng thân thiện** — nhập ngôn ngữ tự nhiên như `today`, `tomorrow`, `next week`
- ⏰ **Phát hiện quá hạn** — công việc đã quá ngày được tô đỏ tự động
- 🔍 **Lọc linh hoạt** — lọc theo trạng thái, độ ưu tiên, tag hoặc công việc đến hạn hôm nay
- 📐 **Bảng tự co giãn** — tự động ẩn/hiện cột phù hợp với độ rộng terminal
- 💾 **Không phụ thuộc runtime** — dữ liệu lưu dưới dạng file JSON trong thư mục cấu hình của hệ điều hành

---

## Cài đặt

### Build từ mã nguồn

```bash
git clone https://github.com/iamminhquan/gotodo.git
cd gotodo

# Build file thực thi
make build

# Hoặc cài vào $GOPATH/bin
make install

# Hoặc dùng go build
go build -o .\gotodo.exe .\cmd\gotodo\main.go
```

### Yêu cầu

- Go 1.21 trở lên

---

## Bắt đầu nhanh

```bash
# Thêm một công việc
gotodo add "Mua đồ ăn" --priority high --due tomorrow --tags cá-nhân,việc-vặt

# Xem danh sách công việc đang chờ (mặc định)
gotodo list

# Đánh dấu hoàn thành (dùng vài ký tự đầu của ID)
gotodo done a1b2c3d4

# Chỉnh sửa công việc
gotodo edit a1b2c3d4 --priority low --due 2026-05-15

# Xóa công việc
gotodo delete a1b2c3d4
```

---

## Các lệnh

### `gotodo` / `gotodo list`

Hiển thị danh sách công việc dưới dạng bảng màu sắc. Chạy `gotodo` không có subcommand tương đương với `gotodo list`.

```
gotodo list [flags]
```

| Flag | Viết tắt | Mô tả |
|------|----------|-------|
| `--pending` | | Chỉ hiện công việc chưa xong **(mặc định)** |
| `--done` | | Chỉ hiện công việc đã hoàn thành |
| `--all` | | Hiện tất cả công việc |
| `--today` | | Hiện công việc đến hạn hôm nay |
| `--priority` | `-p` | Lọc theo độ ưu tiên (`high`, `medium`, `low`) |
| `--tag` | `-t` | Lọc theo tag |

**Ví dụ:**

```bash
gotodo list
gotodo list --all
gotodo list --done
gotodo list --priority high --tag work
gotodo list --today
```

---

### `gotodo add`

Tạo công việc mới với độ ưu tiên, ngày hạn và tag tùy chọn.

```
gotodo add <tiêu đề> [flags]
```

| Flag | Viết tắt | Mô tả |
|------|----------|-------|
| `--priority` | `-p` | Độ ưu tiên: `high`, `medium`, `low` (mặc định: `medium`) |
| `--due` | `-d` | Ngày hết hạn: `YYYY-MM-DD`, `today`, `tomorrow`, `next week` |
| `--tags` | `-t` | Danh sách tag, cách nhau bằng dấu phẩy |

**Ví dụ:**

```bash
gotodo add "Đọc sách Clean Code" --priority high --due 2026-05-01 --tags đọc-sách,học-tập
gotodo add "Mua cà phê" --due tomorrow
gotodo add "Họp nhóm hàng tuần" --priority medium --tags công-việc
```

> **Mẹo:** Nếu không cung cấp `--tags`, gotodo sẽ tự động gán ngẫu nhiên 1–2 tag mặc định từ bộ: `work`, `personal`, `urgent`, `home`, `learning`, `health`, `finance`, `todo`.

---

### `gotodo done`

Đánh dấu một công việc là đã hoàn thành. Dùng `--undo` để chuyển lại trạng thái chờ.

```
gotodo done <id> [flags]
```

| Flag | Mô tả |
|------|-------|
| `--undo` | Chuyển công việc đã xong về trạng thái chờ |

`<id>` có thể là UUID đầy đủ hoặc chỉ vài ký tự đầu (tìm theo tiền tố).

**Ví dụ:**

```bash
gotodo done a1b2c3d4
gotodo done a1b2c3d4 --undo
```

---

### `gotodo edit`

Chỉnh sửa một hoặc nhiều trường của công việc hiện có. Chỉ những trường được cung cấp mới bị cập nhật.

```
gotodo edit <id> [flags]
```

| Flag | Viết tắt | Mô tả |
|------|----------|-------|
| `--title` | `-T` | Tiêu đề mới |
| `--priority` | `-p` | Độ ưu tiên mới: `high`, `medium`, `low` |
| `--due` | `-d` | Ngày hạn mới: `YYYY-MM-DD`, `today`, `tomorrow` |
| `--tags` | `-t` | Tags mới, cách nhau bằng dấu phẩy (thay thế toàn bộ tags cũ) |
| `--clear-due` | | Xóa ngày hạn khỏi công việc |

**Ví dụ:**

```bash
gotodo edit a1b2c3d4 --title "Mua rau củ hữu cơ"
gotodo edit a1b2c3d4 --priority low --due 2026-05-15
gotodo edit a1b2c3d4 --tags công-việc,khẩn-cấp
gotodo edit a1b2c3d4 --clear-due
```

---

### `gotodo delete`

Xóa vĩnh viễn một công việc. Sẽ có bước xác nhận trừ khi dùng `--force`.

```
gotodo delete <id> [flags]
```

| Flag | Viết tắt | Mô tả |
|------|----------|-------|
| `--force` | `-f` | Bỏ qua bước xác nhận |

**Ví dụ:**

```bash
gotodo delete a1b2c3d4
gotodo delete a1b2c3d4 --force
```

---

## Độ ưu tiên

| Ký hiệu | Mức độ | Ý nghĩa |
|---------|--------|---------|
| 🔴 | `high` | Khẩn cấp, cần xử lý ngay |
| 🟡 | `medium` | Mức mặc định cho hầu hết công việc |
| 🔵 | `low` | Làm được thì tốt, không quá cấp bách |

---

## Định dạng ngày hạn

gotodo chấp nhận cả ngày tuyệt đối lẫn ngôn ngữ tự nhiên khi nhập ngày:

| Giá trị nhập | Hiểu là |
|--------------|---------|
| `2026-05-01` | Ngày 1 tháng 5 năm 2026 |
| `today` | Ngày hôm nay |
| `tomorrow` | Ngày mai |
| `next week` | 7 ngày kể từ hôm nay |

Công việc **quá hạn** được tô đỏ với ký hiệu ⚠. Công việc **đến hạn hôm nay** hiển thị màu vàng với 🔔.

---

## Cấu hình

gotodo đọc cấu hình từ file YAML khi khởi động. File cấu hình nằm ở thư mục chuẩn của từng hệ điều hành:

| Hệ điều hành | Thư mục cấu hình |
|-------------|-----------------|
| Linux | `~/.config/gotodo/config.yaml` |
| macOS | `~/Library/Application Support/gotodo/config.yaml` |
| Windows | `%APPDATA%\gotodo\config.yaml` |

Nếu không có file cấu hình, gotodo sẽ chạy với các giá trị mặc định.

### Các tuỳ chọn cấu hình

```yaml
# Backend lưu trữ (hiện chỉ hỗ trợ "json")
storage_type: json

# Thư mục chứa file tasks.json
data_dir: ""   # mặc định theo thư mục data của hệ điều hành

# Độ ưu tiên mặc định khi không truyền --priority
default_priority: medium

# Tắt màu ANSI (hữu ích khi dùng terminal không hỗ trợ màu)
use_color: true

# Tên file lưu công việc trong data_dir
tasks_file_name: tasks.json
```

### Vị trí file dữ liệu

Công việc được lưu dưới dạng JSON. Vị trí mặc định:

| Hệ điều hành | Đường dẫn file |
|-------------|---------------|
| Linux | `~/.local/share/gotodo/tasks.json` |
| macOS | `~/Library/Application Support/gotodo/tasks.json` |
| Windows | `%APPDATA%\gotodo\data\tasks.json` |

Bạn có thể ghi đè đường dẫn này qua biến môi trường `GOTODO_DATA_DIR` hoặc config key `data_dir`.

### Biến môi trường

Tất cả config key đều có thể ghi đè qua biến môi trường với tiền tố `GOTODO_`:

```bash
GOTODO_DEFAULT_PRIORITY=high gotodo add "Việc quan trọng"
GOTODO_USE_COLOR=false gotodo list
GOTODO_DATA_DIR=/đường/dẫn/tùy-chỉnh gotodo list
```

---

## Phát triển

```bash
# Chạy toàn bộ unit test với race detector
make test

# Kiểm tra code với go vet
make lint

# Build file thực thi
make build

# Build và chạy với tham số
make run ARGS="list --all"

# Dọn dẹp file build
make clean
```

---

## Cấu trúc dự án

```
gotodo/
├── cmd/gotodo/         # Điểm vào chương trình (main.go)
├── internal/
│   ├── cli/            # Các Cobra command (add, list, done, edit, delete)
│   ├── config/         # Đọc cấu hình (Viper + configdir)
│   ├── storage/        # Backend lưu trữ JSON
│   ├── task/           # Model dữ liệu và logic nghiệp vụ
│   └── version/        # Chuỗi phiên bản (inject qua ldflags)
├── Makefile
└── go.mod
```

---

## Giấy phép

Dự án được phát hành theo [Giấy phép MIT](LICENSE).
