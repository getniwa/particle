CORE=YOUR_CORE_NAME
GITHASH=$(shell git describe --tags --always)
OUTPUT=firmware-$(GITHASH).bin
MAIN=main.ino
SPARK=particle

all: compile flash

compile:
	$(SPARK) compile $(MAIN) $(OUTPUT)

flash:
	$(SPARK) flash $(CORE) $(OUTPUT)


