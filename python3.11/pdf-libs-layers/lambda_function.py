from aws_lambda_powertools import Logger, Tracer, Metrics
import fitz
import gspread
from yaml import load, dump
# import pymupdf
import requests

print("Required libraries installed successfully")


def handler(event, context):
    return {"Status": "BaseLayer successfully imported: aws_lambda_powertools, fitz, gspread, yaml, pymupdf, requests"}
