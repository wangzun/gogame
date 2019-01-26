#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/load
ios-deploy -r -b load.app
