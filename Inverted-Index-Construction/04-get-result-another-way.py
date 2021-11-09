import json
import sqlite3
database = './table.db'
conn = sqlite3.connect(database)
cursor = conn.cursor()

keywords = ["Volkswagen", "凡爾賽宮", "洛杉磯市", "東京迪士尼樂園", "史丹福橋球場"]

my_answers = []
for k in keywords:
    row = cursor.execute(
        'SELECT index_id FROM data where articles like ?', ('%'+k+'%',))
    res = row.fetchall()
    my_answers.append([str(r[0]) for r in res])

print(json.dumps(my_answers))
