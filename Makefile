.phony: build clean run

PROJECT = chromos
SOURCES = $(wildcard src/*.go)
RUNTIME_RESOURCES = $(wildcard src/resources/*)
ICON_ASSETS = $(wildcard assets/*.icns)

BASEBUILD = build/$(PROJECT)-osx

OSXBUILD = $(BASEBUILD)/$(PROJECT).app/Contents
OSXLIB  = $(wildcard lib/osx/*.dylib)
OSXLIBD = $(subst lib/osx/,$(OSXBUILD)/MacOS/,$(OSXLIB))

YOSBUILD = $(BASEBUILD)-yosemite/$(PROJECT).app/Contents
YOSLIB  = $(wildcard lib/osx-yosemite/*.dylib)
YOSLIBD = $(subst lib/osx-yosemite/,$(YOSBUILD)/MacOS/,$(YOSLIB))

WINBUILD = build/$(PROJECT)-win
WINLIB  = $(wildcard lib/win/*.dll)
WINLIBD = $(subst lib/win/,$(WINBUILD)/,$(WINLIB))

NIXBUILD = build/$(PROJECT)-linux
NIXLIB  = $(wildcard lib/linux/*.*)
NIXLIBD = $(subst lib/linux/,$(NIXBUILD)/lib/,$(NIXLIB))

VERSION = $(shell cat VERSION)
REPLACE = s/9\.9\.9/$(VERSION)/g

clean:
	rm -rf build

$(OSXBUILD)/MacOS/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/Info.plist: pkg/osx/Info.plist
	mkdir -p $(OSXBUILD)
	sed $(REPLACE) $< > $@

$(OSXBUILD)/MacOS/%.dylib: lib/osx/%.dylib
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go
	cd $(OSXBUILD)/MacOS/ && ../../../../../scripts/fix.sh

$(OSXBUILD)/Resources/%.icns: assets/%.icns
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/Resources/resources/%: src/resources/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-osx-$(VERSION).zip: \
	$(OSXBUILD)/MacOS/launch.sh \
	$(OSXBUILD)/Info.plist \
	$(OSXLIBD) \
	$(OSXBUILD)/MacOS/$(PROJECT) \
	$(subst src/resources/,$(OSXBUILD)/Resources/resources/,$(RUNTIME_RESOURCES)) \
	$(subst assets/,$(OSXBUILD)/Resources/,$(ICON_ASSETS))
	cd build && zip -r $(notdir $@) $(PROJECT)-osx

$(YOSBUILD)/MacOS/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(YOSBUILD)/Info.plist: pkg/osx/Info.plist
	mkdir -p $(YOSBUILD)
	sed $(REPLACE) $< > $@

$(YOSBUILD)/MacOS/%.dylib: lib/osx-yosemite/%.dylib
	mkdir -p $(dir $@)
	cp $< $@

$(YOSBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go
	cd $(YOSBUILD)/MacOS/ && ../../../../../scripts/fix-yosemite.sh

$(YOSBUILD)/Resources/%.icns: assets/%.icns
	mkdir -p $(dir $@)
	cp $< $@

$(YOSBUILD)/Resources/resources/%: src/resources/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-osx-yosemite-$(VERSION).zip: \
	$(YOSBUILD)/MacOS/launch.sh \
	$(YOSBUILD)/Info.plist \
	$(YOSLIBD) \
	$(YOSBUILD)/MacOS/$(PROJECT) \
	$(subst src/resources/,$(YOSBUILD)/Resources/resources/,$(RUNTIME_RESOURCES)) \
	$(subst assets/,$(YOSBUILD)/Resources/,$(ICON_ASSETS))
	cd build && zip -r $(notdir $@) $(PROJECT)-osx-yosemite

$(WINBUILD)/$(PROJECT).exe: $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

$(WINBUILD)/%.dll: lib/win/%.dll
	mkdir -p $(dir $@)
	cp $< $@

$(WINBUILD)/resources/%: src/resources/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-win-$(VERSION).zip: \
	$(WINBUILD)/$(PROJECT).exe \
	$(WINLIBD) \
	$(subst src/resources/,$(WINBUILD)/resources/,$(RUNTIME_RESOURCES))
	cd build && /c/Program\ Files/7-Zip/7z.exe a -r $(notdir $@) $(PROJECT)-win

$(NIXBUILD)/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(NIXBUILD)/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

$(NIXBUILD)/lib/%: lib/linux/%
	mkdir -p $(dir $@)
	cp $< $@

$(NIXBUILD)/resources/%: src/resources/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-linux-$(VERSION).zip: \
	$(NIXBUILD)/launch.sh \
	$(NIXBUILD)/$(PROJECT) \
	$(NIXLIBD) \
	$(subst src/resources/,$(NIXBUILD)/resources/,$(RUNTIME_RESOURCES))
	cd build && zip -r $(notdir $@) $(PROJECT)-linux

build: build/$(PROJECT)-win-$(VERSION).zip

run: build
	$(OSXBUILD)/MacOS/launch.sh
