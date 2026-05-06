# OBS で録画した動画を YouTube にアップする手順

[obs_gc551g2.md](./obs_gc551g2.md) の手順で録画した `.mov` ファイルを YouTube にアップするための手順をまとめる。

## 無料でよく使われる動画編集ツール

| ツール | 対応 OS | 特徴 |
| --- | --- | --- |
| [iMovie](https://www.apple.com/jp/imovie/) | macOS / iOS | Mac に標準搭載。直感的な UI でカット編集・BGM 追加が簡単にできる。YouTube への直接書き出しにも対応 |
| [DaVinci Resolve](https://www.blackmagicdesign.com/jp/products/davinciresolve) | macOS / Windows / Linux | 無料版でも機能が豊富。カラーグレーディングや音声編集も本格的にできる |
| [HandBrake](https://handbrake.fr/) | macOS / Windows / Linux | 動画変換専用ツール。`.mov` を YouTube に最適な `.mp4`（H.264）に変換するのに適している |

## 発展的な動画編集：ゆっくり実況・複数動画合成

### ゆっくり実況スタイルの編集（画像・吹き出し追加）

ゆっくり実況のように画像キャラクターや吹き出しを動画に重ねたい場合は、以下のツールが使われている。

| ツール | 対応 OS | 特徴 |
| --- | --- | --- |
| [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) | Windows | 日本の実況・解説動画界で定番の無料動画編集ソフト。拡張編集プラグインを導入することで画像・テキスト・吹き出しをタイムライン上で自由に配置できる |
| [ゆっくりMovieMaker4（YMM4）](https://manjubox.net/ymm4/) | Windows | ゆっくり実況動画の作成に特化した無料ツール。VOICEVOX / AquesTalk 等の音声合成と連携し、キャラクター画像・吹き出し・字幕を簡単に配置できる |
| [VOICEVOX](https://voicevox.hiroshiba.jp/) | macOS / Windows / Linux | 無料の音声合成ソフト。ずんだもん・四国めたん等のキャラクターボイスを生成でき、YMM4 や AviUtl と組み合わせて使う |
| [DaVinci Resolve](https://www.blackmagicdesign.com/jp/products/davinciresolve) | macOS / Windows / Linux | Fusion ページを使うと画像・テキスト・吹き出し等の合成（コンポジット）が可能。慣れれば高品質な仕上がりになる |

> **おすすめの組み合わせ（Windows）**: VOICEVOX で音声生成 → YMM4 でキャラクター・吹き出しを配置 → 動画を書き出し

### 複数動画の結合・2 画面表示（ピクチャー・イン・ピクチャー）

複数の動画を 1 本にまとめたり、メイン動画＋サブ動画を 1 画面に同時表示したい場合は以下の方法がある。

| ツール | 対応 OS | 機能 |
| --- | --- | --- |
| [ffmpeg](https://ffmpeg.org/) | macOS / Windows / Linux | CLI ツール。動画の連結・フィルタを使った 2 画面合成（PiP）が可能。自由度が高い |
| [DaVinci Resolve](https://www.blackmagicdesign.com/jp/products/davinciresolve) | macOS / Windows / Linux | タイムライン上に複数トラックを並べてオーバーレイするだけで 2 画面合成が実現できる。GUI 操作で完結 |
| [OBS Studio](https://obsproject.com/ja) | macOS / Windows / Linux | 録画・配信ソフトだが、シーン機能を使って複数映像ソースを配置した状態で録画することで 2 画面構成の動画を作れる |
| [iMovie](https://www.apple.com/jp/imovie/) | macOS / iOS | ピクチャー・イン・ピクチャー機能を標準搭載。サブ動画を小窓で重ねる程度のシンプルな 2 画面合成ならこれで十分 |

#### ffmpeg で動画を連結する（複数動画 → 1 本）

```sh
# 連結したいファイルのリストを作成
printf "file 'part1.mp4'\nfile 'part2.mp4'\nfile 'part3.mp4'\n" > filelist.txt

# 再エンコードなしで連結（同じ解像度・コーデックの動画同士）
ffmpeg -f concat -safe 0 -i filelist.txt -c copy output.mp4
```

#### ffmpeg で 2 画面表示（ピクチャー・イン・ピクチャー）

メイン動画を背景（フルスクリーン）に、右上にサブ動画を小窓で重ねる例（1920×1080 出力）:

```sh
ffmpeg \
  -i main.mp4 \
  -i sub.mp4 \
  -filter_complex "
    [1:v]scale=480:270[sub];
    [0:v][sub]overlay=W-w-20:20
  " \
  -c:v libx264 -crf 18 -preset slow \
  -c:a aac -b:a 192k \
  output_pip.mp4
```

| パラメータ | 内容 |
| --- | --- |
| `[1:v]scale=480:270` | サブ動画を 480×270（元サイズの 1/4）に縮小 |
| `overlay=W-w-20:20` | `W-w-20` = 右端から 20px、`20` = 上端から 20px に配置（右上） |
| `overlay=20:H-h-20` | 左端から 20px・下端から 20px に配置したい場合（左下）はこちら |

> サブ動画を **左下** に置きたい場合は `overlay=20:H-h-20` に変更する。メインとサブを入れ替えたい場合は `-i` の順序を逆にする。

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
- [スプラトゥーン3 動画投稿・配信ガイドライン](./splatoon3_guidelines.md)
- [YouTube チャンネルハンドル候補](./youtube_channel_handles.md)
