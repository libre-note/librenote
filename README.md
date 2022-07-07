# LibreNote

Libre(Free as in freedom) note is a note taking applications. A alternative to google keep.

## Build & Run
- Local setup
  ```bash
  make help # get full list of make commands
  make run # test binary
  make serve # run the application
  ```
- Using docker
  ```bash
    make docker-build
    make docker-run
    # migrate sqlite inside
    make docker-migrate
  ```
