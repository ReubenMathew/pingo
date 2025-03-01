#!/bin/bash

package="pingo"

platforms=("windows/amd64" "windows/386" "linux/amd64"  "darwin/amd64")

cd ../src

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name='SodaDB-'$GOOS'-'$GOARCH

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    build_name='../build/'$output_name

    echo "Building" $GOOS $GOARCH "executable"

    env GOOS=$GOOS GOARCH=$GOARCH go build  -v -o $build_name .
    
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

done

if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

printf "\nSodaDB succesfully built!\n"
