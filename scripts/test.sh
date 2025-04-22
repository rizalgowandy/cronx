#!/usr/bin/env bash

echo -e "\e[32mRunning: \e[33mtest.\e[0m"

echo -e "\e[32mType: \e[33munit test.\e[0m"
command time -f %E gotestsum --format=testname --hide-summary=skipped -- -failfast -count=1 -v -race ./... || exit 1
echo -e "\e[32mUnit test: \e[33msuccess.\e[0m"

echo -e "\e[32mTest: \e[33msuccess.\e[0m"
