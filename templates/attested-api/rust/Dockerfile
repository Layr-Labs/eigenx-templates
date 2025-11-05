FROM rust:latest AS builder

WORKDIR /app

COPY Cargo.toml ./
COPY src ./src

RUN cargo build --release

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/target/release/eigenx-tee-rust-app .

EXPOSE 8080

CMD ["./eigenx-tee-rust-app"]
