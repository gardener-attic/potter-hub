FROM eu.gcr.io/gardenlinux/gardenlinux:184.0
RUN apt-get -y update && apt-get -y install ca-certificates curl jq && \
    update-ca-certificates && \
    curl -o /bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/v1.17.4/bin/linux/amd64/kubectl && chmod 755 /bin/kubectl

# Disable start of Berkeley DB
# copied installation package files from https://github.wdf.sap.corp/devx-wing/noberkeley/wiki/NoBerkeley-Packages
COPY noberkeley/noberkeley_1.0.0-3_amd64.deb .
COPY noberkeley/noberkeley-dev_1.0.0-3_amd64.deb .
RUN apt-get -y install ./noberkeley_1.0.0-3_amd64.deb ./noberkeley-dev_1.0.0-3_amd64.deb && \
    rm noberkeley_1.0.0-3_amd64.deb && \
    rm noberkeley-dev_1.0.0-3_amd64.deb 

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