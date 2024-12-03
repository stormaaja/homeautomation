#!/bin/bash

project_dir=$1

if [ -z "$project_dir" ]; then
    echo "Usage: $0 <project_dir>"
    exit 1
fi

cd $project_dir

version_file="version.go"

current_version=$(perl -nle 'print $1 if /Version = "([^"]+)"/' $version_file)

echo "Current version: $current_version"

# Increment the version
build_version=$(echo $current_version | cut -d. -f3)
new_build_version=$(($build_version + 1))
new_version=$(echo $current_version | cut -d. -f1-2).$new_build_version

echo "New version: $new_version"

# Update the version in the version file
perl -pi -e "s/Version = \"$current_version\"/Version = \"$new_version\"/" $version_file