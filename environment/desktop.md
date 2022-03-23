# デスク回りの設定
## current
```mermaid
flowchart LR
    Web <-->|USB| Switch
    Switch -->|HDMI| Monitor["モニター (外部出力なし)"]
    Switch -.->|3.5mm| ヘッドフォン["ヘッドフォン WH-1000XM3"]
    Web <-->|無線| MacBook
    MacBook -->|DisplayPort| Monitor
    MacBook -.->|3.5mm| イヤフォン["イヤフォン WI-1000XM2"]
    MacBook -.->|3.5mm| ヘッドフォン
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
    MacBook --->|なにか| 外付けカメラ
    MacBook --->|なにか| 外付けマイク
    Gadget -->|"3.5mm A"| イヤフォン["イヤフォン WI-1000XM2"]
    Gadget -->|"3.5mm B"| ヘッドフォン["ヘッドフォン WH-1000XM3"]
    style Gadget fill:#900,stroke:#f00,stroke-width:4px
    style 外付けカメラ fill:#909
    style 外付けマイク fill:#909
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
