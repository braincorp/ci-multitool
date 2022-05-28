#!/bin/bash

jq 'del(.steps | .[] | .newState,.oldState,.detailedDiff)' $1
