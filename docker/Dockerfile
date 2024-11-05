FROM debian:bookworm

LABEL maintainer="Federico Pereira <lord.basex@gmail.com>"

ENV GO111MODULE=on \
    GOPATH=/home/go \
    VERSION=1.23.2 \
    GO_URL=https://go.dev/dl \
    GOARCH=amd64 \
    GOOS=linux

# Instalar dependencias
RUN apt -y update && apt-get -y install git bash wget file

# Descargar y descomprimir Go
RUN wget ${GO_URL}/go${VERSION}.linux-${GOARCH}.tar.gz -O /usr/src/go${VERSION}.linux-${GOARCH}.tar.gz \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf /usr/src/go${VERSION}.linux-${GOARCH}.tar.gz

# Crear el script de inicio
RUN echo "#!/bin/bash\n\
echo 'Variables de entorno:'\n\
printenv | grep APP_VERSION\n\
printenv | grep GIT_HASH\n\
printenv | grep BUILD_DATE\n\
env GOOS=linux GOARCH=i386 /usr/local/go/bin/go build -ldflags=\"-s -w -X 'github.com/lordbasex/parsewatchdog/config.Version=\$APP_VERSION' -X 'github.com/lordbasex/parsewatchdog/config.DaemonGitBuild=\$GIT_HASH' -X 'github.com/lordbasex/parsewatchdog/config.DaemonGitBuildDate=\$BUILD_DATE'\" -o /root/go/dist/parsewatchdog-x86_64 cmd/main.go \n\
env GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -ldflags=\"-s -w -X 'github.com/lordbasex/parsewatchdog/config.Version=\$APP_VERSION' -X 'github.com/lordbasex/parsewatchdog/config.DaemonGitBuild=\$GIT_HASH' -X 'github.com/lordbasex/parsewatchdog/config.DaemonGitBuildDate=\$BUILD_DATE'\" -o /root/go/dist/parsewatchdog-x86_64 cmd/main.go \n\
env GOOS=linux GOARCH=arm64 /usr/local/go/bin/go build -ldflags=\"-s -w -X 'github.com/lordbasex/parsewatchdog/config.Version=\$APP_VERSION' -X 'github.com/lordbasex/parsewatchdog/config.DaemonGitBuild=\$GIT_HASH' -X 'github.com/lordbasex/parsewatchdog/config.DaemonGitBuildDate=\$BUILD_DATE'\" -o /root/go/dist/parsewatchdog-x86_64 cmd/main.go \n\
/usr/bin/file /root/go/dist/parsewatchdog-i386 \n\
/usr/bin/file /root/go/dist/parsewatchdog-x86_64 \n\
/usr/bin/file /root/go/dist/parsewatchdog-arm64 \n\
exec \"\$@\" \n\
" > /start.sh

# Dar permisos de ejecución al script
RUN chmod +x /start.sh

# Establecer el directorio de trabajo
WORKDIR /root/go/src/parsewatchdog

# Establecer el comando de inicio
CMD ["/bin/sh", "/start.sh"]

