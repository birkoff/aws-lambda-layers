import pyodbc

print("Required libraries pyodbc installed successfully")

def handler(event, context):
    return {"Status": "pyodbc successfully imported"}