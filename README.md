
```
sentra-be
├─ .dockerignore
├─ Dockerfile
├─ LICENSE
├─ Makefile
├─ README.md
├─ bin
│  └─ start
├─ cmd
│  └─ app
│     └─ main.go
├─ coverage.html
├─ database
│  ├─ migrations
│  │  ├─ 20250104035256_user.down.sql
│  │  ├─ 20250104035256_user.up.sql
│  │  ├─ 20250407045854_budget_manager.down.sql
│  │  ├─ 20250407045854_budget_manager.up.sql
│  │  ├─ 20250408185549_wallet.down.sql
│  │  ├─ 20250408185549_wallet.up.sql
│  │  ├─ 20250408185624_wallet_transaction.down.sql
│  │  └─ 20250408185624_wallet_transaction.up.sql
│  └─ postgres
│     └─ postgres.go
├─ docker-compose.backend-only.yml 
├─ docker-compose.yml
├─ go.mod
├─ go.sum
├─ internal
│  ├─ api
│  │  ├─ auth
│  │  │  ├─ dto.go
│  │  │  ├─ error.go
│  │  │  ├─ handler
│  │  │  │  ├─ authentication_hd.go
│  │  │  │  ├─ http.go
│  │  │  │  ├─ oauth_hd.go
│  │  │  │  ├─ password_hd.go
│  │  │  │  └─ profile_hd.go
│  │  │  ├─ repository
│  │  │  │  ├─ query.go
│  │  │  │  ├─ repository.go
│  │  │  │  └─ user_rs.go
│  │  │  └─ service
│  │  │     ├─ auth_sv.go
│  │  │     ├─ biometric_sv.go
│  │  │     ├─ helper.go
│  │  │     ├─ password_sv.go
│  │  │     ├─ service.go
│  │  │     └─ user_sv.go
│  │  ├─ budget_manager
│  │  │  ├─ dto.go
│  │  │  ├─ error.go
│  │  │  ├─ handler
│  │  │  │  ├─ budget_hd.go
│  │  │  │  └─ http.go
│  │  │  ├─ repository
│  │  │  │  ├─ budget_rs.go
│  │  │  │  ├─ query.go
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     ├─ budget_sv.go
│  │  │     └─ service.go
│  │  ├─ detection
│  │  │  ├─ dto.go
│  │  │  ├─ error.go
│  │  │  ├─ handler
│  │  │  │  ├─ detection_hd.go
│  │  │  │  └─ http.go
│  │  │  └─ service
│  │  │     ├─ detection_sv.go
│  │  │     └─ service.go
│  │  └─ sentra_pay
│  │     ├─ dto.go
│  │     ├─ error.go
│  │     ├─ handler
│  │     │  ├─ http.go
│  │     │  └─ sentra_pay_hd.go
│  │     ├─ repository
│  │     │  ├─ query.go
│  │     │  ├─ repository.go
│  │     │  └─ wallet_rs.go
│  │     └─ service
│  │        ├─ sentra_pay_sv.go
│  │        └─ service.go
│  ├─ config
│  │  ├─ fiber.go
│  │  ├─ logrus.go
│  │  ├─ rest.go
│  │  ├─ server.go
│  │  └─ validator.go
│  ├─ entity
│  │  ├─ budget_manager.go
│  │  ├─ face_detection.go
│  │  ├─ ktp_detection.go
│  │  ├─ qris_detection.go
│  │  ├─ session.go
│  │  ├─ user.go
│  │  └─ whatsapp.go
│  ├─ middleware
│  │  ├─ logging.go
│  │  ├─ middleware.go
│  │  ├─ rate_limitter.go
│  │  ├─ request_id.go
│  │  └─ token.go
│  └─ mocks
│     ├─ redis_mock.go
│     ├─ s3_mock.go
│     ├─ smtp_mock.go
│     ├─ websocket_mock.go
│     └─ whatsapp_mock.go
├─ nginx
│  ├─ conf.d
│  │  └─ default.conf
│  ├─ logs
│  └─ ssl
├─ pkg
│  ├─ bcrypt
│  │  ├─ bcrypt.go
│  │  └─ bcrypt_test.go
│  ├─ context
│  │  ├─ context.go
│  │  └─ context_test.go
│  ├─ doku
│  │  └─ doku.go
│  ├─ gemini
│  │  └─ gemini.go
│  ├─ google
│  │  └─ google.go
│  ├─ handlerUtil
│  │  ├─ handler_util.go
│  │  └─ handler_util_test.go
│  ├─ jwt
│  │  ├─ jwt_test.go
│  │  └─ token.go
│  ├─ log
│  │  ├─ log.go
│  │  └─ log_test.go
│  ├─ redis
│  │  ├─ redis.go
│  │  └─ redis_test.go
│  ├─ response
│  │  ├─ response.go
│  │  └─ response_test.go
│  ├─ s3
│  │  ├─ s3.go
│  │  └─ s3_test.go
│  ├─ smtp
│  │  ├─ smtp.go
│  │  └─ smtp_test.go
│  ├─ utils
│  │  ├─ utils.go
│  │  └─ utils_test.go
│  ├─ websocket
│  │  ├─ websocket.go
│  │  └─ websocket_test.go
│  └─ whatsapp
│     ├─ whatsapp.go
│     └─ whatsapp_test.go
├─ storage
│  └─ logs
└─ tmp
   └─ main

```