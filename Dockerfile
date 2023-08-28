FROM golang:1.20-alpine as builder 

WORKDIR /
ADD . .

RUN go build -o /tmp/report-bot main.go

FROM alpine AS runner

COPY --from=builder /tmp/report-bot /bin/report-bot

RUN chmod +x /bin/report-bot


CMD ["/bin/report-bot"]