FROM bitnami/minideb:bullseye as builder

RUN apt-get -y update && apt-get -y install ca-certificates curl jq && \
    update-ca-certificates
    
RUN curl -o /bin/kubectl -L https://dl.k8s.io/release/v1.25.0/bin/linux/amd64/kubectl && chmod 755 /bin/kubectl

#### BASE ####
FROM gcr.io/distroless/static-debian11:nonroot AS base

COPY --from=builder /bin/kubectl /bin/kubectl

ENTRYPOINT [ "kubectl" ]
CMD [ "--help" ]