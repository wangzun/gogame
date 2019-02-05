#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/skybox
ios-deploy -r -b skybox.app
