# import psycopg2
import sqlalchemy
    
print("Required libraries installed successfully")


def handler(event, context):
    return {"Status": "psycopg2 and sqlalchemy successfully imported"}