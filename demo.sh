#!/bin/bash

BK_CMD=bk
if [ -f "./bin/bk" ]; then
    BK_CMD=./bin/bk
fi

$BK_CMD -version

$BK_CMD post -date 20200501 \
    -left 1110/500000/会社設立 -right 3100/500000/会社設立

$BK_CMD post -date 20200502 \
    -left 1110/1000000/設備導入資金 -right 2200/1000000/設備導入資金

$BK_CMD post -date 20200503 \
    -left 7300/50000/事務用品 -right 1110/50000/事務用品

$BK_CMD post -date 20200503 \
    -left 1211/500000/パソコン -right 1110/500000/パソコン

$BK_CMD post -date 20200505 \
    -left 5200/100000/おもちゃ仕入 -right 1110/100000/おもちゃ仕入

$BK_CMD post -date 20200507 \
    -left 1110/200000/おもちゃ販売 -right 4100/200000/おもちゃ販売

$BK_CMD post -date 20200510 \
    -left 1110/1000000/運転資金 -right 2101/1000000/運転資金

$BK_CMD post -date 20200511 \
    -left 5200/2000000/おもちゃ仕入 -right 2100/2000000/おもちゃ仕入

$BK_CMD post -date 20200512 \
    -left 1120/4000000/おもちゃ販売 -right 4100/4000000/おもちゃ販売

$BK_CMD post -date 20200515 \
    -left 2100/2000000/買掛金清算 -right 1110/2000000/買掛金清算

$BK_CMD post -date 20200516 \
    -left 1110/3000000/売掛金回収 -right 1120/3000000/売掛金回収

$BK_CMD post -date 20200520 \
    -left 7200/300000/事務員A給与 -right 1110/290000/給与 \
                                -right 2103/10000/源泉所得税

$BK_CMD post -date 20200521 \
    -left 2101/1000000/返済 -right 1110/1100000/返済 \
    -left 8200/100000/支払利息

$BK_CMD post -date 20200522 \
    -left 7300/200000/旅費交通費 -right 1110/200000/旅費交通費

$BK_CMD post -date 20200531 \
    -left 1130/100000/繰越商品 -right 5300/100000/繰越商品

$BK_CMD post -date 20200502 \
    -left 7300/100000/パソコン減価償却 -right 1211/100000/パソコン減価償却

$BK_CMD post -date 20200531 \
    -left 9000/450000/法人税 -right 2102/450000/法人税

