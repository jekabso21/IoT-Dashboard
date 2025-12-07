import psycopg2

conn = psycopg2.connect(
    host="65.21.6.209",
    port=5432,
    user="postgres",
    password="Xk9mP2vL8nQ4wR7jT3yH6sF1zB5cN0aD",
    database="mydb"
)

cur = conn.cursor()

sql = """
CREATE OR REPLACE FUNCTION generate_pairing_code()
RETURNS TEXT AS $$
DECLARE
    new_code TEXT;
    code_exists BOOLEAN;
BEGIN
    LOOP
        new_code := LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0');

        SELECT EXISTS(
            SELECT 1 FROM pairing_codes
            WHERE code = new_code
            AND expires_at > NOW()
        ) INTO code_exists;

        EXIT WHEN NOT code_exists;
    END LOOP;

    RETURN new_code;
END;
$$ LANGUAGE plpgsql;
"""

cur.execute(sql)
conn.commit()
print("Function updated successfully")

cur.close()
conn.close()
