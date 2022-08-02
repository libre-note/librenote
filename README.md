# LibreNote
[![](https://img.shields.io/badge/Go-1.17+-00ADD8?style=flat&logo=go)](https://golang.org/doc/go1.17)
[![build](https://github.com/libre-note/librenote/actions/workflows/build.yml/badge.svg)](https://github.com/libre-note/librenote/actions?query=workflow%3ABuild)
[![gosec](https://img.shields.io/github/workflow/status/libre-note/librenote/Security?label=%F0%9F%94%91%20gosec&style=flat&color=75C46B)](https://github.com/libre-note/librenote/actions?query=workflow%3ASecurity)
[![codecov](https://codecov.io/gh/libre-note/librenote/branch/master/graph/badge.svg?token=N3JVSRO5NZ)](https://codecov.io/gh/libre-note/librenote)
[![Go Report Card](https://goreportcard.com/badge/github.com/libre-note/librenote)](https://goreportcard.com/report/github.com/libre-note/librenote)

Libre(Free as in freedom) note is a note taking applications. A alternative to google keep.

## ⚡️ Quick start
- Install **`docker`**, **`golang-migrate`**
- Copy config file `mv _doc/config ./` to root directory and change it
- Run project by this command:
  ```bash
    make docker-build
    make docker-run
    # migrate using sqlite
    make docker-migrate
  ```
- Visit **`http://localhost:8000`**
- Stop `CTRL + C`

## 🔨 Development
- Copy config file `mv _doc/config ./` to root directory and change it
- $
  ```bash
  make run # test binary
  make serve # run the application
  make migrate-up-sqlite
  make test-unit
  make test-integration # default sqlite
  make test-integration-mysql
  ``
- Check `make help` to get full list of make commands

## 🗒️ Docs
- [ERD Sqlite](_doc/sqlite_erd.png)
- [ERD Mysql](_doc/mysql_erd.png)
- [ERD Postgresql](_doc/postgresql_erd.png)
- [API Documentation](_doc/swagger.html)

