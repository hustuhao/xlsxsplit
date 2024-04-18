#!/bin/bash

# Fixed parameters
PROJECT_NAME="xlsxsplit" # Project name
VERSION=v0.0.2           # Version

RELEASE_NAME="${PROJECT_NAME}"-"${VERSION}" # Release name
OUTPUT_DIR="bin" # Output directory

# Platforms and architectures array
PLATFORMS=("linux" "windows" "darwin")
ARCHITECTURES=("amd64" "arm64")

# Compile the application for different platforms
for PLATFORM in "${PLATFORMS[@]}"
do
    for ARCH in "${ARCHITECTURES[@]}"
    do
        echo "Building for ${PLATFORM} ${ARCH}..."
        GOOS=${PLATFORM} GOARCH=${ARCH} go build -o ${OUTPUT_DIR}/${PROJECT_NAME}_${PLATFORM}_${ARCH}
    done
done

# Create a tar compression file
echo "Creating release package..."
RELEASE_TAR="${RELEASE_NAME}.tar.gz"
tar -czvf "${RELEASE_TAR}" -C ${OUTPUT_DIR} .

# Create a zip compression file
RELEASE_ZIP="${RELEASE_NAME}.zip"
zip -r "${RELEASE_ZIP}" ${OUTPUT_DIR}/*

# Clean up temporary files
echo "Cleaning up..."
#rm -rf ${OUTPUT_DIR}
#rm "${RELEASE_TAR}"
#rm "${RELEASE_ZIP}"

echo "Release process completed!"