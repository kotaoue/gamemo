# デスク回りの設定
## 環境
```mermaid
flowchart LR
    Web <-->|無線| MacBook
    MacBook -->|Type-C| USBハブ
    USBハブ -->|USB| DAC["Sound Blaster X4"]
    USBハブ -->|USB| マイク["マイク T669"]
    MacBook -->|HDMI| HDMI分配器
    HDMI分配器 -->ゲーミングモニタ["ゲーミングモニタ(23.6) EX-LDGC242HT"]
    MacBook -->|"HDMI(Type-C)"| PCモニタ["PCモニタ(23.8) VTF2401"]

    ゲーミングモニタ -->|3.5mm| DAC
    DAC -->|3.5mm| ヘッドホン["ヘッドホン ATH-AD500X"]
    DAC -->|3.5mm| イヤフォン["イヤホン WI-1000XM2"]
    
    Web <-->|有線| Switch
    Switch -->|HDMI| HDMI分配器

    style モニタ fill:#00c
    style ゲーミングモニタ fill:#00c
    style PCモニタ fill:#00c

    style 端末 fill:#060
    style MacBook fill:#060
    style Switch fill:#060
```
## 変更したいポイント
* 音回り
  * 3.5mmのつなぎ直しを減らしたい
* カメラ回り
  * 内蔵カメラから外付けカメラにしたい
