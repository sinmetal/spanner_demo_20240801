FROM gcr.io/distroless/static-debian11
COPY ./cmd/app /app
COPY ./static/ /static
ENTRYPOINT ["/app"]