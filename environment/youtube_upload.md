# OBS で録画した動画を YouTube にアップする手順

[obs_gc551g2.md](./obs_gc551g2.md) の手順で録画した `.mov` ファイルを YouTube にアップするための手順をまとめる。

## 無料でよく使われる動画編集ツール

| ツール | 対応 OS | 特徴 |
| --- | --- | --- |
| [iMovie](https://www.apple.com/jp/imovie/) | macOS / iOS | Mac に標準搭載。直感的な UI でカット編集・BGM 追加が簡単にできる。YouTube への直接書き出しにも対応 |
| [DaVinci Resolve](https://www.blackmagicdesign.com/jp/products/davinciresolve) | macOS / Windows / Linux | 無料版でも機能が豊富。カラーグレーディングや音声編集も本格的にできる |
| [HandBrake](https://handbrake.fr/) | macOS / Windows / Linux | 動画変換専用ツール。`.mov` を YouTube に最適な `.mp4`（H.264）に変換するのに適している |

## .mov を MP4 に変換する（推奨）

YouTube は `.mov` を直接アップロードできるが、`.mp4`（H.264 + AAC）に変換してからアップロードすると安定しやすい。

### HandBrake を使う場合

1. HandBrake を起動し、録画した `.mov` ファイルをドラッグ＆ドロップする
2. **Preset** で `YouTube HQ 1080p60` を選択する
3. **Format** が `MP4` になっていることを確認する
4. **Video** タブ:
   - **Encoder**: `H.264 (x264)` または `H.265 (x265)`
   - **Framerate**: `60`
5. **Audio** タブ:
   - **Codec**: `AAC`
6. **Start Encode** をクリックして変換する

### ffmpeg を使う場合（CLI）

```sh
ffmpeg -i input.mov -c:v libx264 -crf 18 -preset slow -c:a aac -b:a 192k -movflags +faststart output.mp4
```

| オプション | 内容 |
| --- | --- |
| `-c:v libx264` | 映像コーデックに H.264 を使用 |
| `-crf 18` | 品質設定（0 が最高品質・ロスレス、51 が最低品質。18〜23 が高品質出力の推奨範囲） |
| `-preset slow` | エンコード速度（`slow` にすると品質が上がるが時間がかかる） |
| `-c:a aac -b:a 192k` | 音声コーデックに AAC を使用し、ビットレートを 192kbps に設定 |
| `-movflags +faststart` | moov atom をファイル先頭に移動し、YouTube へのアップロード前から再生できるよう最適化する |

## YouTube へのアップロード手順

1. [YouTube Studio](https://studio.youtube.com/) にアクセスしてログインする
2. 右上の **作成** → **動画をアップロード** をクリックする
3. 変換した `.mp4` ファイルをドラッグ＆ドロップする
4. 以下の情報を入力する:
   - **タイトル**: 内容がわかりやすいタイトルをつける
   - **説明**: 使用機材・ゲーム名・プレイ内容などを記載する
   - **サムネイル**: 任意の画像を設定すると視認性が上がる
   - **公開設定**: `非公開` / `限定公開` / `公開` から選ぶ
5. **次へ** を数回クリックし、最後に **保存** または **公開** をクリックする

## アップロード時のポイント

- **解像度・フレームレート**: 1080p 60fps でアップロードすると高品質な動画になる
- **処理時間**: アップロード直後は画質が粗いことがある。YouTube 側の処理が完了するまで数分〜数十分かかる場合がある
- **ファイルサイズ**: YouTube のアップロード上限は 256GB または 12 時間。長時間の録画は事前にカットしておくと良い
- **著作権**: ゲーム音楽や BGM は著作権に注意する。ゲーム実況については各メーカーのガイドラインを確認する

## 参考

- [YouTube ヘルプ: アップロードできる動画形式](https://support.google.com/youtube/troubleshooter/2888402)
- [YouTube ヘルプ: 動画の詳細設定](https://support.google.com/youtube/answer/57404)
