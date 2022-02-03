#!/bin/bash

echo -e "Current Tags: "
git tag -l

echo -e "**************\n"

NEW_TAG="$1"

if [[ "${NEW_TAG}" == "" ]]; then
    echo -e "Missing the new tag"
    exit 1;
fi

## create and push the new tag
git tag ${NEW_TAG}
git push origin ${NEW_TAG}

# delete a tag
#git tag -d ${TAG_NAME}
#git push origin :refs/tags/${TAG_NAME}