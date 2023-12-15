echo "build for amd64-windows"

GIT_TAG=$(git describe --tags)
GIT_COMMITID="$(git rev-parse --abbrev-ref HEAD)+$(git rev-list -1 HEAD)"
BUILD_TIME=$(date "+%Y-%m%d-%H%M")
echo "GitTag: ${GIT_TAG}"
echo "GitCommitID: ${GIT_COMMITID}"
echo "BuildTime: ${BUILD_TIME}"

rm -r ../../builds/template*

cd ../../cmd/template/

export GOARCH=amd64
export GOOS=windows

go build -ldflags "-X 'template/internal/status.gGitTag=${GIT_TAG}' -X 'template/internal/status.gGitCommitID=${GIT_COMMITID}' -X 'template/internal/status.gBuildTime=${BUILD_TIME}'" -o ../../builds/template_${GIT_TAG}_${GIT_COMMITID}_${BUILD_TIME}/bin/template.exe -buildvcs=false
cp -f -r ../../configs/* ../../builds/template_${GIT_TAG}_${GIT_COMMITID}_${BUILD_TIME}/cfg/
cp -f -r ../../deploy/scripts/* ../../builds/template_${GIT_TAG}_${GIT_COMMITID}_${BUILD_TIME}/