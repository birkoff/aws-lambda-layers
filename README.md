# lambda-layer
A Lambda Layer CLI Tool to Build and Deploy AWS Lambda Layers

# Requirements:
- Docker

# Build and Deploy
```
go build -o lambda-layer main.go
```

# Test
```
go run main.go -name langchain-base-layer
```

# Installation
Download the Lambda Layer Binary
https://github.com/junctionnet/tools/blob/a485ad5612e88868900ec4cedbbe8ed8c9ae3f54/lambda-layer-tool/lambda-layer 

Add the path where you downloaded the CLI Tool on BashRC or ZShrc
`code ~/.zshrc`


Example:
`export PATH="/Users/yomerengues/tools/lambda-layer:$PATH"`



# Usage
Requirements file with the following name convention <lambda layer name>-requirements.txt

```
-> % cat langchain-base-layer-requirements.txt
langchain==0.2.12
langchain-community==0.2.10
langchain-core==0.2.27
langchain-openai==0.1.20
langchain-text-splitters==0.2.2
langsmith==0.1.96
openai==1.37.2
urllib3==1.26.19
typing-inspect==0.9.0
tiktoken==0.7.0
```

```
lambda-layer build <lambda layer name> <runtime>
```

```
-> % lambda-layer build langchain-base-layer python3.10


Building Lambda layer with the following details:
Layer Name: langchain-base-layer
Compatible Runtimes: python3.10
Requirements File: /Users/hector/code/junctionnet/github/apis/entrywriter-api/infrastructure/layers/langchain-base-layer-requirements.txt
################################################################################
Should I proceed building this Layer...is Docker running? (yes/no): yes
Running: docker pull public.ecr.aws/lambda/python:3.10
Running: docker ps -a --format '{{.Names}}'
Container amzn already exists. Removing it...
Running: docker rm -f amzn
Running: docker run --name amzn -d -t --rm public.ecr.aws/lambda/python:3.10 /bin/bash
Copying requirements file from /Users/hector/code/junctionnet/github/apis/entrywriter-api/infrastructure/layers/langchain-base-layer-requirements.txt to container...
Running: docker cp /Users/hector/code/junctionnet/github/apis/entrywriter-api/infrastructure/layers/langchain-base-layer-requirements.txt amzn:/root/requirements.txt
Executing inside container: yum update -y
Running: docker exec amzn sh -c "yum update -y"
Executing inside container: yum install -y zip
Running: docker exec amzn sh -c "yum install -y zip"
Executing inside container: python3.10 -m ensurepip
Running: docker exec amzn sh -c "python3.10 -m ensurepip"
WARNING: Running pip as the 'root' user can result in broken permissions and conflicting behaviour with the system package manager. It is recommended to use a virtual environment instead: https://pip.pypa.io/warnings/venv
Executing inside container: python3.10 -m pip install --upgrade pip
Running: docker exec amzn sh -c "python3.10 -m pip install --upgrade pip"
WARNING: Running pip as the 'root' user can result in broken permissions and conflicting behaviour with the system package manager. It is recommended to use a virtual environment instead: https://pip.pypa.io/warnings/venv
Executing inside container: python3.10 -m pip install -r /root/requirements.txt -t /root/package/python
Running: docker exec amzn sh -c "python3.10 -m pip install -r /root/requirements.txt -t /root/package/python"
WARNING: Running pip as the 'root' user can result in broken permissions and conflicting behaviour with the system package manager, possibly rendering your system unusable.It is recommended to use a virtual environment instead: https://pip.pypa.io/warnings/venv. Use the --root-user-action option if you know what you are doing and want to suppress this warning.
Executing inside container: find /root/package/ -name '*.pyc' -delete
Running: docker exec amzn sh -c "find /root/package/ -name '*.pyc' -delete"
Executing inside container: rm -rf /root/package/urllib3* /root/package/six* /root/package/botocore* /root/package/idna* /root/package/tomli* /root/package/jmespath* /root/package/charset_normalizer*
Running: docker exec amzn sh -c "rm -rf /root/package/urllib3* /root/package/six* /root/package/botocore* /root/package/idna* /root/package/tomli* /root/package/jmespath* /root/package/charset_normalizer*"
Executing inside container: find /root/package/ -name '*.dist-info' -exec rm -rf {} +
Running: docker exec amzn sh -c "find /root/package/ -name '*.dist-info' -exec rm -rf {} +"
Executing inside container: find /root/package/ -type d -name '__pycache__' -exec rm -r {} +
Running: docker exec amzn sh -c "find /root/package/ -type d -name '__pycache__' -exec rm -r {} +"
Executing inside container: cd /root/package && zip -r9qv ../langchain-base-layer.zip *
Running: docker exec amzn sh -c "cd /root/package && zip -r9qv ../langchain-base-layer.zip *"
Copying the resulting ZIP file to /Users/hector/code/junctionnet/github/apis/entrywriter-api/infrastructure/layers/langchain-base-layer.zip...
Running: docker cp amzn:/root/langchain-base-layer.zip /Users/hector/code/junctionnet/github/apis/entrywriter-api/infrastructure/layers/langchain-base-layer.zip
Should I Publish the layer to AWS Lambda? (yes/no):
```