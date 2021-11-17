# 建立反向索引表
# v2: 盡可能避免奇怪的標點符號阻礙判斷

from ckiptagger import WS
import json
import sqlite3
from tqdm import tqdm
import re

# Reference: https://github.com/ckiplab/ckiptagger/wiki/Chinese-README

# If you want to download models for ckiptagger
# data_utils.download_data_gdown("./")

ws = WS("./data")

with open('./wiki_2021_10_05_50000.json', 'r', encoding="utf-8") as f:
    data = json.load(f)

database = './table-v2.db'
conn = sqlite3.connect(database)
cursor = conn.cursor()
cursor.execute('DELETE FROM inverted_table')
conn.commit()

batch_size = 200
reg = "[\s\-，.。\:：！!；;\?？（）\(\)\"\'《》〈〉．～—─\=「」『』、”“·／\#\[\]]"

for index in tqdm(range(len(data) // batch_size)):
    input = []
    for i in range(batch_size):
        input.append(re.sub(reg, " ", data[index * batch_size + i]['articles']))
        input.append(re.sub(reg, " ", data[index * batch_size + i]['title']))

    text = ws(input)

    for i in range(batch_size):
        pos = 0
        for key in text[i*2] + text[i*2+1]:
            key = key.strip()
            if key != "":
                cursor.execute('INSERT INTO inverted_table(term, aid, pos) VALUES(?, ?, ?)',
                               (key, data[index * batch_size + i]['id'], pos))
                pos += 1
    conn.commit()

cursor.close()
conn.close()
