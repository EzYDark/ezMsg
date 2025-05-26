# Android build (for Emulator):
```ps1 PowerShell
gogio -target android -o ./target/android/ezMsg.apk .;
adb uninstall com.github.ezMsg;
adb install .\target\android\ezMsg.apk;
adb shell am start -n com.github.ezMsg/org.gioui.GioActivity;
adb logcat com.github.ezMsg:I *:S
```


# Protobuf generate:
```ps1 PowerShell
protoc --go_out=. --go_opt=default_api_level=API_OPAQUE --go_opt=paths=source_relative .\_test\zero_knowledge_com\protobuf\message.proto
```

```ps1 PowerShell
protoc --go_out=. --go_opt=default_api_level=API_OPAQUE --go_opt=paths=source_relative .\protobuf\message.proto
```
