#!/bin/sh

# Build the project
npm run build

# Create the static directory
rm -rf ../static/
mkdir ../static/
cp build/* ../static/
cp -r build/static/css ../static/
cp -r build/static/js ../static/
cp -r build/static/images ../static/

# cleanup
rm -rf ./build