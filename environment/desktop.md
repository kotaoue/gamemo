# デスク回りの設定
## current
```mermaid
flowchart LR
    Web <-->|USB| Switch
    Switch --->|HDMI| Monitor["モニター (外部出力なし)"]
    Switch -.->|3.5mm| ヘッドフォン["ヘッドフォン WH-1000XM3"]
    Web <-->|無線| MacBook
    MacBook --->|DisplayPort| Monitor
    MacBook -.->|3.5mm| Meeting{会議中}
    Meeting -.->|Yes| イヤフォン["イヤフォン WI-1000XM2"]
    Meeting -.->|No| ヘッドフォン
```

---

## 理想
```mermaid
flowchart LR
    Web <-->|USB| Switch
    Switch --->|HDMI| Monitor["モニター (外部出力なし)"]
    Switch -->|"ケーブル A"| Gadget
    Web <-->|無線| MacBook
    MacBook --->|DisplayPort| Monitor
    MacBook -->|"ケーブル B"| Gadget[なんかすごいアイテム]
    Gadget -->|"ケーブル C"| イヤフォン["イヤフォン WI-1000XM2"]
    Gadget -->|"ケーブル D"| ヘッドフォン["ヘッドフォン WH-1000XM3"]
```

---

## 変更したいポイント
* 音回り
  * 3.5mmのつなぎ直しを減らしたい
  * switch本体に3.5mm刺すのをやめたい
* カメラ回り
  * 内蔵カメラから外付けカメラにしたい
  * 外付けカメラにしたらマイクもほしい
  * 照明を電球からLEDに変更したい

## 現状ポイント
* MacBookのUSBポートは、電源、キーボード
* SwitchドックのUSBポートが有線LANで埋まってる
* モニタからHDMIとかスピーカーは出力できない
* デスク上の余剰スペースは15*25cmくらい
