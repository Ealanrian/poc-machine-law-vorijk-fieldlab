PROJECT_NAME := machine_law
PREBUILD_BASE := ../prebuild

NDK_VERSION ?=28.0.12674087

BUILD_MODE ?= release

ifeq ($(BUILD_MODE),debug)
	#Debug build
	GO_BUILD_FLAGS := -gcflags="all=-N -l"
else
	#Release build
	GO_BUILD_FLAGS := -trimpath -ldflags="-s -w"
endif

# Check os to figure out likely android sdk location
ifeq ($(OS), Windows_NT)
	#Windows platform (not really supported)
	ANDROID_HOME ?= $(LOCALAPPDATA/Android/Sdk)
	NDK_PLATFORM := windows-x86_64
else
	# Unix system Linux or Mac
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		ANDROID_HOME ?= $(HOME)/Android/Sdk
		NDK_PLATFORM := linux-x86_64
	else ifeq ($(UNAME_S),Darwin)  # macOS
		ANDROID_HOME ?= $(HOME)/Library/Android/sdk
		NDK_PLATFORM := darwin-x86_64
	endif
endif


NDK_BIN ?= /home/erik/Android/Sdk/ndk/$(NDK_VERSION)/toolchains/llvm/prebuilt/$(NDK_PLATFORM)/bin
#Check if the NDK is found.

NDK_EXISTS := $(if $(wildcard $(NDK_BIN)),true,false)

LIB_PREFIX = lib
ifeq ($(findstring ios,$(MAKECMDGOALS)),ios)
	PLATFORM = iOS
else ifeq ($(findstring android,$(MAKECMDGOALS)),android)
	PLATFORM = Android
else ifeq ($(findstring macos,$(MAKECMDGOALS)),macos)
	PLATFORM = macOS
else ifeq ($(MAKECMDGOALS),windows)
	PLATFORM = Windows
	LIB_PREFIX :=
else ifeq ($(MAKECMDGOALS),linux)
	PLATFORM = Linux
else ifeq ($(MAKECMDGOALS),web)
	PLATFORM = Web
else ifeq ($(MAKECMDGOALS),all)
else ifeq ($(MAKECMDGOALS),)
else
	PLATFORM = $(MAKECMDGOALS)
endif
PREBUILD_PATH = $(PREBUILD_BASE)/$(PLATFORM)
LIB_NAME := $(LIB_PREFIX)$(PROJECT_NAME)


.PHONY: release
all:
ifeq ($(OS),Windows_NT)
	$(MAKE) windows
else ifeq ($(UNAME_S),Linux)
	$(MAKE) linux
else ifeq ($(UNAME_S),Darwin)  # macOS
	$(MAKE) macos
	# Add the ios ask if on macos
	$(MAKE) ios
else
	$(error Unsupported OS for default build: $(UNAME_S))
endif
	#add android task
	@if [ "$(NDK_EXISTS)" = "true" ]; then \
		$(MAKE) android; \
	fi


debug:
	$(MAKE) BUILD_MODE=debug all

release:
	$(Make) BUILD_MODE=release all

.DEFAULT_GOALS := all


ANDROID_OUT=../android/src/main/jniLibs
NDK_BIN=/home/erik/Android/Sdk/ndk/$(NDK_VERSION)/toolchains/llvm/prebuilt/linux-x86_64/bin
LIB_NAME=machine_law

android-armv7a: CURRENT_ARCH := armeabi-v7a
android-armv7a:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm \
	GOARM=7 \
	CC=$(NDK_BIN)/armv7a-linux-androideabi21-clang \
	go build $(GO_BUILD_FLAGS) -buildmode=c-shared -o $(PREBUILD_PATH)/$(CURRENT_ARCH)/${LIB_NAME}.so .
	rm $(PREBUILD_PATH)/$(CURRENT_ARCH)/${LIB_NAME}.h


android-arm64: CURRENT_ARCH := arm64-v8a
android-arm64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm64 \
	CC=$(NDK_BIN)/aarch64-linux-android21-clang \
	go build $(GO_BUILD_FLAGS) -buildmode=c-shared -o $(PREBUILD_PATH)/$(CURRENT_ARCH)/${LIB_NAME}.so .
	rm $(PREBUILD_PATH)/Android/${LIB_NAME}.h

android-x86: CURRENT_ARCH := x86
android-x86:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=386 \
	CC=$(NDK_BIN)/i686-linux-android21-clang \
	go build $(GO_BUILD_FLAGS) -buildmode=c-shared -o $(PREBUILD_PATH)/$(CURRENT_ARCH)/${LIB_NAME}.so .
	rm $(PREBUILD_PATH)/$(CURRENT_ARCH)/${LIB_NAME}.h

android-x86_64: CURRENT_ARCH := x86_64
android-x86_64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=amd64 \
	CC=$(NDK_BIN)/x86_64-linux-android21-clang \
	go build $(GO_BUILD_FLAGS) -buildmode=c-shared -o $(PREBUILD_PATH)/$(CURRENT_ARCH)/${LIB_NAME}.so .
	rm $(PREBUILD_PATH)/$(CURRENT_ARCH)/${LIB_NAME}.h

android: android-armv7a android-arm64 android-x86 android-x86_64

ios-x86_64-sim:
	GOARCH=amd64 \
	SDK=iphonesimulator \
	./build_ios.sh

ios-arm64-sim:
	GOARCH=arm64 \
	SDK=iphonesimulator \
	./build_ios.sh

ios-arm64:
	GOARCH=arm64 \
	SDK=iphoneos \
	./build_ios.sh

ios: ios-x86_64-sim ios-arm64-sim ios-arm64

macos-arm64:
	GOARCH=arm64 \
	./build_mac.sh

macos-amd64:
	GOARCH=amd64 \
	./build_mac.sh

macos-universal: macos-arm64 macos-amd64
	mkdir -p $(PREBUILD_PATH)/universal/
	lipo \
		-create \
		$(PREBUILD_PATH)/arm64/$(LIB_NAME).a \
		$(PREBUILD_PATH)/x86_64/$(LIB_NAME).a \
		-output $(PREBUILD_PATH)/universal/$(LIB_NAME).a
	rm $(PREBUILD_PATH)/arm64/$(LIB_NAME).a
	rm $(PREBUILD_PATH)/x86_64/$(LIB_NAME).a

macos: macos-arm64 macos-amd64

windows:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	go build $(GO_BUILD_FLAGS) -buildmode=c-shared -o ${PREBUILD_PATH}/AMD64/${LIB_NAME}.dll .
	rm ${PREBUILD_PATH}/AMD64/${LIB_NAME}.h

linux:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	go build $(GO_BUILD_FLAGS) -buildmode=c-shared -o ${PREBUILD_PATH}/x86_64/${LIB_NAME}.so .
	rm ${PREBUILD_PATH}/x86_64/${LIB_NAME}.h
