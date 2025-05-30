# AWS Utilities
from aws_lambda_powertools import Logger, Tracer, Metrics  # aws-lambda-powertools
from aws_xray_sdk.core import xray_recorder, patch_all  # aws-xray-sdk

# HTTP/API
import requests  # requests
from requests.exceptions import RequestException

print("Required libraries installed successfully")


def handler(event, context):
    return {"Status": "BaseLayer successfully imported"}
