#!/bin/bash

./build.sh && ./bin/donkeysqquest 2>&1 | tee quest.log
