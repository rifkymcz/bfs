## Disclaimer
pembuat bot ini tidak bertanggung jawab jika anda kena banned/blokir shopee karena bot ini.

## Install
### Termux
copas command dibawah
```
pkg install wget && wget -O bfs https://github.com/alimsk/bfs/releases/latest/download/bfs-android-arm64 && chmod +x bfs
```
gunakan command yang sama ketika mau update.

### Linux
kalo anda sudah menginstall golang versi 1.17 keatas, run:
```
go install github.com/alimsk/bfs@latest
```
kalo belum:
```
sudo apt install wget && wget -O bfs https://github.com/alimsk/bfs/releases/latest/download/bfs-linux-amd64 && chmod +x bfs
```

## Cara Pake
cek [link tutorial](https://youtu.be/DcRKPCcx-Bs).

ambil cookie dari browser, lalu login dengan command dibawah
```
echo -n "PASTE COOKIE DISINI" | ./bfs login
```
jangan lupa ganti teks `PASTE COOKIE DISINI` dengan cookie yg barusan anda ambil.\
jika login sukses username akan ditampilkan di terminal.

lalu run botnya seperti biasa
```
./bfs
```
run `bfs -h` buat nampilin semua opsi.

### Notes
- lakukan langkah yg sama ketika cookie expired.
- login cukup sekali saja, jika sudah login, gak perlu login lagi setiap mau running.

#
kalo nemu bug pada bot ini, bisa buat issue [disini](https://github.com/alimsk/bfs/issues/new).

speed masih lemot, yg penting jadi dulu.\
next ane percepat waktu checkoutnya.
