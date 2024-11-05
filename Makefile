# Obtener la versión desde el hash de git
git_hash := $(shell git rev-parse --short HEAD || echo 'development')
$(info git_hash: $(git_hash))

# Archivo de versión
version_file := version.txt

# Definir la versión inicial
initial_version := 0.0.0

# Verificar si version.txt existe y crear si no
ifeq ($(wildcard $(version_file)),)
    $(shell echo "initial_version: $(initial_version)" > $(version_file))
    $(shell echo "git_hash: $(git_hash)" >> $(version_file))
endif

# Obtener la última versión y git_hash del archivo version.txt
last_version := $(shell awk -F': ' '/initial_version:/ {print $$2}' $(version_file) | xargs)
last_git_hash := $(shell awk -F': ' '/git_hash:/ {print $$2}' $(version_file) | xargs)

$(info last_version: $(last_version))
$(info last_git_hash: $(last_git_hash))

# Verificar si el hash de git ha cambiado
ifeq ($(strip $(git_hash)), $(strip $(last_git_hash)))
    next_version := $(last_version)
else
    # Función para incrementar la versión
    next_version := $(shell \
      major=$$(echo $(last_version) | awk -F. '{print $$1}'); \
      minor=$$(echo $(last_version) | awk -F. '{print $$2}'); \
      patch=$$(echo $(last_version) | awk -F. '{print $$3}'); \
      if [ $$patch -eq 9 ]; then \
        if [ $$minor -eq 9 ]; then \
          echo "$$(($$major + 1)).0.0"; \
        else \
          echo "$$major.$$((minor + 1)).0"; \
        fi; \
      else \
        echo "$$major.$$minor.$$((patch + 1))"; \
      fi)
    # Actualizar version.txt
    $(shell echo "initial_version: $(next_version)" > $(version_file))
    $(shell echo "git_hash: $(git_hash)" >> $(version_file))
endif
$(info next_version: $(next_version))

# Obtener la fecha actual
current_time = $(shell date +"%Y%m%d-%H%M%S")
$(info current_time: $(current_time))

all: run

build:
	rm -fr dist/parsewatchdog-*
	docker build --progress=plain \
		--build-arg VERSION=$(next_version) \
		--build-arg GIT_HASH=$(git_hash) \
		--build-arg BUILD_DATE=$(current_time) \
		-t parsewatchdog_createbinary \
		-f ./docker/Dockerfile .

run:
	docker run --rm --privileged \
		-v ./:/root/go/src/parsewatchdog \
		-v ./:/usr/local/go/src/parsewatchdog \
		-v ./dist:/root/go/dist \
		-v ./go-data-pkg:/home/go/pkg/mod \
		-e APP_VERSION=$(next_version) \
		-e GIT_HASH=$(git_hash) \
		-e BUILD_DATE=$(current_time) \
		parsewatchdog_createbinary
