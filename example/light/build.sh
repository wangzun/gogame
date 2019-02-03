#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/light
ios-deploy -r -b light.app
