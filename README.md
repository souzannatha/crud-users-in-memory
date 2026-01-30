# Desafio 2 – CRUD de Usuários em Memória (Go)

Projeto desenvolvido como **Desafio 2 da Rocketseat**, com o objetivo de criar uma API REST simples em Go para gerenciamento de usuários em memória.

## 🎯 Objetivo

Praticar os fundamentos de:

* Go
* net/http
* Rotas REST
* Manipulação de dados em memória
* Uso correto de status HTTP

## 🛠 Tecnologias

* Go
* net/http
* chi
* uuid

## 📦 Estrutura

```
.
├── main.go
└── api/
    └── api.go
```

## 🚀 Como rodar

```bash
go mod init

go mod tidy

go run main.go
```

API disponível em:

```
http://localhost:8080
```

## 📚 Endpoints

* `POST /api/users`
* `GET /api/users`
* `GET /api/users/{id}`
* `PUT /api/users/{id}`
* `DELETE /api/users/{id}`

## ⚠️ Observações

* Dados armazenados apenas em memória
* Projeto voltado exclusivamente para estudo e prática

---

Desafio concluído 🚀
