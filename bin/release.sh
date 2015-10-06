#!/bin/bash

export GH_ORG=${GH_ORG:-simonleung8}
export GH_REPO=${GH_REPO:-diego-enabler}
export NAME=${NAME:-"Diego Enabler"}
export DESCRIPTION=${DESCRIPTION:-"Turn on Diego runtime support for application"}
export PKG_DIR=${PKG_DIR:=out}

VERSION=$(<VERSION)

if [[ "$(which github-release)X" == "X" ]]; then
  echo "Please install github-release. Read https://github.com/aktau/github-release#readme"
  exit 1
fi


echo "Creating tagged release v${VERSION} of $GH_ORG/$GH_REPO."
read -n1 -r -p "Ok to proceed? (Ctrl-C to cancel)..." key

github-release release \
    --user $GH_ORG --repo $GH_REPO \
    --tag v${VERSION} \
    --name "v${VERSION} - ${NAME}" \
    --description "${DESCRIPTION}"

os_arches=( darwin_amd64 linux_386 linux_amd64 windows_386 windows_amd64 )
for os_arch in "${os_arches[@]}"; do
  asset=$(ls ${PKG_DIR}/${GH_REPO}_${os_arch}* | head -n 1)
  filename="${asset##*/}"

  echo "Uploading $filename..."
  github-release upload \
    --user $GH_ORG --repo $GH_REPO \
    --tag v${VERSION} \
    --name $filename \
    --file ${asset}
done

echo "Release complete: https://github.com/$GH_ORG/$GH_REPO/releases/tag/v$VERSION"
