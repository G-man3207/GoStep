# GoStep 🎥

Windows app for screen recording and input capture. Built with Fyne. Replaces Windows Step Recorder.

## ✨ Features

- 🖱️ Mouse tracking + clicks  
- 📸 Screenshots with highlights  
- 📄 HTML/PDF export  
- 🎨 Fyne UI  
- 🔒 No keyboard capture  
- 📜 MIT License  

## 🔧 Requirements

- **Run**: Windows  
- **Build**: Linux/WSL, Go 1.23+, Git, GCC, MinGW-w64  

## 📥 Install

1. Clone: `git clone https://github.com/G-man3207/GoStep.git && cd GoStep`  
2. Dependencies: `go mod download`  
3. Build (Linux/WSL): `./build.sh`  

## 🎮 Usage

1. Run `gostep.exe` (Windows)  
2. Hit "Start Recording"  
3. Do stuff (mouse only)  
4. Hit "Stop Recording"  
5. Edit steps:  
   - Tweak text  
   - Delete/reorder  
6. Export: HTML or PDF  
7. Files in `Documents/GoStep`  

## 📊 Output

- 🌐 **HTML**: Web page, interactive  
- 📑 **PDF**: Steps with screenshots + timestamps  

## 🛠️ Build

Linux/WSL: `chmod +x build.sh && ./build.sh`  
Outputs `gostep.exe`.  

## ❓ Troubleshooting

- **"Can’t run"**:  
  - Use `gostep.exe`  
  - Run as admin  
  - Unblock (Properties)  
  - Exclude in Defender  

## 🤝 Contribute

PRs on GitHub.  

## 📜 License

MIT. See [LICENSE](LICENSE).  

## 🪟 Windows Notes

- Unsigned app warnings:  
  - Unblock: Properties  
  - Admin mode  
  - Defender exclusion  
