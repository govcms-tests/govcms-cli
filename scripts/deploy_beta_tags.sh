#!/bin/sh
GOVCMS_RELEASE_VERSION=$1
LAGOON_RELEASE_VERSION=$2


####### UPDATE GOVCMS RELEASE VERSION #################################################################################

git clone https://github.com/govCMS/GovCMS.git
cd GovCMS || return
git checkout 3.x-develop
git checkout -b release/3.x/"$GOVCMS_RELEASE_VERSION"

# Replace GovCMS version.
# Note we avoid in-place editing using '-i' as it is not POSIX compliant and does not work on OS X without modifications.
sed -E "s/version: '[0-9]+.[0-9]+.[0-9]+'/version: '$GOVCMS_RELEASE_VERSION'/" govcms.info.yml > /tmp/file.$$ && mv /tmp/file.$$ govcms.info.yml

git add govcms.info.yml
git commit -m "Update release version"
git push --set-upstream origin release/3.x/"$GOVCMS_RELEASE_VERSION"

rm -rf ../GovCMS

######### UPDATE LAGOON RELEASE VERSION ###############################################################################

git clone https://github.com/govCMS/lagoon.git
cd lagoon || return
git checkout 3.x-develop
git checkout -b release/3.x/"$LAGOON_RELEASE_VERSION"

sed -E "s|GOVCMS_PROJECT_VERSION=.*|GOVCMS_PROJECT_VERSION=$GOVCMS_RELEASE_VERSION|g" .env.default > /tmp/file.$$ && mv /tmp/file.$$ .env.default

git add .env.default
git commit -m "Update release version"
git push --set-upstream origin release/3.x/"$LAGOON_RELEASE_VERSION"

rm -rf ../lagoon

####### PUSH BETA TAGS TO DOCKERHUB ##################################################################################
