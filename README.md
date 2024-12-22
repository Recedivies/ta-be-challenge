<div align="center" style="padding-bottom: 20px">
  <h1>Paper.id Take Home Test</h1>
</div>

This is a simple wallet application that allows users to manage their wallet balance: top-up and withdrawal.

## Setup local development

### Install tools

- [Golang](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### How to run

- Clone the Repository

```bash
git clone https://github.com/Recedivies/ta-be-challenge.git
cd ta-be-challenge
```

Make sure Docker Engine is running.

- Build and Run Server with Docker Compose

  ```bash
  docker compose up --build
  ```

Application should be up and running: backend `127.0.0.1:9090`, postgres `127.0.0.1:5432`

- To run the test cases

  ```bash
  go test -v
  ```

- To stop and remove containers

  ```bash
  docker compose down
  ```
