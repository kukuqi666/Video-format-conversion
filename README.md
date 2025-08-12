# 视频格式转换工具 / Video Format Conversion Tool

### 项目背景
灵感来源于给我爷爷下载戏曲， 因为唱戏的播放器有很多格式限制 没办法只能进行格式转换，我在网上找了些现成的工具，但是很符合国人习惯 一度的乱收费，查了一下壳是用c++和QT做的，本来想着逆向呢，但是还挺麻烦的 结果在GitHub上面找到了FFmpeg， 是一个开源的、功能强大的多媒体处理工具，它主要用于处理音频和视频文件，然后Video-Format-Conversion就诞生了，为大家提供免费的帮助，如果有懂golang的大佬可以二次开发，但是要标明原创作者，谢谢配合

### ffmpeg介绍
FFmpeg 是一个开源的、功能强大的多媒体处理工具，它主要用于处理音频和视频文件。以下是关于 FFmpeg 的详细介绍：

1. **功能**
- **音视频转换**：可以将一种格式的音视频文件转换为另一种格式，例如将 MP4 转换为 AVI，或将 WAV 转换为 MP3。
- **音视频剪辑**：可以对音视频文件进行裁剪，提取其中的片段。
- **音视频合并**：可以将多个音视频文件合并成一个完整的文件。
- **添加字幕**：可以将字幕文件嵌入到视频文件中。
- **调整音视频参数**：可以调整音视频的分辨率、帧率、比特率等参数。
- **屏幕录制**：可以用来录制屏幕操作，生成视频文件。
- **直播推流**：可以将本地的音视频内容推送到直播服务器，用于直播。

2. **组成**
- **FFmpeg**：核心程序，用于处理音视频文件的编解码和转换。
- **FFplay**：一个简单的音视频播放器，用于播放音视频文件。
- **FFprobe**：用于分析音视频文件的元数据，例如时长、分辨率、编码格式等信息。

3. **应用场景**
- **视频编辑**：在视频编辑软件中，FFmpeg 常被用作底层的音视频处理引擎。
- **在线视频平台**：用于视频的转码和格式适配，以便在不同设备和网络环境下播放。
- **直播系统**：用于直播内容的采集、编码和推流。
- **多媒体开发**：开发者可以使用 FFmpeg 的库（如 libavcodec、libavformat 等）来开发自己的多媒体应用程序。

4. **使用方式**
- **命令行工具**：FFmpeg 提供了强大的命令行接口，用户可以通过命令行来执行各种音视频处理操作。例如：
  - 转码：`ffmpeg -i input.mp4 -c:v libx264 output.avi`
  - 剪辑：`ffmpeg -i input.mp4 -ss 00:00:10 -t 00:01:00 output.mp4`
- **编程接口**：开发者可以通过编程语言（如 C、Python 等）调用 FFmpeg 的库来实现更复杂的音视频处理功能。

5. **优势**
- **开源免费**：FFmpeg 是开源软件，用户可以免费使用，也可以根据需要修改源代码。
- **跨平台**：支持 Windows、Linux、macOS 等多种操作系统。
- **功能强大**：几乎涵盖了所有常见的音视频处理功能，且性能出色。
- **社区支持**：拥有庞大的开发者社区，遇到问题时可以很容易地找到解决方案。

6. **安装**
- **Windows**：可以从 FFmpeg 官方网站下载预编译的二进制文件，解压后即可使用。
- **Linux**：可以通过包管理器（如 apt、yum）安装，例如：`sudo apt-get install ffmpeg`。
- **macOS**：可以使用 Homebrew 安装，命令为：`brew install ffmpeg`。


### 项目介绍
这是一个基于 Go、Gin 框架和 FFmpeg 构建的 Web 视频格式转换工具。用户可以通过浏览器上传大视频文件（最高支持 10GB），并将其转换为 MP4、AVI、MKV、MOV、WMV 和 RMVB 等格式。工具提供响应式的 Bootstrap 5 界面，通过 WebSocket 实现实时转换进度显示，并支持高质量输出和转换时间记录。

### 功能特性
- **大文件支持**：支持高达 10GB 的视频文件上传和转换。
- **多种格式**：支持转换为 MP4、AVI、MKV、MOV、WMV 和 RMVB。
- **实时进度**：通过进度条实时显示转换进度。
- **高质量输出**：AVI 格式使用 `libx264` 编码，CRF 18 确保接近无损质量。
- **响应式界面**：基于 Bootstrap 5，适配桌面和移动设备。
- **转换时间**：显示转换耗时（以秒为单位）。
- **跨平台**：通过 GitHub Actions 在 Linux 和 Windows 上测试。

---
#### 浏览器显示页面
[![1](https://origin.picgo.net/2025/08/12/1aac498a4c77adb34.png)](https://www.picgo.net/image/1.0WQuAw)

#### 终端页面
[![2](https://origin.picgo.net/2025/08/12/2e8b9fb358a268f5f.png)](https://www.picgo.net/image/2.0WQAbp)

#### 最初测试页面
[![3](https://origin.picgo.net/2025/08/12/3e064d7d35991af0f.png)](https://www.picgo.net/image/3.0WQEeh)

---


### 前置条件
- **Go**：1.16 或更高版本 ([下载](https://go.dev/dl/))。
- **FFmpeg**：安装 FFmpeg 和 ffprobe，确保在系统 PATH 中 ([下载](https://ffmpeg.org/download.html))。
- **Git**：用于克隆和管理仓库。
- **Debian系列安装ffmpeg** ：`sudo apt install ffmpeg`
### 安装步骤
1. **克隆仓库**：
   
   ```bash
   git clone https://github.com/yourusername/video-format-conversion.git
   cd video-format-conversion
   ```
   
2. **安装依赖**：
   
   ```bash
   go get -u github.com/gin-gonic/gin github.com/gorilla/websocket
   ```
   
3. **创建目录**：
   为上传和输出文件创建目录，并设置写权限：
   
   ```bash
   mkdir Uploads outputs
   chmod -R 755 Uploads outputs
   ```
   
4. **运行或者编译应用**：
   ```bash
   go run main.go
   ```
   
   ```bash
   go build -o Video-format-conversion.exe main.go
   ```
   

在浏览器访问 `http://localhost:8080`。


### 使用方法
1. **上传视频**：
   - 打开 `http://localhost:8080`。
   - 选择视频文件（例如 900MB 的 RMVB 文件）和输出格式（例如 AVI）。
   - 点击“开始转换”。

2. **监控进度**：
   - 转换过程中，进度条实时更新。
   - 控制台打印输入文件大小、总时长、进度百分比和转换耗时。

3. **下载结果**：
   - 转换完成后，下载输出文件（例如 `output_converted.avi`）。
   - 界面显示文件大小和转换耗时。

### 示例
- 输入：900MB RMVB 文件（720x576，2小时）。
- 输出：AVI 文件（~950MB，高质量，`-crf 18`）。
- 转换时间：约 23 秒（取决于硬件）。
- 控制台输出：
  ```
  输入文件: example.rmvb, 大小: 900.50 MB
  视频总时长: 7909.64 秒
  转换进度: 50.23%
  转换耗时: 23.45 秒
  输出文件: outputs/example_converted.avi, 大小: 950.30 MB
  ```

### 项目结构
```
video-format-conversion/
├── main.go              # 主程序逻辑
├── templates/
│   └── index.html       # Bootstrap 界面
├── Uploads/             # 上传文件目录
├── outputs/             # 输出文件目录
├── .gitignore           # Git 忽略文件
├── .github/workflows/   # GitHub Actions 工作流
└── README.md            # 自述文件
```

### 开发说明
- **依赖**：
  - `github.com/gin-gonic/gin`：Web 框架。
  - `github.com/gorilla/websocket`：WebSocket 进度更新。
  - FFmpeg：视频转换核心。
- **自定义**：
  - 在 `main.go` 中调整 `-crf`（如 15，文件更大）或 `-b:v`（如 4000k）。
  - 修改 `index.html` 自定义界面。

### GitHub Actions
项目包含 Linux 和 Windows 的 CI/CD 工作流：
- **Linux**：在 Ubuntu 上构建和测试，验证 Go 和 FFmpeg。
- **Windows**：在 Windows 上构建和测试，验证 FFmpeg 兼容性。
详情见 `.github/workflows/`。

### 贡献指南
欢迎贡献！请：
1. Fork 仓库。
2. 创建特性分支（`git checkout -b feature/your-feature`）。
3. 提交更改（`git commit -m "Add your feature"`）。
4. 推送分支（`git push origin feature/your-feature`）。
5. 提交 Pull Request。

### 许可证
[MIT License](LICENSE)

### 联系方式
如有问题或建议，请在 GitHub 上提交 Issue 或联系 [你的邮箱或用户名]。

---