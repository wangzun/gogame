#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/basic
ios-deploy -r -b basic.app
