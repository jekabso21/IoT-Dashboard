import psycopg2
import sys
import os

def main():
    try:
        conn = psycopg2.connect(
            host="65.21.6.209",
            port=5432,
            user="postgres",
            password="Xk9mP2vL8nQ4wR7jT3yH6sF1zB5cN0aD",
            database="mydb"
        )

        print("Connected to PostgreSQL successfully")

        schema_path = os.path.join(os.path.dirname(__file__), "..", "database", "schema.sql")
        with open(schema_path, 'r', encoding='utf-8') as f:
            schema = f.read()

        print("Applying database schema...")

        cur = conn.cursor()
        cur.execute(schema)
        conn.commit()

        print("Database schema applied successfully!")

        cur.execute("SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public'")
        count = cur.fetchone()[0]
        print(f"Verified: {count} tables created in database")

        cur.close()
        conn.close()

    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
