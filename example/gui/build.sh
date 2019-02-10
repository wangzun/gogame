#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/gui
ios-deploy -r -b gui.app
