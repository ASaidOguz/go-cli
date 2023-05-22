## CLI for checking Mem-usage in windovs-os  

Using [cobra](https://github.com/spf13/cobra) cli package to get windovs os mem usage values and write into designated text file in desired folder

## Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o memCli-windows.exe

after the build use .\memCli-windows.exe <--foldername--> <--filename-->

all mem usage will be written on your file