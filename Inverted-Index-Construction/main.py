
from ckiptagger import WS
import json
import sqlite3
from tqdm import tqdm

# 下載模型
# data_utils.download_data_gdown("./")

ws = WS("./data")

with open('./wiki_2021_10_05_50000.json', 'r', encoding="utf-8") as f:
    data = json.load(f)

database = './table.db'
conn = sqlite3.connect(database)
cursor = conn.cursor()

delimiter_set = {",", "，", 
                 ".", "。", 
                 ":", "：", 
                 "！", "!", 
                 "；", ";", 
                 "?", "？", 
                 "（", "）",
                 "(", ")",
                 "\"", "'"}

batch_size = 100
for index in tqdm(range(len(data) // batch_size)):
    input = []
    for i in range(batch_size):
        input.append(data[index * batch_size + i]['articles'])
        input.append(data[index * batch_size + i]['title'])

    text = ws(input, segment_delimiter_set=delimiter_set)

    for i in range(batch_size):
        for key in text[i*2] + text[i*2+1]:
            if key not in delimiter_set:
                cursor.execute('INSERT INTO inverted_table(term, id) VALUES(?, ?)',
                            (key, data[index * batch_size + i]['id']))
    conn.commit()

cursor.close()
conn.close()
