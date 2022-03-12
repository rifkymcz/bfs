## Disclaimer
pembuat bot ini tidak bertanggung jawab jika anda kena banned/blokir shopee karena bot ini.

## Install
### Termux
copas command dibawah
```
pkg install wget && wget -O bfs https://github.com/alimsk/bfs/releases/latest/download/bfs && chmod +x bfs
```

### Linux
kalo anda sudah menginstall golang versi 1.17 keatas, run:
```
go install github.com/alimsk/bfs@latest
```
kalo belum:
```
sudo apt install wget && wget -O bfs https://github.com/alimsk/bfs/releases/latest/download/bfs && chmod +x bfs
```

## Cara Pake
ambil cookie dari browser, terus simpan ke file `cookie` didalam folder yang sama dengan file `bfs`.

hasilnya akan terlihat seperti ini
```
$ ls
bfs  cookie
```

terus run botnya seperti biasa
```
$ ./bfs
```
run `bfs -h` buat nampilin semua opsi.

file `cookie` bisa berisi cookie biasa seperti yg ada di request header,
atau `[]*http.Cookie` yang diserialisasi menggunakan [gob](https://pkg.go.dev/encoding/gob).

#
kalo nemu bug pada bot ini, bisa buat issue [disini](https://github.com/alimsk/bfs/issues/new).

speed masih lemot, yg penting jadi dulu.\
next ane percepat waktu checkoutnya.
