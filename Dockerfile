FROM ubuntu:20.10
RUN mkdir /app
COPY ./bin/prometheus-rancher-exporter /app/
WORKDIR /app
RUN chmod +x /app/prometheus-rancher-exporter
CMD ["/app/prometheus-rancher-exporter"]
EXPOSE 8080
