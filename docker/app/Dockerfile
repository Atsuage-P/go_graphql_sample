FROM golang:1.22.5 as dev

RUN go install github.com/volatiletech/sqlboiler/v4@latest \
  && go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest

