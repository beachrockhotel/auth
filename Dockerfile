FROM golang:1.20.3-alpine AS builder

ENTRYPOINT ["top", "-b"]