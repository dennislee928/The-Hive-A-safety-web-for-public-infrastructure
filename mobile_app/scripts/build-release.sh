#!/bin/bash
# Build script for Flutter app releases
# Usage: ./scripts/build-release.sh [android|ios|all]

set -e

BUILD_TYPE=${1:-all}

echo "Building Flutter app for: $BUILD_TYPE"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Navigate to mobile_app directory
cd "$(dirname "$0")/.." || exit 1

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
flutter clean
flutter pub get

# Build Android
if [ "$BUILD_TYPE" = "android" ] || [ "$BUILD_TYPE" = "all" ]; then
  echo -e "${GREEN}Building Android APK (Release)...${NC}"
  flutter build apk --release
  
  echo -e "${GREEN}Building Android AAB (Release)...${NC}"
  flutter build appbundle --release
  
  # Generate checksums
  echo -e "${GREEN}Generating checksums...${NC}"
  cd build/app/outputs/flutter-apk
  sha256sum app-release.apk > app-release.apk.sha256
  cd ../../bundle/release
  sha256sum app-release.aab > app-release.aab.sha256
  
  echo -e "${GREEN}Android builds completed!${NC}"
  echo "APK: build/app/outputs/flutter-apk/app-release.apk"
  echo "AAB: build/app/outputs/bundle/release/app-release.aab"
fi

# Build iOS
if [ "$BUILD_TYPE" = "ios" ] || [ "$BUILD_TYPE" = "all" ]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
    echo -e "${GREEN}Building iOS (Release)...${NC}"
    flutter build ios --release --no-codesign
    
    # Create device-oriented file
    echo -e "${GREEN}Creating iOS device file...${NC}"
    cd build/ios/iphoneos
    if [ -d "Runner.app" ]; then
      zip -r ios-release-device.zip Runner.app
      shasum -a 256 ios-release-device.zip > ios-release-device.zip.sha256 || sha256sum ios-release-device.zip > ios-release-device.zip.sha256
      echo -e "${GREEN}iOS build completed!${NC}"
      echo "Device file: build/ios/iphoneos/ios-release-device.zip"
    fi
  else
    echo -e "${YELLOW}iOS builds can only be done on macOS${NC}"
  fi
fi

echo -e "${GREEN}Build process completed!${NC}"

