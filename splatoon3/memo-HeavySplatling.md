# バレル

<img src="./images/HeavySplatling.png" width="256px" alt="バレル" />&nbsp;<img src="./images/is.png" width="128px" alt="is" />&nbsp;<img src="./images/JustInCaseSomethingLikeThisHappens.png" width="192px" alt="こんなこともあろうかと" />

## 立ち回り

- バレルの真髄は「事前準備」
- ❌️ 一回のチャンスで勝つ ⭕️ 少しずつアドバンテージを増やす
- 相手の土俵で戦わずこちらの有利を押し付ける

### キホンのif-then

| 状況 | if-then |
| -- | -- |
| 原則 | <ol><li>塗る</li><li>[if]危険を感じたら → 引く</li><li>[if]ミカタがSP使ったら → 合わせる</li><li>[if]2人有利になったら → 前に出る</li></ol> |
| 確認 | <ul><li>[if]情報手に入れたら → わかったことを口に出す<ul><li>e.g. ミカタ落ちた</li><li>e.g. 壁裏ブラスターいる</li><li>e.g. SP溜まったけどミカタまだ</li></ul></li></ul> |
| 索敵 | <ul><li>[if]対面していない時は → 注意するポイントを順番に見る<ul><li>❌️ 全部を見る。※ それは無理</li><li>⭕️ 先に決めた注意すべき場所を順番に見ていく。※ 危ない順に目視確認をしていく</li></ul></ul> |
| 射撃 | <ul><li>[if]チャージする前に → 周りを見る</li><li>[if]チャージするなら → 壁裏</li><li>[if]チャージしたなら → 周りを見る</li><li>[if]撃つなら → 動き続ける</li></ul> |

### 余裕があったら考える

| 状況 | if-then |
| -- | -- |
| 上げ判断 | <ul><li>[if]ミカタ2人以上が前にいるなら → 索敵する</li><li>[if]人数有利だったら → ラインを上げる</li></ul> |
| 下げ判断 | <ul><li>[if]テキインクまで2本以内なら → 下がる</li><li>[if]画面にミカタの名前が無いなら → 下がる</li></ul> |
| ライン上げ下げ | <ul><li>[if]ライン上げるなら → 次の遮蔽まで</li><li>[if]ライン下げるなら → 一つ前の遮蔽まで</li></ul> |
| カバー | <ul><li>[if]チャージがあれば → 積極的にカバーする</li><li>[if]チャージがなければ → 無理をしない</li></ul> |
| 視線 | <ul><li>[if]自キャラを見てるなら → 視線を上げる</li><li>[if]射撃を開始するなら → レティクルを見る</li><li>[if]射撃してるなら → レティクル以外を見る</li><li>[if]画面上部に視線移動して索敵したなら → ついでにイカランプ見る</li></ul> |

### 突き詰めた立ち回り

```mermaid
%%{ 
  init: { 
    'theme': 'dark',
    'themeVariables': { 
      'edgeLabelBackground':'#000000'
    },
    'flowchart': { 'curve': 'linear' } 
  }
}%%

graph TB
  if_weigh_pros_cons{状況確認}
  if_weigh_pros_cons--2人有利-->action_push
  if_weigh_pros_cons--1人有利-->action_paint
  if_weigh_pros_cons--拮抗-->action_paint
  if_weigh_pros_cons--1人不利-->action_stay
  if_weigh_pros_cons--2人不利-->action_pull
  action_push[ラインを上げる]-->action_return
  action_paint[塗り広げ]-->action_return
  action_stay[塗り維持]-->action_return
  action_pull[ラインを下げる]-->action_return
  action_return[状況確認に戻る]

  %% 開始終了
  style if_weigh_pros_cons color:#ffffff,fill:#000000
  style action_return color:#ffffff,fill:#000000

  %% 2人有利
  style action_push color:#ffffff,fill:#00a300
  linkStyle 0 stroke-width:4px,stroke:#00a300,color:#00a300
  linkStyle 5 stroke-width:4px,stroke:#00a300

  %% 1人有利
  style action_paint color:#ffffff,fill:#7b8c00
  linkStyle 1 stroke-width:4px,stroke:#7b8c00,color:#7b8c00
  linkStyle 2 stroke-width:4px,stroke:#7b8c00,color:#7b8c00
  linkStyle 6 stroke-width:4px,stroke:#7b8c00

  %% 1人不利
  style action_stay color:#ffffff,fill:#d26900
  linkStyle 3 stroke-width:4px,stroke:#d26900,color:#d26900
  linkStyle 7 stroke-width:4px,stroke:#d26900

  %% 2人不利
  style action_pull color:#ffffff,fill:#dc143c
  linkStyle 4 stroke-width:4px,stroke:#dc143c,color:#dc143c
  linkStyle 8 stroke-width:4px,stroke:#dc143c
```

### その他、立ち回りのポイント

- Ver 11.0.0 で明確に塗りブキになった。生存しているだけでアド！！

<details>

<summary>Before Ver.11.0.0</summary>

- 全体
  - なんでもできるからなんでもしないとダメ
  - 相手より有利なところを押し付ける
    - 有利なところとはバレルとしての絶対的な強みではなく、相手と比較した時の相対的な強み
- 立ち位置
  - 射程のギリギリを押し付ける
    - 撃ちながら前に進むのが強い。バレルの弾の先で突くイメージ
    - クリアリングも射程の先で行う ※耳かきみたいなイメージ
  - 前線を上げるんじゃなくて、前線に引っ張られる状況もある
  - 長射程がいたら高台に立たない
- 戦術
  - 敵が来たそうな場所に置き撃ちしておく
  - 「テキが通りたいルートを塗っておく」「ミカタが通りたいルートを塗っておく」
  - 射程とダメージ変わらないので一周チャージを有効に使う
    - フルチャばっかりだと前線に遅れるし、インクもなくなりがち
  - 短射程相手はブロック近くで戦わず、広い場所で射程を押し付ける
  - チャージ中に周りを見る
  - チャージャーは倒すんじゃなくて倒されないのスタンス
  - チャージ中は遅いので、チャージしながら歩くよりもイカで移動してからチャージしたほうが早い
    - 移動して壁裏チャージ
  - 塗りはまん中が抜ける
    - 最高効率は上下に振ることだが実戦だと難しいので円を描くように回す
    - 塗りブキと正面から塗り合っても勝てないので塗りブキは先端で撃って追い払うことを重視

</details>

## 気をつけるブキ

| ポイント | 詳細 | ブキ |
| ---- | -- | -- |
| バレル以上の射程 | > 0.5 | <img src="./images/Bamboozler14MkI.png" width="32px" alt="竹" /><img src="./images/GooTuber.png" width="32px" alt="ソイチュ" />(<img src="./images/SplatanaStamper.png" width="32px" alt="ジム" />)<img src="./images/JetSquelcher.png" width="32px" alt="ジェット" /><img src="./images/Explosher.png" width="32px" alt="エクス" />(<img src="./images/BallpointSplatling.png" width="32px" alt="クゲ" />) |
| - | >= 1.0 | (<img src="./images/DynamoRoller.png" width="32px" alt="ダイナモ" />)<img src="./images/HydraSplatling.png" width="32px" alt="ハイドラ" /><img src="./images/SplatCharger.png" width="32px" alt="スプチャ" /><img src="./images/WellstringV.png" width="32px" alt="フルイド" /><img src="./images/Bloblobber.png" width="32px" alt="オフロ" /> |
| - | >= 1.5 | (<img src="./images/SplatCharger.png" width="32px" alt="チャースコ" />)<img src="./images/Snipewriter5H.png" width="32px" alt="ペン" /><img src="./images/Tri-Stringer.png" width="32px" alt="弓" /> |
| - | >= 2.0 | <img src="./images/E-liter4K.png" width="32px" alt="リッター" /> |
| バレル以上のキル速 | !1確 | <img src="./images/Sploosh-O-Matic.png" width="32px" alt="ボールド" /><img src="./images/52Gal.png" width="32px" alt="52" /><img src="./images/DappleDualies.png" width="32px" alt="スパ" />(<img src="./images/GloogaDualies.png" width="32px" alt="ケルビン" />)(<img src="./images/BallpointSplatling.png" width="32px" alt="クゲ" />)<img src="./images/Nautilus47.png" width="32px" alt="ノーチ" />(<img src="./images/HydraSplatling.png" width="32px" alt="ハイドラ" />) |
| バレル以上の射程&キル速 | - | <img src="./images/HydraSplatling.png" width="32px" alt="ハイドラ" /><img src="./images/SplatCharger.png" width="32px" alt="スプチャ" /><img src="./images/Tri-Stringer.png" width="32px" alt="弓" /><img src="./images/E-liter4K.png" width="32px" alt="リッター" /> |

## links

- [勝ち負けトントンから始めるスプラトゥーン攻略メモ with バレルスピナー](https://note.com/kotaoue/n/nb3e3413c759a)
- [スピナーで簡単にキルが取れる方法解説【バレルスピナー】](https://www.youtube.com/watch?v=kH2dufXwR0o)
- [【スピナー】誰でも上級者の動きになるキャラコン解説【スプラ３】](https://www.youtube.com/watch?v=-rm13bAuax8)
- [【必見】王冠勢がバレルで「最も効率よく」塗る方法教えます【バレルスピナー】【バレリミ】【ガチエリア】](https://www.youtube.com/watch?v=SpR-de0AIMM)
- [有効キルは立ち位置で決まる【オーバーフロッシャー】](https://www.youtube.com/watch?v=cbxmMGWsROc) ※ バレルの考え方とも一致
- [難所で崩れないハイドラの立ち回り。私の”安定ルーティン”公開【スプラトゥーン３】](https://www.youtube.com/watch?v=2iWkbRPVMPQ) ※ 下がり方が参考になる
- ステージ
  - [スプラ気をつけるポイント:バレル:ステージ](https://docs.google.com/spreadsheets/d/1JlSD02TxEMjEiwJZ0Kd5NAbjlgZKQplCbdxq-N9SC8s/edit?gid=1579070159#gid=1579070159)
  - [この動画を見るだけでヒラメが丘団地バレルスピナーの立ち回りが分かります【スプラトゥーン3】](https://www.youtube.com/watch?v=sSO1gZKmWus)
- [3D Model Heavy Splatling](https://sketchfab.com/3d-models/heavy-splatling-e7317dddfa594fbebd75dc3593455533)
