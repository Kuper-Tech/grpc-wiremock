#!/bin/bash

set -euo pipefail


BUILD_IMAGE=./dev/build_image.sh
REF=master

TAG="sbermarkettech/grpc-wiremock:dev"
BUILD_ARGS=$(grep -v -E "^#.*|^$" .image-build-args | sed 's@^@--build-arg @g' | xargs)

COMMAND="docker build . $BUILD_ARGS --tag=$TAG"


cat <<EOF > ${BUILD_IMAGE}
#!/bin/bash

set -euo pipefail

echo ${COMMAND}

${COMMAND}
EOF

perl -i -pe 's/ --/ \\\n    --/g' ${BUILD_IMAGE}

chmod +x ${BUILD_IMAGE}


CompileDaemon \
    -command="${BUILD_IMAGE}" \
    -directory=. -build='true' \
    -color -log-prefix=false \
    -pattern='(Dockerfile|.image-build-args|.*\.sh|.*\.conf|.*\.tpl|.*\.list)$'
