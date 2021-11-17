import json
from ckiptagger import data_utils, WS
import sqlite3
database = './table-v2.db'
conn = sqlite3.connect(database)
cursor = conn.cursor()

inputs = ["美國航空公司", "洛杉磯快艇", "足利義輝", "底特律老虎", "格萊美", "奧斯卡最佳導演獎", "臺北都會區", "塞拉耶佛", "洛杉磯市", "閃電十一人",
          "阿姆斯壯", "石垣島", "Volkswagen", "正規化", "英格蘭足球超級聯賽", "東京迪士尼樂園", "攻殼機動隊", "帝國飯店", "凡爾賽宮", "曾文溪"]
inputs = [i.lower() for i in inputs]

ws = WS("./data")
keywords = ws(inputs, sentence_segmentation=False)
print(keywords)

# 資料庫搜尋時必尋處理連續性問題
# 因為可能出現文字分割後關鍵字前後不連續的問題
# 要注意，sqlite 與一般資料庫不同，本身有區分大小寫，因此在查找時因加入 COLLATE NOCASE

answers = []
for key in keywords:
    like = '%%%s%%' % (''.join(key))
    sql = 'SELECT DISTINCT aid FROM inverted_table WHERE term=? COLLATE NOCASE '
    if len(key) > 1:
        for i in range(len(key) - 1):
            sql += 'INTERSECT SELECT DISTINCT aid FROM inverted_table WHERE term=? COLLATE NOCASE '
    sql = """
        SELECT a.aid FROM (""" + sql + """) as a, data as d
        WHERE (d.articles like ? or d.title like ?) and a.aid = d.aid 
    """

    rows = cursor.execute(sql, tuple(key)+(like, like))
    res = rows.fetchall()

    print(''.join(key), [r[0] for r in res])
    answers.append([str(r[0]) for r in res])


print("--------------------------------")
print(json.dumps(answers))
print("--------------------------------")
