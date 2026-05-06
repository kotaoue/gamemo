# デスク回りの設定

## 環境

```mermaid
flowchart LR
    MacBook -->|HDMI| HDMI分配器
    MacBook -->|"HDMI(Type-C)"| サブモニタ
    
    メインモニタ["メインモニタ(23.6) EX-LDGC242HT"]
    サブモニタ["サブモニタ(23.8) VTF2401"]
    HDMI分配器-->メインモニタ

    MacBook <--Type-C--> USBハブ
    subgraph USB機器
        キャプチャーデバイス["GC551G2"] -->|USB| USBハブ    
        USBハブ -->|USB| DAC["Sound Blaster X4"]
        USBハブ -->|USB| マイク["マイク HyperX SoloCast"]
        USBハブ -->|USB| カメラ["カメラ Logicool MX Brio"]
    end

    Switch2 -->|HDMI| HDMI分配器
    HDMI分配器 -->|HDMI| キャプチャーデバイス

    メインモニタ -->|3.5mm| DAC
    subgraph ヘッドセット
        DAC -->|3.5mm| ヘッドホン["ヘッドホン ATH-AD500X"]
        DAC -->|3.5mm| イヤフォン["イヤホン WI-1000XM2"]
    end
    
    style メインモニタ fill:#00c
    style サブモニタ fill:#00c
    style MacBook fill:#060
    style Switch2 fill:#060
```
