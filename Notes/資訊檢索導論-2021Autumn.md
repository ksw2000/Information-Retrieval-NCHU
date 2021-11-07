---
tags: 大四
---

# 資訊檢索導論 (Introduction to Information Retrieval) Autumn 2021

[![](https://img.shields.io/badge/dynamic/json?style=flat-square&color=green&label=%E7%B8%BD%E8%A7%80%E7%9C%8B%E6%95%B8&query=%24.viewcount&suffix=%E4%BA%BA%E6%AC%A1&url=https%3A%2F%2Fhackmd.io%2F01qNrz7xSF-MHv2YjlVotQ%2Finfo)](https://hackmd.io/01qNrz7xSF-MHv2YjlVotQ/info)

![](https://i.imgur.com/EhXkW1A.png)

資訊檢索：如何在一群資料中找到你要找的東西

[TOC]

## 成績計算

> Assignments (20%)
> Midterm Exam (20%)
> Final Exam I (30%)
> Final Exam II(30%)
>
> or
>
> Assignments (20%)
> Midterm Exam (30%)
> Final Exam I (30%)
> Final Exam II(20%)

## Overview

### 為什麼 Google 那麼快？

1. Google 有爬蟲 (crawler)
2. Tokenizer (適當切割語義)
    > Tokenization: Cut character sequence into word tokens
    > 英文斷詞很容易，但中文就不容易
    > Normalization: USA 與 U.S.A 是一樣的
    > Stemming: authorize, authorization
    > stop words: the, a, to, of, 了, 呃, 那個呃
    > python: NLTK 套件
3. Index (建立索引)
    > ![](https://i.imgur.com/rmMUSPr.png  =250x)


Openfind(台灣蕃薯藤) 與 Google(美國) 皆以 Page Rank 做為搜尋引擎的演算法，Openfind 失敗的其中一個原因可能是因為 tokenization， 中文的斷詞比英文的斷詞來的困難太多

## Boolean Queries

Boolean retrieval model is able to ask a query that is a boolean expression.

Boolean retrieval 是用集合來做計算，集合是無序的，無法分辨「我愛你」「你愛我」的差別

Example: 
```
Which news contain the words 蔡英文 AND 郭台銘 but NOT 韓國瑜
```

Question: Too slow (無駄無駄無駄無駄無駄無駄)


### Term-document incidence matrices

|document\key| key1  | key2 | key3 |
|------------|--|--|--|
| doc1 | 1 | 0 | 0 |
| doc2 | 0 | 0 | 1 |
| doc3 | 1 | 1 | 0 |
| doc4 | 1 | 1 | 1 |

1 if doc contains keywords, otherwise 0


Question: Waste a lot of space, 可以用 adjacency list 來表示

### Inverted index 反向索引

以關鍵字當 key 值，一串文件 list 當作 value

> 如果不知道 Inverted Index ，請說你的IR是曾學文教的
> \-\- 范耀中, 2021

![](https://i.imgur.com/i1I4Wj3.png)

註: 該方法屬於 pre processing, 圖中是老師在問是 pre- 還是 post- processing 

公平性問題：這樣子比較稀有的 keyword 查起來會比較快 (茶

<!--![](https://i.imgur.com/L6N9CZV.png)-->
<!--![](https://i.imgur.com/G3xL37E.png)-->


### Query Processing with Inverted Index

比較兩個 posting list (linked-list) 取交集: O(x+y)

![](https://i.imgur.com/lwbhojj.png =500x)

Question: Too slow!! (如果太多關鍵字會非常久)
Solution: 偷看後面的 (**skip pointer**)

#### Skip pointer

![](https://i.imgur.com/Uk8VDd5.png =500x)

那要跳幾個呢？J個94商業機密了(Google只說長度 $L$ 可以設為 $L^{0.5}$)

+ 一次 skip 太多(只要跳比較少次)：跳太多的話很可能要回頭找
+ 一次 skip 太少(要跳比較多次)：跳太少就跟沒跳一樣

#### And Not / Or Not

> **思考一個例子**
> 
> AND NOT 與 OR NOT 的差別
> 
> 蔡英文 AND NOT 韓國瑜 → O(x+y)
> 蔡英文 OR NOT 韓國瑜 → O(n) 因為 "NOT 韓國瑜" 太大
>
> x: 含有蔡英文的資料
> y: 含有韓國瑜的資料
> n: 所有資料
>
> 舉例：台科大圖書館僅能使用 and not 不能使用 or not

#### Query optimization

如何有效率地 Merge 兩個搜尋結果？

![](https://i.imgur.com/8C3wD0H.png)

先處理較短的 list 比如上圖中若要處理三個的交集，則先尋找 Calpurnia 和 Brutus 的交集，再交和 Caesar 的交集



## Phrase Queries

### Biword indexes 

> 將片語也當作是單字來處理建立反向索引，也就是會需要 $C^n_2\times2!$ 個(有序) key
> 
> Question: 要建立太多 key，而且如果片語是三個字或更多字呢？

### Positional indexes

除了記得單字他在哪篇文，另外還記得他在文章中的哪個位置

![](https://i.imgur.com/NIGjdM0.png)

to 在 2 號文章的 1, 17, 74, 222, 551 的位置

<!--![](https://i.imgur.com/eitBBWq.png)-->


### Combination schemes

常用的片語照樣建成 biword indexes，但仍然使用 positional indexes

<!--9/29-->

## Tolerant Retrieval

如果打錯字了，Google 要怎麼看懂你打錯字呢？

+ 常見的錯包括 Typo, context error。
+ 在做 spelling correction 前比需先把現有的字詞編成字典

### Dictionary data structure

1. Hash table
    + 搜尋複雜度 O(1)
    + 無法解決 wildcard
2. Tree
    + 搜尋複雜度 O(logN)
    + 可以解決簡單的 wildcard 比如 `Mom*`, `*Mom`, `Mom*mom` 這種的，但是 `*Mom*` 這類的則無法處理
3. K-Gram Index
    + 可以解決 wildcard queries

### k-Gram Index

舉個例子：[water blue new world](https://www.youtube.com/watch?v=obW81vi3COw) 拆成
> `$w` `wa` `at` `te` `er` `r$`
> `$b` `bl` `lu` `ue` `e$`
> `$n` `ne` `ew` `w$`
> `$w` `wo` `or` `rl` `ld` `d$`

透過這個手法建立索引

`$w` -> `we`, `wait`, `word`
`wa` -> `water`, `await`, `taiwan`, `awaken`

以上示範為 k=2 又稱 bigram index，當然也可以用 k=3

> `$wa` `wat` `ate` `ter` `er$`
> `$bl` `blu` `lue` `ue$`
> `$ne` `new` `ew$`
> `$wo` `wor` `orl` `rld` `ld$`

### Spelling correction ─ Typo


![](https://i.imgur.com/lt6GxVR.png)


#### Edit distance (編輯距離)

將使用者打的字與字典的字做計算，計算最少新增、刪除、取代的次數

+ d(dof, dog) = 1
+ d(cat, act) = 2
+ d(dog, cat) = 3

> 複雜度 $O(m\times n \times M)$
> m: 單字A長度, n: 單字B長度, M: 全部字典有多少字

> 延伸 → **weighted edit distance**
> (m, n) < (m, q)
> 因為 q 在鍵盤上離 m 比較遠，比較不太像打錯

#### K-Gram overlap

把單字依成k個一組再比對

> **november**
> nov, ove, vem, **emb**, **mbe**, **ber**.
> **december**
> dec, ece, cem, **emb**, **mbe**, **ber**.
> 
> december 和 november 相似度為 3

> Query **lord**
> ![](https://i.imgur.com/24bJkXM.png)

> 2: lord, border
> 1: alone, sloth, morbid, card,...


> 複雜度 $O((m+n) \times M)$
> m: 單字A長度, n: 單字B長度, M: 全部字典有多少字

> 但是有些字太長了要怎麼比？
> 
> **延伸** → Jaccard
> K-Gram 算出的距離為 0 到無限大，希望可以 mapping 到 0 到 1之間，因此可以計算 Jaccard coefficient
> $$
> \frac{A交集B}{A聯集B}
> $$

### Spelling correction ─ Context error

![](https://i.imgur.com/fnOgTtE.png)

> 爆力解法：
> 把每個字的相近單字都納入考慮
> 把 form 的相近字 taipei 的相近字 to 的相近字 tokyo 的相近字做排列組合尋找
> 
> (無駄無駄無駄)
> 
> A better approach: Break phrase query into biwords. Look for biwords that need only one term corrected. 


### Dictionary Compression

#### 可能有的問題

1. 固定字串長度太浪費
2. Postings list的ID欄位佔太大

#### sol

固定字串長度太浪費：

+ 把字串全部串在一起只存指標

Postings list的ID欄位佔太大

+ 只存第一個，後面的存差值

ex:

> 33, 47, 154, 159, 202 ->
> 33, 14, 107, 5, 43

## Rank Retrieval

Boolean search 有個問題 → Difficult to rank output.

### TFIDF

> 如果修過 IR 不知道 TFIDF 那就是 co-domain 的問題囉
> \-\- 范耀中, 2021

$$
w_{t,d}=(1+\log{tf}_{td})(log\frac{N}{df_t})
$$

**TF: Term Frequency**

$$
w_{t,d} = \begin{cases}
   1+\log tf_{t,d} &\text{if }\ \  tf_{t,d}>0 \\
   0 &\text{otherwise } 
\end{cases}
$$

$tf_{t,d}$ 文字 t 在文章 d 中出現的次數
$w_{t,d}$ 文字 t 在文章 d 中的重要性

**IDF: Inverse Document Frequency**

但是僅有 TF 這個方法有個問題，比如 `the`, `on`, `my`, `I`, `your`, `to` 所占的頻率還是比較高。

Document frequency: 偵對某一個字在文章出現的比率，比如某個字常常在各篇文章出現，就代表這個字並沒有很重要；反之，如果有個字比較少在各篇文章出現，但某篇文章出現特別多，就代表這篇文章是重要的

$$
idf_t=\log\frac{N}{df_t} \\
\text{N: 文章數量}
$$

![](https://i.imgur.com/WMUpF9e.png)

**Example**

+ D1: “Shipment of gold damaged in a fire”
+ D2: “Delivery of silver arrived in a silver truck”
+ D3: “Shipment of gold arrived in a truck”

> 計算 TF 和 IDF

![](https://i.imgur.com/iifNgRQ.png)

$$
w_{silver, D2}=(tf)(idf)=(1+log2)(0.4771) = 0.62
$$

### OKAPI-BM25

OKAPI-BM25 比起 TFIDF 來得更好

#### 改用除法處理 term 的飽和度

捨棄log的方法，因為除法比較快，而且值會值會介於 0~1 之間

$$
w_{t,f} = \frac{tf}{tf+k}
$$

![](https://i.imgur.com/Kim7Gij.png)
> 藍線為 $y= 1+log\, x$
> 綠線為 $y=\frac{x}{x+2}$

k 是自己設的參數，通常在2到3之間

#### 考慮文章長度

$$
w_{t,f} = \frac{tf}{tf+k\times\frac{dl}{adl}}
$$

> dl: document length
> adl: average document length
>
> 也就是說當文章太長時會使最後 w 數值下降

#### 考慮文章長度(允許給定參數調整其影響程度)

$$
w_{t,f} = \frac{tf}{tf+k\times(1-b+b\times\frac{dl}{adl})}
$$

## Vector-space Model

> A collection of n documents can be represented in the vector space model by a term-document matrix.

| Document \ Term | T<sub>1</sub> | T<sub>2</sub> | T<sub>3</sub> | ... | T<sub>t</sub>|
|-----------------|----|----|----|---|---|
| D<sub>1</sub> | w<sub>11</sub> | w<sub>21</sub> | w<sub>31</sub>| ... | w<sub>t1</sub>|
| D<sub>2</sub> | w<sub>12</sub> | w<sub>22</sub> | w<sub>32</sub>|...| w<sub>t2</sub> |


![](https://i.imgur.com/IoBUGYI.png =250x)

### Method 1: Euclidean Distance 歐式距離 (垃圾)

※ TFIDF 沒有一定要在 (0, 1) 之間

![](https://i.imgur.com/6VRNyYE.png =250x)

$q$ 理論上更接近 $d_2$ 但用歐式距離算起來會比較靠近 $d_1$, $d_3$

### Method 2: Cosine Similarity Measure

比起用距離，用夾角計算相似度會更接近實際情況 (高二數學)

![](https://i.imgur.com/C7vU1CM.png)

<!--![](https://i.imgur.com/VnLZfam.png)-->

效率問題，cosine similarity 計算太複雜，若要比較兩篇文章是否相似

#### TAAT (Term At A Time)

Scores for all docs computed concurrently, one query term at a time.

先看第一個 term，算出每個文件的 partial score，再看第二個 term 再慢慢加起來...

+ 優點是能減少硬碟的讀取
+ 缺點是很耗記憶體，因為要儲存所有候選文件的 partial score

#### DAAT (Document At A Time)

Total score for each doc (include all query terms) computed, before proceeding to the next.

先看第一個文件所存在的 term

+ 優點是能節省記憶體
+ 缺點是必需把整個 inverse index 的 list 讀取出來

### Optimization on top-K

#### Impact-ordered postings for TAAT

+ 用 $wf_{t,d}$ 倒序
+ 直到 $wf_{t,d}$ 小於某個threshold

TOP-1 範例

<img src="https://i.imgur.com/E49DhHR.gif" width="80%">

#### for DAAT: WAND method

+ 按`docID`正排
+ finger: 指向目前的`docID`
+ Upper Bound(UB):還沒看過的裡面$w_t$最小的`doc`

詳細流程可見[Day 15: 神奇的法杖 - 提高效率的WAND演算法](https://ithelp.ithome.com.tw/articles/10215669)，有很詳盡的說明

### Non Safe Ranking

#### High idf query terms only

Basic cosine computation algorithm only
considers docs containing at least one query
term. e.g. 移除太常出現的字 "the", "a"

#### Docs containing many query terms

Any doc with at least one query term is a candidate for the top K output list.For multi term queries, only compute scores
for docs containing several of the query terms.

+ Say, at least 3 out of 4 (4個關鍵字至少要有3個有 match 到，以 [Bloom filter](https//zh.wikipedia.org/wiki/布隆过滤器) 實作，bloom filter 雖然有可能會挑錯但效能和記憶體都相當快)

#### Cluster Pruning

利用 K-means 先分群，再挑最接近結果的群來抓

![](https://i.imgur.com/9YGSYVf.png =400x)

## Link Analysis Model

### Page Rank

將網路上的網頁畫成一個 Graph
1. 看 out-degree? 無駄: 可以自己指向其他很多人
2. 看 in-degree? 無駄: 可以有其他垃圾指向他

> 利用其他人的 page rank 來決定自己的 page rank
> 另外，如果指向越多人，那每個被指的人分數也會被均攤掉
>
> 也就是說，自己的 page rank 由別人決定
>
> ![](https://i.imgur.com/xWnZ4er.png =300x)
>
> 我們可以把上述的式子轉換成矩陣形勢，跟高二學的轉移矩陣是相同道理
> $$
> \begin{bmatrix}
> \frac{1}{2} & \frac{1}{2} & 0 \\
> \frac{1}{2} & 0 & 1 \\
> 0 & \frac{1}{2} & 0 \\
> \end{bmatrix}
> \begin{bmatrix}y\\a\\m\end{bmatrix}
> $$

#### Random teleports

Google 的爬蟲爬一爬可能就爬不出去了，因此 Google 設計一套機制，每 $\beta$ 的機率會延著連結爬，$1-\beta$ 的機率隨機跳到他網站，修正過後的 graph 會長這個樣子 (假設 $\beta = 0.8$ )

![](https://i.imgur.com/yE8F0gD.png =300x)


### HITS

HITS 與 Page Rank 想法類似，但 HITS 會將網頁分成兩個分數，一個是 authority score，一個是 hub score

一個好的 hub 會指向好的 authority
一個好的 authority 會由很多好的 hub 指向

![](https://i.imgur.com/2tz5Fev.png =400x)

**缺點**

SEO 人士可以先偽造分數高的 Hub，再把高的 Hub 指向自己的垃圾網頁

> Example:
> ![](https://i.imgur.com/rKmwmAy.png =200x)
>
> 跟 Page Rank 相比，要反過來看，因為 hub score 是由指向的網站來評分
> 
> 也就是說雖然 y 指向 y, a, m 但實際上 y 的 hub score 是來自 y, a, m 三者
>
> | | y | a | m |
> |-|---|---|---|
> |y| 1 | 1 | 1 |
> |a| 1 | 0 | 1 |
> |m| 0 | 1 | 0 |
> 
> $$
> h = \lambda Aa \\
> a =\mu A^Th
> $$
> h: hub score
> a: authority score
> A: 矩陣
> $\lambda , \mu$: scale factor (因為算出來會有單一項超過 1，如果算出來是(1, 2, 3)就把 scale factor 設成 1/3，把結果mapping到(1/3, 2/3 , 1))

## Evaluate

### Evaluate unranked results

Precision = P(relevant|retrieved) (回傳的答案裡的確有相關的)
Recall = P(retrieved|relevant) (所有有相關的答案裡確實有被挑出來的)

|           | relevant | non relevant |
|-----------|:-----------:|:-------------:|
| retrieved |    tp     |     fp      |
| not retrieved | fn     |  tn  |

Precision = tp/(tp+fp)
Recall = tp/(tp+fn)
Accuracy = (tp+tn)/(tp+fp+fn+tn)

Accuracy 在資訊檢索系統上沒有意義，因為實際上 tn 非常大，只要都不怎麼回傳，就會有很大的 Accuracy

Precision 只要刻意減少回傳的量，也能達到很高的 Precision

Recall 只要將所有文章都查找回來，就會有很高的 Recall

#### F-measure

$$
F = \frac{1}{\frac{1}{R}+\frac{1}{P}}
$$

#### E-measure

$$
E = \frac{1+\beta^2}{\frac{\beta^2}{R}+\frac{1}{P}}
$$

$\beta$ 用來控制 trade-off 
+ $\beta>1$: weight recall more
+ $\beta<1$: weight precision more

### Evaluate ranked results

#### precision-recall curve

1. 計算當回傳 n 筆資料時所造成的 recall 與 precision![](https://i.imgur.com/FUeEwnE.png =400x)
2. 再將這些點描繪出來 ![](https://i.imgur.com/XckkiWJ.png =400x)



#### Precision@K

大家說好只回傳 k 筆，再去比。如下圖紅色表示錯誤，綠色表示正確

![](https://i.imgur.com/r5fkXZu.png =250x)


#### R-Precision

k 不好決定於就將所有 relevant 的數量當做 k，比如 relevant 有 6 篇就變成 precision@6

#### Mean average precision (MAP)

![](https://i.imgur.com/qMoVApv.png)

將有抓到 relevant document 時的 precision 取平均

> + ranking #1: avg(0.1, 0.67, 0.5, 0.44, 0.5) = 0.62
> + ranking #2: avg(0.5, 0.4, 0.43) = 0.44
> 
> MAP = mean average precision = avg(0.62, 0.44)

#### Mean Reciprocal Rank (MRR)

+ 把第一個有相關的位置定出來, K (比如在第2筆才是第一個有相關的文章那麼 K=2)
+ Reciprocal Rank = 1 / K
+ MRR = The mean RR across multiple queries

#### Discounted Cumulative Gain (DCG)
`lg: log 以 2 為底`
+ 每個 relavent 都會有相應的分數 $r_i$(相關度)
+ 越前面的應該要佔較重比例，後面則越少

$$
DCG_p=r_1+\sum^p_{i=2}\frac{r_i}{lg\,i}
$$ 

#### Normalized DCG (NDCG)

> **For example**:
> ![](https://i.imgur.com/B4tPbkA.png)
> 
> **Ground Truth**:
> $$
> DCG = 2 + \frac{2}{lg2}+ \frac{1}{lg3}+ \frac{0}{lg4} = 4.6309
> $$
> **Ranking function 2**:
> $$
> DCG = 2 + \frac{1}{lg2}+ \frac{2}{lg3}+ \frac{0}{lg4} = 4.2619
> $$

+ 使用 Ground Truth (之前提到的標準)的 DCG 當作 1
+ 若 $DCG_{GT}=4.6309$ , $DCG_{RF} = 4.2619$ 
+ 則 $NGDC_{RF} = \frac{4.2619}{4.6309}=0.9203$


<!-- 
這麼晚了還這麼多人在看哈哈
(留言不登入，此風不可長)
就是不登入 特地登出來留言
 -->


#### 期中考題露出
 ![](https://i.imgur.com/qreDUXq.png)
