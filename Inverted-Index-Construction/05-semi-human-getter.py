# 綜合兩種方法將兩者不同之處之有爭議之處 print 出來
# 以工人智慧判定結果

import json
from ckiptagger import data_utils, WS
import sqlite3
database = './table.db'
conn = sqlite3.connect(database)
cursor = conn.cursor()

inputs = ["apple", "Volkswagen", "凡爾賽宮", "洛杉磯市", "東京迪士尼樂園", "史丹福橋球場"]
#inputs = ["人心惶惶", "縱火犯", "印尼羽毛球", "生長激素", "胰高血糖素"]
inputs = [i.lower() for i in inputs]

ws = WS("./data")
keywords = ws(inputs, sentence_segmentation=False)

print(keywords)

database = './table-v2.db'
conn = sqlite3.connect(database)
cursor = conn.cursor()

# 資料庫搜尋時必尋處理連續性問題
# 因為可能出現文字分割後關鍵字前後不連續的問題
# 另外，由於老師給的答案是 index_id 不是 article_id 所以回傳 index_id
# 要注意，sqlite 與一般資料庫不同，本身有區分大小寫，因此在查找時因加入 COLLATE NOCASE

my_answers = []
for key in keywords:
    like = '%%%s%%' % (''.join(key))
    sql = 'SELECT DISTINCT aid FROM inverted_table WHERE term=? COLLATE NOCASE '
    if len(key) > 1:
        for i in range(len(key) - 1):
            sql += 'INTERSECT SELECT DISTINCT aid FROM inverted_table WHERE term=? COLLATE NOCASE '
    sql = """
        SELECT d.index_id FROM (""" + sql + """) as a, data as d
        WHERE (d.articles like ? or d.title like ?) and a.aid = d.aid 
    """

    rows = cursor.execute(sql, tuple(key)+(like, like))
    res = rows.fetchall()

    print(''.join(key), [r[0] for r in res])
    my_answers.append([r[0] for r in res])

# Way 2
my_answers2 = []
for k in inputs:
    row = cursor.execute(
        'SELECT index_id FROM data where articles like ?', ('%'+k+'%',))
    res = row.fetchall()
    print(k, [r[0] for r in res])
    my_answers2.append([r[0] for r in res])

def cmp(list1, list2):
    list1, list2 = sorted(list1), sorted(list2)
    i, j = 0, 0
    intersection = []
    different = []
    different_type = []
    while i < len(list1) and j < len(list2):
        if list1[i] == list2[j]:
            intersection.append(str(list1[i]))
            i += 1
            j += 1
        elif list1[i] < list2[j]:
            different.append(str(list1[i]))
            different_type.append(0)
            i += 1
        elif list1[i] > list2[j]:
            different.append(str(list2[j]))
            different_type.append(1)
            j += 1
    while i < len(list1):
        different.append(str(list1[i]))
        different_type.append(0)
        i += 1
    while j < len(list2):
        different.append(str(list2[j]))
        different_type.append(1)
        j += 1
    return different, different_type, intersection

def split(word):
    return [char for char in word]

print("--------------------------------")
print("way1:", my_answers)
print("way2:", my_answers2)
print("--------------------------------")

final_answer = []
for i in range(len(my_answers)):
    different, different_type, intersection = cmp(
        my_answers[i], my_answers2[i])
    print("\033[0;32m"+inputs[i]+"\033[0m")

    believe = False

    # 尋問使用者要不要將該比資料加入 intersection
    for idx, d in enumerate(different):
        # 使用者選擇相信，該關鍵字不做檢查
        if believe:
            intersection.append(d)
            continue

        row = cursor.execute(
            'SELECT articles FROM data where index_id = ?', (d,))
        res = row.fetchall()
        print("--------------------------------")
        if different_type[idx] == 0:
            print("該筆資料僅源自於「\033[0;36m斷詞法\033[0m」比較不會錯")
        else:
            print("該筆資料僅源自於「\033[0;36m資料庫爆力搜尋\033[0m」要小心檢查")
        print("關鍵字:", "\033[0;32m"+inputs[i]+"\033[0m")
        mark = split(inputs[i])
        print("正在啟動「斷詞標記輔助系統」...", end='\r')
        # 斷詞輔助標記系統
        words = ws([str(res[0])])[0]
        for word in words:
            word = word.strip(" ,，.。:：！!；;?？（）()\"'《》〈〉．～—─=「」『』、”“·／")
            word = word.lower()
            if word == '': continue
            # 標記整個詞「綠色」
            if word in keywords[i]:
                print("\033[0;32m"+word+"\033[0m ", end=' ')
            # 標記字「藍色」
            else:
                find = split(word)
                for c in find:
                    if c in mark:
                        print("\033[0;36m"+c+"\033[0m", end='')
                    else:
                        print(c, end='')
                print(end=' ')

        print("\n加入該筆資料 b:該關鍵字總是相信 u:該關鍵字總是不相信", d, "[y/n/b/u]", end=' ')
        yes_or_no = input()
        if yes_or_no == 'b':
            believe = True
            intersection.append(d)
        elif yes_or_no == 'u':
            break
        elif yes_or_no == 'y':
            intersection.append(d)
    final_answer.append(intersection)

print("--------------------------------")
print(json.dumps(final_answer))
print("--------------------------------")
