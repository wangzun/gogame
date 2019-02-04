#!/bin/bash
gomobile build -target=ios github.com/wangzun/gogame/example/material
ios-deploy -r -b material.app
