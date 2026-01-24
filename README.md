# 🚀 TaskFlow - RESTful API Task Management

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![Gin Framework](https://img.shields.io/badge/Gin-Framework-ff6b6b?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean_Architecture-brightgreen)
![Database](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)

**TaskFlow** là hệ thống Backend quản lý công việc (Task Management System) được xây dựng với hiệu năng cao, bảo mật và dễ bảo trì. Dự án áp dụng triệt để mô hình **Clean Architecture** (Layered Architecture) để tách biệt các tầng xử lý, giúp code dễ dàng mở rộng và viết Unit Test.

## 🛠 Tech Stack (Công nghệ sử dụng)

* **Core:** Golang (Go 1.23)
* **Framework:** Gin Gonic (Web Framework)
* **Database:** PostgreSQL, GORM (ORM Code-first)
* **Authentication:** JWT (JSON Web Token)
* **Security:** Rate Limiting (Chống Spam), CORS, Password Hashing (Bcrypt)
* **Infrastructure:** Docker, Docker Compose
* **Testing:** Testify (Unit Test with Mocking)

## ✨ Tính năng chính

* 🔐 **Authentication:** Đăng ký, Đăng nhập, JWT Authorization.
* 👥 **Team Management:** Tạo nhóm, thêm/xóa thành viên, phân quyền.
* 📋 **Task Management:** CRUD Task, gán người làm (Assignee), chuyển trạng thái (Todo -> In Progress -> Done).
* 🛡️ **Rate Limiting:** Giới hạn số lượng request để bảo vệ hệ thống.
* 🐳 **Dockerized:** Đóng gói sẵn sàng để triển khai chỉ với 1 lệnh.

## 📂 Cấu trúc dự án (Project Structure)

Dự án tuân thủ chuẩn Clean Architecture của Golang:

```text
TaskFlow/
├── cmd/
│   └── api/            # Entry point (Main), nơi khởi chạy server
├── internal/
│   ├── config/         # Load biến môi trường (.env)
│   ├── models/         # Định nghĩa Database Schemas (GORM struct)
│   ├── handlers/       # (Controller) Nhận request, validate dữ liệu đầu vào
│   ├── services/       # (UseCase) Chứa Logic nghiệp vụ (Business Logic)
│   ├── repositories/   # (Data Access) Tương tác trực tiếp với Database
│   ├── middlewares/    # Xử lý trung gian (Auth, CORS, RateLimit)
│   └── routes/         # Định nghĩa đường dẫn API (Endpoints)
├── pkg/                # Các thư viện dùng chung (Logger, Utils)
├── Dockerfile          # Cấu hình đóng gói Image
└── docker-compose.yml  # Cấu hình chạy toàn bộ hệ thống (App + DB)


## 🛠 Kịch bản Demo

🟢 GIAI ĐOẠN 1: TẠO TÀI KHOẢN (Setup User)
URL: POST http://localhost:8080/api/v1/auth/register
{
    "email": "hung.leader@gmail.com",
    "password": "password123",
    "full_name": "Nguyen Phi Hung"
}
----------------------------------------------------------------
Đăng ký Member
POST http://localhost:8080/api/v1/auth/register
{
    "email": "nam.dev@gmail.com",
    "password": "password123",
    "full_name": "Tran Van Nam"
}
----------------------------------------------------------------
Đăng nhập (Lấy Token)
POST http://localhost:8080/api/v1/auth/login
{
    "email": "hung.leader@gmail.com",
    "password": "password123"
}
----------------------------------------------------------------
🔵 GIAI ĐOẠN 2: THIẾT LẬP NHÓM (Team)
POST http://localhost:8080/api/v1/teams
{
    "name": "Backend Super Squad",
    "description": "Team chuyên trị bug khó"
}
----------------------------------------------------------------
Thêm vào Team
POST http://localhost:8080/api/v1/teams/{{TEAM_ID}}/members
{
    "user_id": "{{USER_ID_B}}"
}
----------------------------------------------------------------
🟠 GIAI ĐOẠN 3: QUẢN LÝ CÔNG VIỆC
POST http://localhost:8080/api/v1/tasks
{
    "title": "Fix bug login",
    "description": "Login đang bị lỗi 500 khi sai pass",
    "priority": "high",
    "status": "todo",
    "team_id": "{{TEAM_ID}}",
    "assigned_to": "{{USER_ID_B}}",
    "due_date": "2026-02-01T17:00:00Z"
}
----------------------------------------------------------------
Comment chỉ đạo
POST http://localhost:8080/api/v1/tasks/{{TASK_ID}}/comments
{
    "content": "Nam ơi, task này gấp nhé, xong trước 5h chiều!"
}
----------------------------------------------------------------
Upload ảnh minh họa lỗi
POST http://localhost:8080/api/v1/tasks/{{TASK_ID}}/attachments
----------------------------------------------------------------
🟣 GIAI ĐOẠN 4: CẬP NHẬT & BÁO CÁO
PATCH http://localhost:8080/api/v1/tasks/{{TASK_ID}}/status
{
    "status": "in_progress"
}
----------------------------------------------------------------
Xem Dashboard (Thống kê)
URL: GET http://localhost:8080/api/v1/tasks/dashboard/stats?team_id={{TEAM_ID}}
[
    { "status": "in_progress", "count": 1 }
]
----------------------------------------------------------------