@echo off
echo "build for amd64-windows"

@REM for /f "usebackq" %%t in (`git describe --tags --abbrev=0`) do set "GIT_TAG=%%t"
for /f "usebackq" %%t in (`git describe --tags`) do set "GIT_TAG=%%t"
for /f "usebackq" %%c in (`git rev-parse --abbrev-ref HEAD`) do set "GIT_COMMITID=%%c"
for /f "usebackq" %%d in (`git rev-list -1 HEAD`) do set "GIT_COMMITID=%GIT_COMMITID%+%%d"

for /f "usebackq" %%b in (`powershell Get-Date -Format "yyyy-MMdd-HHmm"`) do set "BUILD_TIME=%%b"
echo "GitTag: %GIT_TAG%"
echo "GitCommitID: %GIT_COMMITID%"
echo "BuildTime: %BUILD_TIME%"

for /d %%i in ("..\..\builds\template*") do (
    rmdir /s /q "%%i"
)

cd ..\..\cmd\template\

set "GOARCH=amd64"
set "GOOS=windows"

go build -ldflags "-X 'template/internal/status.gGitTag=%GIT_TAG%' -X 'template/internal/status.gGitCommitID=%GIT_COMMITID%' -X 'template/internal/status.gBuildTime=%BUILD_TIME%'" -o "..\..\builds\template_%GIT_TAG%_%GIT_COMMITID%_%BUILD_TIME%\bin\template.exe" -buildvcs=false
xcopy /y /e "..\..\configs\*" "..\..\builds\template_%GIT_TAG%_%GIT_COMMITID%_%BUILD_TIME%\"\cfg\
xcopy /y /e "..\..\deploy\scripts\*" "..\..\builds\template_%GIT_TAG%_%GIT_COMMITID%_%BUILD_TIME%\"

xcopy /y /e "..\..\builds\template_%GIT_TAG%_%GIT_COMMITID%_%BUILD_TIME%\bin\*" "..\..\builds\GS\Ver_5.3.16_202311091726_f381c4450_447\bin\"
xcopy /y /e "..\..\builds\template_%GIT_TAG%_%GIT_COMMITID%_%BUILD_TIME%\config\*" "..\..\builds\GS\Ver_5.3.16_202311091726_f381c4450_447\config\"
