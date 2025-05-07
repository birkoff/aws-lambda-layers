import os
import pyodbc

def handler(event, context=None):
    conn_str = (
        "DRIVER=/opt/microsoft/msodbcsql18/lib64/libmsodbcsql-18.5.so.1.1;"
        f"SERVER={os.getenv('DB_HOST', 'db')};"
        f"DATABASE={os.getenv('DB_NAME', 'master')};"
        f"UID={os.getenv('DB_USER', 'sa')};"
        f"PWD={os.getenv('DB_PASSWORD')};"
        "TrustServerCertificate=yes"
    )
    
    try:
        with pyodbc.connect(conn_str) as conn:
            with conn.cursor() as cursor:
                cursor.execute("SELECT @@VERSION")
                return str(cursor.fetchone()[0])
    except Exception as e:
        return str(e)