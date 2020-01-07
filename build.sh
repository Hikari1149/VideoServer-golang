#! /bin/bash

# Build web UI

cd ~/go/src//videoServer/web
go install
cp ~/go/bin/web ~/go/bin/videoServer_web_ui/web
cp -R ~/go/src/videoServer/template ~/go/bin/videoServer_web_ui