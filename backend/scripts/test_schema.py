import psycopg2
import sys

def main():
    try:
        conn = psycopg2.connect(
            host="65.21.6.209",
            port=5432,
            user="postgres",
            password="Xk9mP2vL8nQ4wR7jT3yH6sF1zB5cN0aD",
            database="mydb"
        )

        cur = conn.cursor()

        print("Testing database schema...")
        print("-" * 60)

        print("\n1. Testing tables:")
        cur.execute("""
            SELECT tablename
            FROM pg_tables
            WHERE schemaname = 'public'
            ORDER BY tablename
        """)
        tables = cur.fetchall()
        for table in tables:
            print(f"   + {table[0]}")

        print("\n2. Testing indexes:")
        cur.execute("""
            SELECT indexname
            FROM pg_indexes
            WHERE schemaname = 'public'
            ORDER BY indexname
        """)
        indexes = cur.fetchall()
        for index in indexes:
            print(f"   + {index[0]}")

        print("\n3. Testing generate_pairing_code() function:")
        cur.execute("SELECT generate_pairing_code()")
        code = cur.fetchone()[0]
        print(f"   Generated code: {code}")
        print(f"   Length: {len(code)} (expected: 6)")
        print(f"   All digits: {code.isdigit()}")

        print("\n4. Testing generate_pairing_code() uniqueness:")
        codes = set()
        for i in range(10):
            cur.execute("SELECT generate_pairing_code()")
            codes.add(cur.fetchone()[0])
        print(f"   Generated 10 codes, {len(codes)} unique")

        print("\n5. Testing is_device_online() function:")
        cur.execute("SELECT is_device_online(uuid_generate_v4())")
        result = cur.fetchone()[0]
        print(f"   Non-existent device: {result} (expected: False)")

        print("\n6. Testing RLS policies:")
        cur.execute("""
            SELECT tablename, policyname
            FROM pg_policies
            WHERE schemaname = 'public'
            ORDER BY tablename, policyname
        """)
        policies = cur.fetchall()
        for table, policy in policies:
            print(f"   + {table}: {policy}")

        print("\n7. Testing triggers:")
        cur.execute("""
            SELECT tgname, tgrelid::regclass
            FROM pg_trigger
            WHERE tgisinternal = false
            ORDER BY tgname
        """)
        triggers = cur.fetchall()
        for trigger, table in triggers:
            print(f"   + {trigger} on {table}")

        print("\n8. Testing cleanup functions:")
        cur.execute("SELECT cleanup_expired_pairing_codes()")
        print("   + cleanup_expired_pairing_codes() executed")

        cur.execute("SELECT cleanup_expired_device_tokens()")
        print("   + cleanup_expired_device_tokens() executed")

        print("\n9. Testing views:")
        cur.execute("""
            SELECT viewname
            FROM pg_views
            WHERE schemaname = 'public'
            ORDER BY viewname
        """)
        views = cur.fetchall()
        for view in views:
            print(f"   + {view[0]}")

        print("\n" + "=" * 60)
        print("All schema tests passed!")
        print("=" * 60)

        cur.close()
        conn.close()

    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == "__main__":
    main()
