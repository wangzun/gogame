#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/sprite
ios-deploy -r -b sprite.app
