# 為方便後續處理直接將原 json 檔轉成 table

import sqlite3
import json
with open('./wiki_2021_10_05_50000.json', 'r', encoding="utf-8") as f:
    data = json.load(f)

database = './table.db'
conn = sqlite3.connect(database)
cursor = conn.cursor()
cursor.execute('DELETE FROM data')

for i in range(len(data)):
    cursor.execute('''
        INSERT INTO data(index_id, aid, title, articles)
        VALUES(?, ?, ?, ?)
    ''', (i, data[i]['id'], data[i]['title'], data[i]['articles']))

conn.commit()
cursor.close()
conn.close()
