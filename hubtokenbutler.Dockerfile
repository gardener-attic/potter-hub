FROM eu.gcr.io/gardener-project/3rd/golang:1.17.5 as builder
WORKDIR /app

# Copy relevant code
COPY go.mod go.sum ./
COPY pkg/ pkg/
COPY cmd/hub-token-butler/ cmd/hub-token-butler/

# Compile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o hub-token-butler cmd/hub-token-butler/main.go

# Create appuser
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN
# and https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

RUN adduser \
--disabled-password \
--gecos "" \
--home "/nonexistent" \
--shell "/sbin/nologin" \
--no-create-home \
--uid "${UID}" \
"${USER}"

FROM eu.gcr.io/sap-gcp-cp-k8s-stable-hub/3rd/kubeapps-chart-repo:1.10.0
WORKDIR /

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder  /app/hub-token-butler /

# Use an unprivileged user.
USER ${USER}:${USER}

CMD [ "/hub-token-butler" ]
