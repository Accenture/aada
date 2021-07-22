all: aada_mac_x64

release: aada_mac_x64.zip aada_mac_m1.zip aada_win_x64.exe aada_win_arm.exe aada_linux_x64 aada_linux_arm32 aada_linux_arm64
	cp aada_mac_m1.zip ../binaries/
	cp aada_mac_x64.zip ../binaries/
	cp -f aada_win_x64.exe aada.exe
	zip -j ../binaries/aada_win_x64.zip aada.exe
	cp -f aada_win_arm.exe aada.exe
	zip -j ../binaries/aada_win_arm.zip aada.exe
	cp -f aada_linux_x64 aada
	zip -j ../binaries/aada_linux_x64.zip aada
	cp -f aada_linux_arm32 aada
	zip -j ../binaries/aada_linux_arm32.zip aada
	cp -f aada_linux_arm64 aada
	zip -j ../binaries/aada_linux_arm64.zip aada

aada_mac_x64.zip: aada_mac_x64 apple.hcl
	cp -f aada_mac_x64 aada
	gon apple.hcl
	mv aada.zip aada_mac_x64.zip

aada_mac_m1.zip: aada_mac_arm64 apple.hcl
	cp -f aada_mac_arm64 aada
	gon apple.hcl
	mv aada.zip aada_mac_m1.zip

aada_mac_arm64: *.go go.mod go.sum
	GOOS=darwin GOARCH=arm64 go build -o aada_mac_arm64 .

aada_mac_x64: *.go go.mod go.sum
	GOOS=darwin GOARCH=amd64 go build -o aada_mac_x64 .

aada_win_x64.exe: *.go go.mod go.sum
	GOOS=windows GOARCH=amd64 go build -o aada_win_x64.exe .

aada_win_arm.exe: *.go go.mod go.sum
	GOOS=windows GOARCH=arm go build -o aada_win_arm.exe .

aada_linux_arm32: *.go go.mod go.sum
	GOOS=linux GOARCH=arm go build -o aada_linux_arm32 .

aada_linux_arm64: *.go go.mod go.sum
	GOOS=linux GOARCH=arm64 go build -o aada_linux_arm64 .

aada_linux_x64: *.go go.mod go.sum
	GOOS=linux GOARCH=amd64 go build -o aada_linux_x64 .

clean:
	rm -f aada_* aada aada.zip
