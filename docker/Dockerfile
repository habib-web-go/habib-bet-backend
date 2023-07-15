FROM golang:1.20.4-bullseye AS builder

ADD go.mod go.sum Makefile /app/
WORKDIR /app
RUN make deps
ADD . /app
RUN make build

FROM debian:bullseye AS production


COPY --from=builder /app/habib-bet-backend /app/
ADD config/production.yaml  /app/config/
WORKDIR /app
ENTRYPOINT ["./habib-bet-backend"]
CMD ["--profile", "production", "--command", "run_server"]