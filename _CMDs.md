# Android build (for Emulator):

```ps1 PowerShell
gogio -target android -o ./target/android/ezMsg.apk .;
adb uninstall com.github.ezMsg;
adb install .\target\android\ezMsg.apk;
adb shell am start -n com.github.ezMsg/org.gioui.GioActivity;
adb logcat com.github.ezMsg:I *:S
```
