FROM eu.gcr.io/gardenlinux/gardenlinux:590.0-276f22-amd64-base-slim
RUN apt-get -y update && apt-get -y install ca-certificates curl jq && \
    update-ca-certificates && \
    curl -o /bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/v1.17.4/bin/linux/amd64/kubectl && chmod 755 /bin/kubectl

# Create appuser
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN

# DEBUG
RUN cat /etc/passwd

RUN adduser \
--disabled-password \
--gecos "" \
--home "/nonexistent" \
--shell "/sbin/nologin" \
--no-create-home \
--uid "${UID}" \
"$USER"

# Use an unprivileged user.
USER ${USER}:${USER}