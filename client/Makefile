all: aada_mac_x64 aada_mac_arm64 aada_win_x64.exe aada_linux_x64

release: aada_bundle
	cp aada ../binaries/mac/x64/aada
	cp aada.zip ../binaries/mac/x64/aada.zip
	cp aada_mac_arm64 ../binaries/mac/m1/aada
	cp aada_win_x64.exe ../binaries/windows/aada.exe
	zip -j ../binaries/windows/aada.zip ../binaries/windows/aada.exe
	cp aada_linux_x64 ../binaries/linux/aada

aada_bundle: aada apple.hcl aada_mac_x64
	cp aada_mac_x64 aada
	gon apple.hcl

aada_mac_arm64: *.go go.mod go.sum
	GOOS=darwin GOARCH=arm64 go build -o aada_mac_arm64 .

aada_mac_x64: *.go go.mod go.sum
	GOOS=darwin GOARCH=amd64 go build -o aada_mac_x64 .

aada_win_x64.exe: *.go go.mod go.sum
	GOOS=windows GOARCH=amd64 go build -o aada_win_x64.exe .

aada_linux_x64: *.go go.mod go.sum
	GOOS=linux GOARCH=amd64 go build -o aada_linux_x64 .

clean:
	rm -f aada_linux_x64 aada_win_x64.exe aada_mac_x64 aada aada.zip
