#!/usr/bin/env bash

echo -e "\e[32mRunning:\e[33m setup.\e[0m\n"

echo -e "\e[32mInstalling:\e[33m air for live reload.\e[0m"
command -v air 2>/dev/null || go install github.com/air-verse/air@v1.61.7
echo ""

echo -e "\e[32mInstalling:\e[33m mockgen for mock generator.\e[0m"
command -v mockgen 2>/dev/null || go install -v github.com/golang/mock/mockgen@v1.6.0 # v1.6.0
echo ""

echo -e "\e[32mInstalling:\e[33m golangci-lint for linter.\e[0m"
command -v golangci-lint 2>/dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.2
echo ""

echo -e "\e[32mInstalling:\e[33m gomodifytags for generating tags.\e[0m"
command -v gomodifytags 2>/dev/null || go install -v github.com/fatih/gomodifytags@v1.17.0
echo ""

echo -e "\e[32mInstalling:\e[33m gotestsum.\e[0m"
command -v gotestsum 2>/dev/null || go install -v gotest.tools/gotestsum@v1.12.1
echo ""

echo -e "\e[32mSetup:\e[33m pre-commit hook.\e[0m"
file=.git/hooks/pre-commit
cp scripts/pre-commit.sh $file
chmod +x $file
test -f $file && echo "$file exists."
echo ""

echo -e "\e[32mSetup:\e[33m success.\e[0m"
