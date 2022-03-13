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
cara ambil cookie
1. login di browser seperti biasa
2. kalo sudah login, buka [shopee.co.id](https://shopee.co.id)
3. copy script dibawah dan pastekan ke kolom url
   ```js
   javascript:setTimeout(()=>navigator.clipboard.writeText(document.cookie),400)
   ```
4. NOTE: chrome secara otomatis menghapus teks `javascript:` ketika anda mem-paste script tersebut di kolom url,
   yang membuat scriptnya tidak berjalan. ketik aja secara manual kalo gak ada.

lalu run command 
```
echo -n "PASTE COOKIE DISINI" | ./bfs login
```
jangan lupa ganti teks `PASTE COOKIE DISINI` dengan cookie yang barusan anda copy.\
jika login sukses username akan ditampilkan di terminal.

lalu run botnya seperti biasa
```
$ ./bfs
```
run `bfs -h` buat nampilin semua opsi.

### Notes
- lakukan langkah yg sama ketika cookie expired.
- login cukup sekali saja, jika sudah login, gak perlu login lagi setiap mau running.

#
kalo nemu bug pada bot ini, bisa buat issue [disini](https://github.com/alimsk/bfs/issues/new).

speed masih lemot, yg penting jadi dulu.\
next ane percepat waktu checkoutnya.
