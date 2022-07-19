#!/bin/sh

# 이 프로젝트를 운영하는데에 필요한 모든 프로그램을 검진하는 script입니다.

if [[ ! -f `which docker` ]]; then
  echo "Required 'docker' cli. https://www.docker.com/products/docker-desktop"
fi

if [[ ! -f `which mockgen` ]]; then
  echo "Required 'mockgen' cli. https://github.com/golang/mock"
fi

if [[ ! -f `which jq` ]]; then
  echo "Required 'godepgraph' cli https://stedolan.github.io/jq/"
  echo "jq is utility at terminal environment for manupulating JSON objects"
  echo "[Quick installation]"
  echo "brew install jq"
fi

if [[ ! -f `which direnv` ]]; then
  echo "Required 'direnv' cli. https://github.com/direnv/direnv"
  echo "direnv inject environments at current terminal. commonly useful tool"
fi

if [[ ! -f `which oapi-codegen` ]]; then
  echo "Required 'oapi-codegen' cli."
  echo "go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest"
fi

if [[ ! -f `which mockgen` ]]; then
  echo "Required 'mockgen' cli."
  echo "go install github.com/golang/mock/mockgen"
fi


