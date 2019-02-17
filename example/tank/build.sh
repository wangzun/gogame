#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/tank
ios-deploy -r -b tank.app
