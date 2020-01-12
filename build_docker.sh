#!/usr/bin/env bash
ZCM="miner"
ZCS="sharder"
read -p "Provide the docker image tag name: " TAG
read -p "Provide the github organisation name[default:-0chaintest]: " organisation
echo "${organisation:-0chaintest}/${ZCM}:$TAG"

REGISTRY_MINER="${organisation:-0chaintest}/${ZCM}"
REGISTRY_SHARDER="${organisation:-0chaintest}/${ZCM}"
ZCHAIN_BUILDBASE="zchain_build_base"
ZCHAIN_BUILDRUN="zchain_run_base"
GIT_COMMIT=$(git rev-list -1 HEAD)
echo $GIT_COMMIT
if [ -n "$TAG" ]; then
echo " $TAG is the tage name provided"
echo -e " Creating 0chain docker the base images..\n"
docker build -f docker.local/build.base/Dockerfile.build_base . -t ${ZCHAIN_BUILDBASE}
docker build -f docker.local/build.base/Dockerfile.run_base   docker.local/build.base -t ${ZCHAIN_BUILDRUN}
 
echo -e "${ZCM}: Docker image build is started.. \n"
sudo docker build --build-arg GIT_COMMIT=$GIT_COMMIT -t ${REGISTRY_MINER}:${TAG} -f docker.local/build.miner/Dockerfile .
sudo docker pull ${REGISTRY_MINER}:latest
sudo docker tag ${REGISTRY_MINER}:latest ${REGISTRY_MINER}:stable_latest
echo "Re-tagging the remote latest tag to stable_latest"
sudo docker push ${REGISTRY_MINER}:stable_latest
sudo docker tag ${REGISTRY_MINER}:${TAG} ${REGISTRY_MINER}:latest
echo "Pushing the new latest tag to dockerhub"
sudo docker push ${REGISTRY_MINER}:latest
echo "Pushing the new tag to dockerhub tagged as ${REGISTRY_MINER}:${TAG}"
sudo docker push ${REGISTRY_MINER}:${TAG}

echo -e "${ZCS}: Docker image build is started.. \n"
sudo docker build --build-arg GIT_COMMIT=$GIT_COMMIT -t ${REGISTRY_SHARDER}:${TAG} -f docker.local/build.sharder/Dockerfile .
sudo docker pull ${REGISTRY_SHARDER}:latest
sudo docker tag ${REGISTRY_SHARDER}:latest ${REGISTRY_SHARDER}:stable_latest
echo "Re-tagging the remote latest tag to stable_latest"
sudo docker push ${REGISTRY_SHARDER}:stable_latest
sudo docker tag ${REGISTRY_SHARDER}:${TAG} ${REGISTRY_SHARDER}:latest
echo "Pushing the new latest tag to dockerhub"
sudo docker push ${REGISTRY_SHARDER}:latest
echo "Pushing the new tag to dockerhub tagged as ${REGISTRY_SHARDER}:${TAG}"
sudo docker push ${REGISTRY_SHARDER}:${TAG}
fi