# CMDs

## SurrealDB server
```ps1 Powershell
surreal start --bind 0.0.0.0:6565 --username ezy --password 1234 surrealkv://C:\Users\EzY\Desktop\ezMsg.db
```

## FlatBuffers generate
```ps1 Powershell
$fbsFiles = (Get-ChildItem -Filter *.fbs -Recurse).FullName
$includePaths = Get-ChildItem -Directory -Recurse | ForEach-Object { "-I", $_.FullName }
flatc --go -o ./generated @includePaths -I . $fbsFiles
```

flatc --go -I .\communication\ -o ./generated/ .\communication\outer_shell.fbs 
flatc --go -I .\communication\ -o ./generated/ .\communication\message.fbs